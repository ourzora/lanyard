package api

import (
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
)

type getProofResp struct {
	UnhashedLeaf hexutil.Bytes   `json:"unhashedLeaf"`
	Proof        []hexutil.Bytes `json:"proof"`
}

func (s *Server) GetProof(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		root = common.FromHex(r.URL.Query().Get("root"))
		leaf = r.URL.Query().Get("unhashedLeaf")
		addr = r.URL.Query().Get("address")
	)
	if len(root) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}
	if leaf == "" && addr == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing leaf")
		return
	}

	const q = `
		WITH tree AS (
			SELECT jsonb_array_elements(proofs) proofs
			FROM merkle_trees
			WHERE root = $1
		)
		SELECT
			proofs->'leaf',
			proofs->'proof'
		FROM tree
		WHERE (
			--eth addresses contain mixed casing to
			--accommodate checksums. we sidestep
			--the casing issues for user queries
			lower(proofs->>'addr') = lower($2)
			OR lower(proofs->>'leaf') = lower($3)
		)
	`
	var (
		resp = &getProofResp{}
		row  = s.db.QueryRow(ctx, q, root, addr, leaf)
		err  = row.Scan(&resp.UnhashedLeaf, &resp.Proof)
	)
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, nil, http.StatusNotFound, "tree not found")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting proof")
		return
	}

	// cache for 1 year if we're returning an unhashed leaf proof
	if leaf != "" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	s.sendJSON(r, w, resp)
}
