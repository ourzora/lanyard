package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
)

func proofURLToDBQuery(param string) []string {
	q, err := json.Marshal(strings.Split(param, ","))
	if err != nil {
		return []string{}
	}

	// encode again
	q, err = json.Marshal(string(q))
	if err != nil {
		return []string{}
	}

	return []string{string(q)}
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
	if proof == "" || len(dbQuery) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing list of proofs")
		return
	}

	const q = `
		SELECT root
		FROM trees
		WHERE proofs_array(proofs) @> $1
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
