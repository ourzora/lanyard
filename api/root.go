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
		WHERE proofs_array(proofs) <@ proofs_array($1);
	`
	rr := rootResp{}

	// should only return one row, using QueryFunc to verify
	// that's the case and return an error if not (we've had
	// issues with this in the past)
	hasRow := false
	_, err := s.db.QueryFunc(ctx, q, []interface{}{dbQuery}, []interface{}{&rr.Root}, func(qfr pgx.QueryFuncRow) error {
		if hasRow {
			return errors.New("multiple rows returned")
		}
		hasRow = true
		return nil
	})

	if len(rr.Root) == 0 {
		s.sendJSONError(r, w, nil, http.StatusNotFound, "root not found for proofs")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting root")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")
	s.sendJSON(r, w, rr)
}
