package api

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestAddrUnpacked(t *testing.T) {
	cases := []struct {
		leaf []byte
		ltd  []string
		want string
	}{
		{
			common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"),
			[]string{"uint32", "address"},
			"0x0000000000000000000000000000000000000001",
		},
		{
			common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000"),
			[]string{"address", "uint32"},
			"0x0000000000000000000000000000000000000001",
		},
	}

	for _, c := range cases {
		res := addrUnpacked(c.leaf, c.ltd)
		if hexutil.Encode(res) != c.want {
			t.Errorf("expected: %v got: %v", "0xaa...", hexutil.Encode(res))
		}
	}
}

func TestAddrPacked(t *testing.T) {
	cases := []struct {
		leaf []byte
		ltd  []string
		want string
	}{
		{
			common.Hex2Bytes("000000000000000000000000000000000000000000000001"),
			[]string{"uint32", "address"},
			"0x0000000000000000000000000000000000000001",
		},
		{
			common.Hex2Bytes("000000000000000000000000000000000000000100000000"),
			[]string{"address", "uint32"},
			"0x0000000000000000000000000000000000000001",
		},
	}

	for _, c := range cases {
		res := addrPacked(c.leaf, c.ltd)
		if hexutil.Encode(res) != c.want {
			t.Errorf("expected: %v got: %v", c.want, hexutil.Encode(res))
		}
	}
}
