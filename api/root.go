package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
)

func proofURLToDBQuery(param string) string {
	type proofLookup struct {
		Proof []string `json:"proof"`
	}

	lookup := proofLookup{
		Proof: strings.Split(param, ","),
	}

	q, err := json.Marshal([]proofLookup{lookup})
	if err != nil {
		return ""
	}

	return string(q)
}

func (s *Server) GetRoot(w http.ResponseWriter, r *http.Request) {
	type rootResp struct {
		Root hexutil.Bytes `json:"root"`
	}

	var (
		ctx     = r.Context()
		proof   = r.URL.Query().Get("proof")
		dbQuery = proofURLToDBQuery(proof)
	)
	if proof == "" || dbQuery == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing list of proofs")
		return
	}

	const q = `
		SELECT root
		FROM trees
		WHERE proofs_array(proofs) @> proofs_array($1)
		LIMIT 1
	`

	rr := rootResp{}
	err := s.db.QueryRow(ctx, q, dbQuery).Scan(
		&rr.Root,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, err, http.StatusNotFound, "root not found for proofs")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting root")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	s.sendJSON(r, w, rr)
}
