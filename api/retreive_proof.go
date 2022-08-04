package api

import (
	"errors"
	"net/http"

	"github.com/contextart/al/api/db/queries"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type retreiveProofResponseBody struct {
	Proof []string `json:"proof"`
}

func (s *Server) RetrieveProof(w http.ResponseWriter, r *http.Request) {
	var (
		rootStr = r.URL.Query().Get("root")
		addrStr = r.URL.Query().Get("address")
	)
	if rootStr == "" || addrStr == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "need to provide both a merkle root and an address")
		return
	}

	var root common.Hash
	if err := root.UnmarshalText([]byte(rootStr)); err != nil {
		s.sendJSONError(r, w, err, http.StatusBadRequest, "invalid merkle root")
		return
	}

	var addr common.Address
	if err := addr.UnmarshalText([]byte(addrStr)); err != nil {
		s.sendJSONError(r, w, err, http.StatusBadRequest, "invalid address")
		return
	}

	proof, err := s.dbq.GetMerkleProof(r.Context(), queries.GetMerkleProofParams{
		Root:    root.Bytes(),
		Address: addr.Bytes(),
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "failed to retrieve proof")
		return
	}

	// cache for 1 year
	w.Header().Set("Cache-Control", "public, max-age=31536000")

	proofStrs := make([]string, 0, len(proof))
	for _, proofHash := range proof {
		proofStrs = append(proofStrs, common.BytesToHash(proofHash).String())
	}

	s.sendJSON(r, w, retreiveProofResponseBody{
		Proof: proofStrs,
	})
}
