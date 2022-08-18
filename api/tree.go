package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/contextwtf/lanyard/merkle"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
)

func (s *Server) TreeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.CreateTree(w, r)
		return
	case http.MethodGet:
		s.GetTree(w, r)
		return
	default:
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}
}

func leaf2AddrBytes(leaf []byte, ltd []string, packed bool) []byte {
	if len(ltd) == 0 || (len(ltd) == 1 && ltd[0] == "address") {
		return common.BytesToAddress(leaf).Bytes()
	}

	if packed {
		return addrPacked(leaf, ltd)
	}
	return addrUnpacked(leaf, ltd)
}

func addrUnpacked(leaf []byte, ltd []string) []byte {
	var addrStart, pos int
	for _, desc := range ltd {
		if desc == "address" {
			addrStart = pos
			break
		}

		pos += 32
	}
	if len(leaf) >= addrStart+32 {
		return common.BytesToAddress(leaf[addrStart:(addrStart + 32)]).Bytes()
	}
	return nil
}

func addrPacked(leaf []byte, ltd []string) []byte {
	var addrStart, pos int
	for _, desc := range ltd {
		t, err := abi.NewType(desc, "", nil)
		if err != nil {
			return nil
		}
		if desc == "address" {
			addrStart = pos
			break
		}

		pos += int(t.GetType().Size())
	}
	if addrStart == 0 && pos != 0 {
		return nil
	}
	if len(leaf) >= addrStart+20 {
		return common.BytesToAddress(leaf[addrStart:(addrStart + 20)]).Bytes()
	}
	return nil
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

type createTreeReq struct {
	Leaves []hexutil.Bytes `json:"unhashedLeaves"`
	Ltd    []string        `json:"leafTypeDescriptor"`
	Packed jsonNullBool    `json:"packedEncoding"`
}

type createTreeResp struct {
	MerkleRoot string `json:"merkleRoot"`
}

func (s *Server) CreateTree(w http.ResponseWriter, r *http.Request) {
	var (
		req createTreeReq
		ctx = r.Context()
	)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(r, w, err, http.StatusBadRequest, "unhashedLeaves must be a list of hex strings")
		return
	}
	if len(req.Leaves) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "No leaves provided")
		return
	}
	if len(req.Leaves) == 1 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "You must provide at least two values")
		return
	}

	dbtx, err := s.db.Begin(ctx)
	defer dbtx.Rollback(ctx)
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	var leaves [][]byte
	for i := range req.Leaves {
		leaves = append(leaves, req.Leaves[i])
	}
	tree := merkle.New(leaves)
	const q1 = `
		INSERT INTO merkle_trees (root, unhashed_leaves, ltd, packed)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (root)
		DO NOTHING
	`
	_, err = dbtx.Exec(ctx, q1, tree.Root(), leaves, req.Ltd, req.Packed)
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to insert merkle tree")
		return
	}

	var batch = &pgx.Batch{}
	for _, leaf := range leaves {
		proof := tree.Proof(leaf)
		if len(proof) == 0 {
			s.sendJSONError(r, w, nil, http.StatusBadRequest, "Must provide addresses that result in a proof")
			return
		}
		const q2 = `
			INSERT INTO merkle_proofs (root, unhashed_leaf, address, proof)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (root, unhashed_leaf)
			DO NOTHING;
		`
		batch.Queue(q2,
			tree.Root(),
			leaf,
			leaf2AddrBytes(leaf, req.Ltd, req.Packed.Bool),
			proof,
		)
	}
	br := dbtx.SendBatch(ctx, batch)
	for i := 0; i < len(leaves); i++ {
		_, err := br.Exec()
		if err != nil {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "inserting proofs")
			return
		}
	}
	err = br.Close()
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "inserting proof batch")
		return
	}

	err = dbtx.Commit(ctx)
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "committing tx")
		return
	}
	s.sendJSON(r, w, createTreeResp{
		MerkleRoot: fmt.Sprintf("0x%s", hex.EncodeToString(tree.Root())),
	})
}

type getTreeResp struct {
	UnhashedLeaves []hexutil.Bytes `json:"unhashedLeaves"`
	LeafCount      int             `json:"leafCount"`
	Ltd            []string        `json:"leafTypeDescriptor"`
	Packed         jsonNullBool    `json:"packedEncoding"`
}

func (s *Server) GetTree(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		root = r.URL.Query().Get("root")
	)
	if root == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}
	const q = `
		SELECT unhashed_leaves, ltd, packed
		FROM merkle_trees
		WHERE root = $1
	`
	tr := getTreeResp{}
	err := s.db.QueryRow(ctx, q, common.FromHex(root)).Scan(
		&tr.UnhashedLeaves,
		&tr.Ltd,
		&tr.Packed,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, err, http.StatusNotFound, "tree not found for root")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting tree")
		return
	}

	tr.LeafCount = len(tr.UnhashedLeaves)

	w.Header().Set("Cache-Control", "public, max-age=3600")
	s.sendJSON(r, w, tr)
}
