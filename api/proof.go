package api

import (
	"bytes"
	"context"
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

type cachedTree struct {
	r getTreeResp
	t merkle.Tree
}

func (s *Server) getCachedTree(ctx context.Context, root common.Hash) (cachedTree, error) {
	r, ok := s.tlru.Get(root)
	if ok {
		return r, nil
	}

	td, err := getTree(ctx, s.db, root.Bytes())
	if err != nil {
		return cachedTree{}, err
	}

	leaves := [][]byte{}
	for _, l := range td.UnhashedLeaves {
		leaves = append(leaves, l[:])
	}

	t := merkle.New(leaves)
	ct := cachedTree{
		r: td,
		t: t,
	}

	s.tlru.Add(root, ct)
	return ct, nil
}

func (s *Server) GetProof(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		root = common.HexToHash(r.URL.Query().Get("root"))
		leaf = common.FromHex(r.URL.Query().Get("unhashedLeaf"))
		addr = common.FromHex(r.URL.Query().Get("address"))
	)

	if len(root) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing root")
		return
	}
	if len(leaf) == 0 && len(addr) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "missing leaf")
		return
	}

	ct, err := s.getCachedTree(ctx, root)
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, nil, http.StatusNotFound, "tree not found")
		w.Header().Set("Cache-Control", "public, max-age=60")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "selecting proof")
		return
	}

	var (
		target []byte
	)
	// check if leaf is in tree and error if not
	for _, l := range ct.r.UnhashedLeaves {
		if len(target) == 0 {
			if len(leaf) > 0 {
				if bytes.Equal(l, leaf) {
					target = l
				}
			} else if bytes.Equal(leaf2Addr(l, ct.r.Ltd, ct.r.Packed), addr) {
				target = l
			}
		}
	}

	if len(target) == 0 {
		s.sendJSONError(r, w, nil, http.StatusNotFound, "leaf not found in tree")
		return
	}

	var (
		p    = ct.t.Proof(ct.t.Index(target))
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
