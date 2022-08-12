package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/contextwtf/lanyard/api/db/queries"
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
	defer r.Body.Close()
	var req createTreeReq
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

	tx, err := s.db.Begin(r.Context())
	defer tx.Rollback(r.Context())
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	q := s.dbq.WithTx(tx)

	var leaves [][]byte
	for i := range req.Leaves {
		leaves = append(leaves, req.Leaves[i])
	}
	tree := merkle.New(leaves, merkle.SortPairs)
	err = q.InsertTree(r.Context(), queries.InsertTreeParams{
		Root:           tree.Root(),
		UnhashedLeaves: leaves,
		Ltd:            req.Ltd,
		Packed:         req.Packed.NullBool,
	})
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to insert merkle tree")
		return
	}

	for _, leaf := range leaves {
		proof := tree.Proof(leaf)
		if len(proof) == 0 {
			s.sendJSONError(r, w, nil, http.StatusBadRequest, "Must provide addresses that result in a proof")
			return
		}
		err := q.InsertProof(r.Context(), queries.InsertProofParams{
			Root:         tree.Root(),
			UnhashedLeaf: leaf,
			Address:      leaf2AddrBytes(leaf, req.Ltd, req.Packed.Bool),
			Proof:        proof,
		})
		if err != nil {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to persist merkle proofs")
			return
		}
	}

	err = tx.Commit(r.Context())
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to persist")
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
	root := r.URL.Query().Get("root")
	if root == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}
	row, err := s.dbq.SelectTree(r.Context(), common.FromHex(root))
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, err, http.StatusNotFound, "tree not found for root")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to select tree")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")

	var l []hexutil.Bytes
	for i := range row.UnhashedLeaves {
		l = append(l, row.UnhashedLeaves[i])
	}
	s.sendJSON(r, w, getTreeResp{
		UnhashedLeaves: l,
		LeafCount:      len(row.UnhashedLeaves),
		Ltd:            row.Ltd,
		Packed:         jsonNullBool{row.Packed},
	})
}
