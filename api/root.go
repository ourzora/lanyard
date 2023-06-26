package api

import (
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
)

func (s *Server) GetRoot(w http.ResponseWriter, r *http.Request) {
	type rootResp struct {
		Root hexutil.Bytes `json:"root"`
		Note string        `json:"note"`
	}

	type rootsResp struct {
		Roots []hexutil.Bytes `json:"roots"`
	}

	var (
		ctx   = r.Context()
		err   error
		proof = r.URL.Query().Get("proof")
		ps    = strings.Split(proof, ",")
		pb    = [][]byte{}
	)

	for _, s := range ps {
		var b []byte
		b, err = hexutil.Decode(s)
		if err != nil {
			break
		}
		pb = append(pb, b)
	}

	if len(pb) == 0 || err != nil {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing or malformed list of proofs")
		return
	}

	const q = `
		SELECT root
		FROM proofs_hashes
		WHERE hash = $1;
	`
	var (
		roots []hexutil.Bytes
		rb    hexutil.Bytes
		ph    = hashProof(pb)
	)

	_, err = s.db.QueryFunc(ctx, q, []interface{}{&ph}, []interface{}{&rb}, func(qfr pgx.QueryFuncRow) error {
		roots = append(roots, rb)
		return nil
	})

	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting root")
		return
	} else if len(roots) == 0 { // db.QueryFunc doesn't return pgx.ErrNoRows
		w.Header().Set("Cache-Control", "public, max-age=60")
		s.sendJSONError(r, w, nil, http.StatusNotFound, "root not found for proofs")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")

	if strings.HasPrefix(r.URL.Path, "/api/v1/roots") {
		s.sendJSON(r, w, rootsResp{Roots: roots})
	} else {
		// The original functionality of this endpoint, getting one root for
		// a given proof, is deprecated. This is because for smaller trees,
		// there are often collisions with the same root for different proofs.
		// This bit of code is for backwards compatibility.
		const note = `This endpoint is deprecated. For smaller trees, there are often collisions with the same root for different proofs. Please use the /v1/api/roots endpoint instead.`
		s.sendJSON(r, w, rootResp{Root: roots[0], Note: note})
	}
}
