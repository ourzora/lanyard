package api

import (
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
)

type getProofResp struct {
	UnhashedLeaf string          `json:"unhashedLeaf"`
	Proof        []hexutil.Bytes `json:"proof"`
}

func (s *Server) GetProof(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		rootStr = r.URL.Query().Get("root")
		root    = common.FromHex(rootStr)
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

	const q1 = `
		SELECT 1
		FROM merkle_trees
		WHERE root = $1;
	`
	var empty int
	err := s.db.QueryRow(ctx, q1, root).Scan(&empty)
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, nil, http.StatusNotFound, "tree not found")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "failed to check tree exists")
		return
	}

	resp := &getProofResp{}
	if leaf != "" {
		resp.UnhashedLeaf = leaf
		const q2 = `
			SELECT proof
			FROM merkle_proofs
			WHERE root = $1
			AND unhashed_leaf = $2
		`
		row := s.db.QueryRow(ctx, q2, root, common.FromHex(leaf))
		err = row.Scan(&resp.Proof)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting proof for unhashedLeaf")
			return
		}
	} else {
		// since this isn't guaranteed to be unique,
		// we take the first proof available
		const q3 = `
			SELECT proof, unhashed_leaf
			FROM merkle_proofs
			WHERE root = $1
			AND address = $2
			LIMIT 1
		`
		row := s.db.QueryRow(ctx, q3, root, common.HexToAddress(addr).Bytes())
		err = row.Scan(&resp.Proof, &resp.UnhashedLeaf)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting proof for address")
			return
		}
	}

	// cache for 1 year if we're returning an unhashed leaf proof
	if leaf != "" {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	s.sendJSON(r, w, resp)
}
