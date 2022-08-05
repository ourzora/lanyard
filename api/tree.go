package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/contextart/al/api/db/queries"
	"github.com/contextart/al/merkle"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

func (s *Server) TreeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.CreateTree(w, r)
		return
	case http.MethodGet:
		s.GetTree(w, r)
		return
	default:
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}
}

type createTreeReq struct {
	AllowedAddresses []common.Address `json:"allowedAddresses"`
}

type createTreeResp struct {
	MerkleRoot string `json:"merkleRoot"`
}

func (s *Server) CreateTree(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req createTreeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSONError(r, w, err, http.StatusBadRequest, "addresses must be a list of hex strings")
		return
	}

	if len(req.AllowedAddresses) == 0 {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "No addresses provided")
		return
	}

	addrs := make([][]byte, 0, len(req.AllowedAddresses))
	for _, addr := range req.AllowedAddresses {
		addrs = append(addrs, addr.Bytes())
	}

	tree := merkle.New(addrs)

	err := s.dbq.InsertMerkleTree(r.Context(), queries.InsertMerkleTreeParams{
		Root:      tree.Root(),
		Addresses: addrs,
	})
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to insert merkle tree")
		return
	}

	for _, addr := range addrs {
		proof := tree.Proof(addr)

		if len(proof) == 0 {
			s.sendJSONError(r, w, nil, http.StatusBadRequest, "Must provide addresses that result in a proof")
			return
		}

		err := s.dbq.InsertMerkleProof(r.Context(), queries.InsertMerkleProofParams{
			Root:    tree.Root(),
			Address: addr,
			Proof:   proof,
		})
		if err != nil {
			s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to persist merkle proofs")
			return
		}
	}

	s.sendJSON(r, w, createTreeResp{
		MerkleRoot: fmt.Sprintf("0x%s", hex.EncodeToString(tree.Root())),
	})
}

const maxAddressesPerPage = 10000

type getTreeResp struct {
	AllowedAddresses  []common.Address `json:"allowedAddresses"`
	Cursor            *string          `json:"cursor"`
	TotalAddressCount int              `json:"totalAddressCount"`
}

func (s *Server) GetTree(w http.ResponseWriter, r *http.Request) {
	rootStr := r.URL.Query().Get("root")
	if rootStr == "" {
		s.sendJSONError(r, w, nil, http.StatusBadRequest, "No merkle root provided")
		return
	}

	var root common.Hash
	if err := root.UnmarshalText([]byte(rootStr)); err != nil {
		s.sendJSONError(r, w, err, http.StatusBadRequest, "Invalid merkle root")
		return
	}

	addresses, err := s.dbq.GetAddressesForMerkleTree(r.Context(), root.Bytes())
	if errors.Is(err, pgx.ErrNoRows) {
		s.sendJSONError(r, w, err, http.StatusNotFound, "No merkle root found")
		return
	} else if err != nil {
		s.sendJSONError(r, w, err, http.StatusInternalServerError, "Failed to retrieve addresses")
		return
	}

	if len(addresses) <= maxAddressesPerPage {
		s.sendJSON(r, w, getTreeResp{
			AllowedAddresses:  addressBytesToAddresses(addresses),
			Cursor:            nil,
			TotalAddressCount: len(addresses),
		})
		return
	}

	currentCursor, err := parseCursor(r.URL.Query().Get("cursor"), len(addresses))
	if err != nil {
		s.sendJSONError(r, w, err, http.StatusBadRequest, err.Error())
		return
	}

	var nextCursor string
	nextCursorOffset := currentCursor + maxAddressesPerPage
	if nextCursorOffset < len(addresses) {
		nextCursor = strconv.Itoa(nextCursorOffset)
	}

	endOfPageIndex := currentCursor + maxAddressesPerPage
	if endOfPageIndex > len(addresses) {
		endOfPageIndex = len(addresses)
	}

	addrBytes := addresses[currentCursor:endOfPageIndex]

	s.sendJSON(r, w, getTreeResp{
		AllowedAddresses:  addressBytesToAddresses(addrBytes),
		Cursor:            &nextCursor,
		TotalAddressCount: len(addresses),
	})
}

func addressBytesToAddresses(addrBytes [][]byte) []common.Address {
	addrs := make([]common.Address, len(addrBytes))
	for i, addr := range addrBytes {
		addrs[i] = common.BytesToAddress(addr)
	}
	return addrs
}

func parseCursor(cursor string, max int) (int, error) {
	if cursor == "" {
		return 0, nil
	}
	num, err := strconv.Atoi(cursor)
	if err != nil {
		return 0, errors.New("invalid cursor")
	}

	if num < 0 {
		return 0, nil
	}

	if num >= max {
		return 0, errors.New("cursor out of range")
	}

	return num, nil
}
