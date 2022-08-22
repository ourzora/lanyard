package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
	"github.com/ryandotsmith/jh"
)

type getProofResp struct {
	UnhashedLeaf hexutil.Bytes   `json:"unhashedLeaf"`
	Proof        []hexutil.Bytes `json:"proof"`
}

func (s *Server) GetProof(ctx context.Context) (*getProofResp, error) {
	var (
		root = common.FromHex(jh.Request(ctx).URL.Query().Get("root"))
		leaf = jh.Request(ctx).URL.Query().Get("unhashedLeaf")
		addr = jh.Request(ctx).URL.Query().Get("address")
	)
	if len(root) == 0 {
		return nil, apiError(http.StatusBadRequest, "missing root")
	}
	if leaf == "" && addr == "" {
		return nil, apiError(http.StatusBadRequest, "missing leaf")
	}

	const q = `
		WITH tree AS (
			SELECT jsonb_array_elements(proofs) proofs
			FROM trees
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
		return nil, apiError(http.StatusNotFound, "tree not found")
	} else if err != nil {
		return nil, apiError(http.StatusInternalServerError, "selecting proof")
	}

	// cache for 1 year if we're returning an unhashed leaf proof
	if leaf != "" {
		jh.ResponseWriter(ctx).Header().Set("Cache-Control", "public, max-age=31536000")
	}
	return resp, nil
}
