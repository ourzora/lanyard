package api

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/contextwtf/lanyard/merkle"
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
		leaf = common.FromHex(r.URL.Query().Get("unhashedLeaf"))
		addr = common.HexToAddress(r.URL.Query().Get("address"))
	)
	if len(root) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}
	if r.URL.Query().Get("unhashedLeaf") == "" && r.URL.Query().Get("address") == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing leaf")
		return
	}

	td, err := getTree(ctx, s.db, root)
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, nil, http.StatusNotFound, "tree not found")
		w.Header().Set("Cache-Control", "public, max-age=60")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting proof")
		return
	}

	var (
		leaves [][]byte
		target []byte
	)
	// check if leaf is in tree and error if not
	for _, l := range td.UnhashedLeaves {
		if len(target) == 0 {
			if len(leaf) > 0 {
				if bytes.Equal(l, leaf) {
					target = l
				}
			} else if leaf2Addr(l, td.Ltd, td.Packed).Hex() == addr.Hex() {
				target = l
			}
		}

		leaves = append(leaves, l)
	}

	if len(target) == 0 {
		s.sendJSONError(r, w, nil, http.StatusNotFound, "leaf not found in tree")
		return
	}

	var (
		p    = merkle.New(leaves).Proof(target)
		phex = []hexutil.Bytes{}
	)

	// convert [][]byte to []hexutil.Bytes
	for _, p := range p {
		phex = append(phex, p)
	}

	// cache for 1 year if we're returning an unhashed leaf proof
	// or 60 seconds for an address proof
	if len(leaf) > 0 {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=60")
	}
	s.sendJSON(r, w, getProofResp{
		UnhashedLeaf: target,
		Proof:        phex,
	})
}
