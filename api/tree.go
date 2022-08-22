package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/contextwtf/lanyard/merkle"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
	"github.com/ryandotsmith/jh"
)

func leaf2Addr(leaf []byte, ltd []string, packed bool) common.Address {
	if len(ltd) == 0 || (len(ltd) == 1 && ltd[0] == "address") {
		return common.BytesToAddress(leaf)
	}
	if packed {
		return addrPacked(leaf, ltd)
	}
	return addrUnpacked(leaf, ltd)
}

func addrUnpacked(leaf []byte, ltd []string) common.Address {
	var addrStart, pos int
	for _, desc := range ltd {
		if desc == "address" {
			addrStart = pos
			break
		}
		pos += 32
	}
	if len(leaf) >= addrStart+32 {
		return common.BytesToAddress(leaf[addrStart:(addrStart + 32)])
	}
	return common.Address{}
}

func addrPacked(leaf []byte, ltd []string) common.Address {
	var addrStart, pos int
	for _, desc := range ltd {
		t, err := abi.NewType(desc, "", nil)
		if err != nil {
			return common.Address{}
		}
		if desc == "address" {
			addrStart = pos
			break
		}
		pos += int(t.GetType().Size())
	}
	if addrStart == 0 && pos != 0 {
		return common.Address{}
	}
	if len(leaf) >= addrStart+20 {
		return common.BytesToAddress(leaf[addrStart:(addrStart + 20)])
	}
	return common.Address{}
}

type jsonNullBool struct {
	sql.NullBool
}

func (jnb *jsonNullBool) UnmarshalJSON(d []byte) error {
	var b *bool
	if err := json.Unmarshal(d, &b); err != nil {
		return err
	}
	if b == nil {
		jnb.Valid = false
		return nil
	}

	jnb.Valid = true
	jnb.Bool = *b
	return nil
}

func (jnb jsonNullBool) MarshalJSON() ([]byte, error) {
	if jnb.Valid {
		return json.Marshal(jnb.Bool)
	}
	return json.Marshal(nil)
}

func encodeProof(p [][]byte) []string {
	var res []string
	for i := range p {
		res = append(res, hexutil.Encode(p[i]))
	}
	return res
}

type createTreeReq struct {
	Leaves []hexutil.Bytes `json:"unhashedLeaves"`
	Ltd    []string        `json:"leafTypeDescriptor"`
	Packed jsonNullBool    `json:"packedEncoding"`
}

type createTreeResp struct {
	MerkleRoot string `json:"merkleRoot"`
}

func (s *Server) CreateTree(
	ctx context.Context,
	req createTreeReq,
) (*createTreeResp, error) {
	if len(req.Leaves) == 1 {
		return nil, apiError(http.StatusBadRequest, "Must provide at least 2 leaves")
	}

	//convert []hexutil.Bytes to [][]byte
	var leaves [][]byte
	for i := range req.Leaves {
		leaves = append(leaves, req.Leaves[i])
	}
	tree := merkle.New(leaves)

	type proofItem struct {
		Leaf  string   `json:"leaf"`
		Addr  string   `json:"addr"`
		Proof []string `json:"proof"`
	}
	var proofs = []proofItem{}
	for _, l := range req.Leaves {
		pf := tree.Proof(l)
		if !merkle.Valid(tree.Root(), pf, l) {
			return nil, apiError(http.StatusBadRequest, "generating proof for tree")
		}
		proofs = append(proofs, proofItem{
			Leaf:  hexutil.Encode(l),
			Addr:  leaf2Addr(l, req.Ltd, req.Packed.Bool).Hex(),
			Proof: encodeProof(pf),
		})
	}
	const q = `
		INSERT INTO trees(
			root,
			unhashed_leaves,
			ltd,
			packed,
			proofs
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (root)
		DO NOTHING
	`
	_, err := s.db.Exec(ctx, q,
		tree.Root(),
		req.Leaves,
		req.Ltd,
		req.Packed.NullBool,
		proofs,
	)
	if err != nil {
		return nil, apiError(http.StatusInternalServerError, "inserting tree")
	}

	return &createTreeResp{hexutil.Encode(tree.Root())}, nil
}

type getTreeResp struct {
	UnhashedLeaves []hexutil.Bytes `json:"unhashedLeaves"`
	LeafCount      int             `json:"leafCount"`
	Ltd            []string        `json:"leafTypeDescriptor"`
	Packed         jsonNullBool    `json:"packedEncoding"`
}

func (s *Server) GetTree(ctx context.Context) (*getTreeResp, error) {
	root := jh.Request(ctx).URL.Query().Get("root")
	if root == "" {
		return nil, apiError(http.StatusBadRequest, "missing root")
	}
	const q = `
		SELECT unhashed_leaves, ltd, packed
		FROM trees
		WHERE root = $1
	`
	tr := &getTreeResp{}
	err := s.db.QueryRow(ctx, q, common.FromHex(root)).Scan(
		&tr.UnhashedLeaves,
		&tr.Ltd,
		&tr.Packed,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apiError(http.StatusNotFound, "no tree found for root")
	} else if err != nil {
		return nil, apiError(http.StatusInternalServerError, "selecting tree")
	}

	tr.LeafCount = len(tr.UnhashedLeaves)
	jh.ResponseWriter(ctx).Header().Set("Cache-Control", "public, max-age=3600")
	return tr, nil
}
