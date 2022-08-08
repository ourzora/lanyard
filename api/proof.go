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
		addr = r.URL.Query().Get("address")
	)
	if root == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}
	if leaf == "" && addr == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing leaf")
		return
	}

	var (
		proof [][]byte
		err   error
	)

	if leaf != "" {
		proof, err = s.dbq.SelectProofByUnhashedLeaf(r.Context(), queries.SelectProofByUnhashedLeafParams{
			Root:         common.FromHex(root),
			UnhashedLeaf: common.FromHex(leaf),
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "failed to retrieve proof")
			return
		}
	} else {
		rows, err := s.dbq.SelectProofByAddress(r.Context(), queries.SelectProofByAddressParams{
			Root:    common.FromHex(root),
			Address: common.HexToAddress(addr).Bytes(),
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "failed to retrieve proof")
			return
		}

		// since this isn't guaranteed to be unique, we take the first proof available
		if len(rows) > 0 {
			proof = rows[0]
		}
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
