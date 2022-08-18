package api

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestAddrUnpacked(t *testing.T) {
	cases := []struct {
		leaf []byte
		ltd  []string
		want common.Address
	}{
		{
			common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"),
			[]string{"uint32", "address"},
			common.HexToAddress("0x0000000000000000000000000000000000000001"),
		},
		{
			common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000"),
			[]string{"address", "uint32"},
			common.HexToAddress("0x0000000000000000000000000000000000000001"),
		},
	}

	for _, c := range cases {
		addr := addrUnpacked(c.leaf, c.ltd)
		if addr != c.want {
			t.Errorf("expected: %v got: %v", c.want, addr)
		}
	}
}

func TestAddrPacked(t *testing.T) {
	cases := []struct {
		leaf []byte
		ltd  []string
		want common.Address
	}{
		{
			common.Hex2Bytes("000000000000000000000000000000000000000000000001"),
			[]string{"uint32", "address"},
			common.HexToAddress("0x0000000000000000000000000000000000000001"),
		},
		{
			common.Hex2Bytes("000000000000000000000000000000000000000100000000"),
			[]string{"address", "uint32"},
			common.HexToAddress("0x0000000000000000000000000000000000000001"),
		},
	}

	for _, c := range cases {
		addr := addrPacked(c.leaf, c.ltd)
		if addr != c.want {
			t.Errorf("expected: %v got: %v", c.want, addr)
		}
	}
}
