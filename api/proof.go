package api

import (
	"errors"
	"net/http"

	"github.com/contextart/al/api/db/queries"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type getProofResp struct {
	Proof []string `json:"proof"`
}

func (s *Server) GetProof(w http.ResponseWriter, r *http.Request) {
	var (
		root = r.URL.Query().Get("root")
		leaf = r.URL.Query().Get("unhashedLeaf")
	)
	if root == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}
	if leaf == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing leaf")
		return
	}
	proof, err := s.dbq.SelectProof(r.Context(), queries.SelectProofParams{
		Root:         common.Hex2Bytes(root),
		UnhashedLeaf: common.Hex2Bytes(leaf),
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "failed to retrieve proof")
		return
	}

	resp := &getProofResp{
		Proof: make([]string, 0, len(proof)),
	}
	for i := range proof {
		resp.Proof = append(resp.Proof, common.BytesToHash(proof[i]).String())
	}

	// cache for 1 year
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	s.sendJSON(r, w, resp)
}
