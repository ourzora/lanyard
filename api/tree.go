package api

import (
	"bytes"
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
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

//TODO(ryan): remove this once the table has been updated
func MigrateProofs(ctx context.Context, db *pgxpool.Pool) error {
	const q1 = `
		select root, unhashed_leaves, ltd, packed
		from merkle_trees
		where proofs is null
	`
	rows, err := db.Query(ctx, q1)
	if err != nil {
		return err
	}
	defer rows.Close()

	//for the updates
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //noop for successful commit

	var migrated int
	for rows.Next() {
		var (
			tree   merkle.Tree
			root   []byte
			leaves [][]byte
			ltd    []string
			packed sql.NullBool
		)
		err = rows.Scan(&root, &leaves, &ltd, &packed)
		if err != nil {
			return err
		}
		tree = merkle.New(leaves)
		if !bytes.Equal(root, tree.Root()) {
			return errors.New("mismatched root")
		}

		type proofItem struct {
			Leaf  string   `json:"leaf"`
			Addr  string   `json:"addr"`
			Proof []string `json:"proof"`
		}
		var proofs = []proofItem{}
		for _, l := range leaves {
			pf := tree.Proof(l)
			if !merkle.Valid(tree.Root(), pf, l) {
				return errors.New("invalid proof")
			}
			proofs = append(proofs, proofItem{
				Leaf:  hexutil.Encode(l),
				Addr:  leaf2Addr(l, ltd, packed.Bool).Hex(),
				Proof: encodeProof(pf),
			})
		}

		const q2 = `
			update merkle_trees
			set proofs = $1
			where root = $2
		`
		_, err = tx.Exec(ctx, q2, proofs, root)
		if err != nil {
			return err
		}
		migrated++
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	log.Ctx(ctx).Info().Int("migrated", migrated).Msg("success")
	return nil
}

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
	switch len(req.Leaves) {
	case 0:
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "No leaves provided")
		return
	case 1:
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "You must provide at least two values")
		return
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
			s.sendJSONError(r, w, nil, http.StatusBadRequest, "Unable to generate proof for tree")
			return
		}
		proofs = append(proofs, proofItem{
			Leaf:  hexutil.Encode(l),
			Addr:  leaf2Addr(l, req.Ltd, req.Packed.Bool).Hex(),
			Proof: encodeProof(pf),
		})
	}
	const q = `
		INSERT INTO merkle_trees(
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
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "inserting tree")
		return
	}

	s.sendJSON(r, w, createTreeResp{hexutil.Encode(tree.Root())})
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
