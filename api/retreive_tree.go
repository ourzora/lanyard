package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/contextart/al/api/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

const maxAddressesPerPage = 10000

type retreiveHashResponseBody struct {
	AllowedAddresses  []common.Address `json:"allowedAddresses"`
	Cursor            *string          `json:"cursor"`
	TotalAddressCount int              `json:"totalAddressCount"`
}

func (s *Server) RetrieveTree(w http.ResponseWriter, r *http.Request) {
	var (
		vars    = mux.Vars(r)
		rootStr = vars["root"]
	)
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
		s.sendJSON(r, w, retreiveHashResponseBody{
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

	var nextCursor *string
	nextCursorOffset := currentCursor + maxAddressesPerPage
	if nextCursorOffset < len(addresses) {
		nextCursor = utils.Ptr(strconv.Itoa(nextCursorOffset))
	}

	endOfPageIndex := currentCursor + maxAddressesPerPage
	if endOfPageIndex > len(addresses) {
		endOfPageIndex = len(addresses)
	}

	addrBytes := addresses[currentCursor:endOfPageIndex]

	s.sendJSON(r, w, retreiveHashResponseBody{
		AllowedAddresses:  addressBytesToAddresses(addrBytes),
		Cursor:            nextCursor,
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
