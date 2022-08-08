package api

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestAddrUnpacked(t *testing.T) {
	leaf := common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	res := addrUnpacked(leaf, []string{"uint32", "address"})
	if common.Bytes2Hex(res) != "0000000000000000000000000000000000000001" {
		t.Errorf("expected: %v got: %v", "0xaa...", common.Bytes2Hex(res))
	}
}

func TestAddrPacked(t *testing.T) {
	leaf := common.Hex2Bytes("000000000000000000000000000000000000000000000001")
	res := addrPacked(leaf, []string{"uint32", "address"})
	if common.Bytes2Hex(res) != "0000000000000000000000000000000000000001" {
		t.Errorf("expected: %v got: %v", "0xaa...", common.Bytes2Hex(res))
	}
}
