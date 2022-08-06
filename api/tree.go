package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/contextart/al/api/db/queries"
	"github.com/contextart/al/merkle"
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
		switch desc {
		case "address":
			addrStart = pos
			break
		default:
			pos += 32
		}
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
		switch desc {
		case "address":
			addrStart = pos
			break
		default:
			pos += int(t.GetType().Size())
		}
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

type createTreeReq struct {
	leaves []hexutil.Bytes `json:"unhashedLeaves"`
	ltd    []string        `json:"leafTypeDescriptor"`
	packed jsonNullBool    `json:"packedEncoding"`
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
	if len(req.leaves) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "No leaves provided")
		return
	}

	var leaves [][]byte
	for i := range req.leaves {
		leaves = append(leaves, req.leaves[i])
	}
	tree := merkle.New(leaves)
	err := s.dbq.InsertTree(r.Context(), queries.InsertTreeParams{
		Root:           tree.Root(),
		UnhashedLeaves: leaves,
		Ltd:            req.ltd,
		Packed:         req.packed.NullBool,
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
		err := s.dbq.InsertProof(r.Context(), queries.InsertProofParams{
			Root:         tree.Root(),
			UnhashedLeaf: leaf,
			Address:      leaf2Addr(leaf, req.ltd, req.packed.Bool).Bytes(),
			Proof:        proof,
		})
		if err != nil {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to persist merkle proofs")
			return
		}
	}

	s.sendJSON(r, w, createTreeResp{
		MerkleRoot: fmt.Sprintf("0x%s", hex.EncodeToString(tree.Root())),
	})
}

const maxAddressesPerPage = 10000

type getTreeResp struct {
	UnhashedLeaves []hexutil.Bytes `json:"unhashedLeaves"`
	LeafCount      int             `json:"leafCount"`
}

func (s *Server) GetTree(w http.ResponseWriter, r *http.Request) {
	root := r.URL.Query().Get("root")
	if root == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}

	leaves, err := s.dbq.SelectLeaves(r.Context(), common.Hex2Bytes(root))
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, err, http.StatusNotFound, "tree not found for root")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to select leaves")
		return
	}

	var l []hexutil.Bytes
	for i := range leaves {
		l = append(l, leaves[i])
	}
	s.sendJSON(r, w, getTreeResp{
		UnhashedLeaves: l,
		LeafCount:      len(leaves),
	})
}
