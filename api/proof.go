package api

import (
	"errors"
	"net/http"

	"github.com/contextwtf/lanyard/api/db/queries"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
)

type getProofResp struct {
	UnhashedLeaf *string  `json:"unhashedLeaf"`
	Proof        []string `json:"proof"`
}

func (s *Server) GetProof(w http.ResponseWriter, r *http.Request) {
	var (
		rootStr = r.URL.Query().Get("root")
		leaf    = r.URL.Query().Get("unhashedLeaf")
		addr    = r.URL.Query().Get("address")
	)
	if rootStr == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}
	if leaf == "" && addr == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing leaf")
		return
	}

	root := common.FromHex(rootStr)

	exists, err := s.dbq.SelectTreeExists(r.Context(), root)
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "failed to check tree exists")
		return
	}

	if !exists {
		s.sendJSONError(r, w, nil, http.StatusNotFound, "tree not found")
		return
	}

	var (
		unhashedLeaf *string
		proof        [][]byte
	)

	if leaf != "" {
		unhashedLeaf = &leaf
		proof, err = s.dbq.SelectProofByUnhashedLeaf(r.Context(), queries.SelectProofByUnhashedLeafParams{
			Root:         root,
			UnhashedLeaf: common.FromHex(leaf),
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "failed to retrieve proof")
			return
		}
	} else {
		rows, err := s.dbq.SelectProofByAddress(r.Context(), queries.SelectProofByAddressParams{
			Root:    root,
			Address: common.HexToAddress(addr).Bytes(),
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "failed to retrieve proof")
			return
		}

		// since this isn't guaranteed to be unique, we take the first proof available
		if len(rows) > 0 {
			ul := hexutil.Encode(rows[0].UnhashedLeaf)
			unhashedLeaf = &ul
			proof = rows[0].Proof
		}
	}

	resp := &getProofResp{
		Proof:        make([]string, 0, len(proof)),
		UnhashedLeaf: unhashedLeaf,
	}
	for i := range proof {
		resp.Proof = append(resp.Proof, common.BytesToHash(proof[i]).String())
	}

	// cache for 1 year if we're returning an unhashed leaf proof
	if leaf != "" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}

	s.sendJSON(r, w, resp)
}
