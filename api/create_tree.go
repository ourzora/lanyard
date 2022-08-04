package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/contextart/al/api/db/queries"
	"github.com/contextart/al/merkle"
	"github.com/ethereum/go-ethereum/common"
)

type createTreeRequestBody struct {
	AllowedAddresses []common.Address `json:"allowedAddresses"`
}

type createTreeResponseBody struct {
	MerkleRoot string `json:"merkleRoot"`
}

func (s *Server) CreateTree(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request createTreeRequestBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.sendJSONError(r, w, err, http.StatusBadRequest, "addresses must be a list of hex strings")
		return
	}

	if len(request.AllowedAddresses) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "No addresses provided")
		return
	}

	addrs := make([][]byte, 0, len(request.AllowedAddresses))
	for _, addr := range request.AllowedAddresses {
		addrs = append(addrs, addr.Bytes())
	}

	tree := merkle.New(addrs)

	err := s.dbq.InsertMerkleTree(r.Context(), queries.InsertMerkleTreeParams{
		Root:      tree.Root(),
		Addresses: addrs,
	})
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to insert merkle tree")
		return
	}

	for _, addr := range addrs {
		proof := tree.Proof(addr)

		if len(proof) == 0 {
			s.sendJSONError(r, w, nil, http.StatusBadRequest, "Must provide addresses that result in a proof")
			return
		}

		err := s.dbq.InsertMerkleProof(r.Context(), queries.InsertMerkleProofParams{
			Root:    tree.Root(),
			Address: addr,
			Proof:   proof,
		})
		if err != nil {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to persist merkle proofs")
			return
		}
	}

	s.sendJSON(r, w, createTreeResponseBody{
		MerkleRoot: fmt.Sprintf("0x%s", hex.EncodeToString(tree.Root())),
	})
}
