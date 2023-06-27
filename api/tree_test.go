package api

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestAddrUnpacked(t *testing.T) {
	cases := []struct {
		leaf []byte
		ltd  []string
		want []byte
	}{
		{
			common.FromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"),
			[]string{"uint32", "address"},
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000"),
			[]string{"address", "uint32"},
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
	}

	for _, c := range cases {
		addr := addrUnpacked(c.leaf, c.ltd)
		if !bytes.Equal(addr, c.want) {
			t.Errorf("expected: %v got: %v", c.want, addr)
		}
	}
}

func TestAddrPacked(t *testing.T) {
	cases := []struct {
		leaf []byte
		ltd  []string
		want []byte
	}{
		{
			common.FromHex("000000000000000000000000000000000000000000000001"),
			[]string{"uint32", "address"},
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("000000000000000000000000000000000000000100000000"),
			[]string{"address", "uint32"},
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
	}

	for _, c := range cases {
		addr := addrPacked(c.leaf, c.ltd)
		if !bytes.Equal(addr, c.want) {
			t.Errorf("expected: %v got: %v", c.want, addr)
		}
	}
}

func TestLeaf2Addr(t *testing.T) {
	cases := []struct {
		leaf   []byte
		ltd    []string
		packed bool
		want   []byte
	}{
		{
			common.FromHex("000000000000000000000000000000000000000000000001"),
			[]string{"uint32", "address"},
			true,
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("000000000000000000000000000000000000000100000000"),
			[]string{"address", "uint32"},
			true,
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("0x0000000000000000000000000000000000000001"),
			[]string{"address"},
			false,
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"),
			[]string{"uint32", "address"},
			false,
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000"),
			[]string{"address", "uint32"},
			false,
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"),
			[]string{"uint256", "address"},
			false,
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002d"),
			[]string{"address", "uint256"},
			false,
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
		{
			common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001"),
			[]string{"address"},
			true,
			common.FromHex("0x0000000000000000000000000000000000000001"),
		},
	}

	for _, c := range cases {
		addr := leaf2Addr(c.leaf, c.ltd, c.packed)
		if !bytes.Equal(addr, c.want) {
			t.Errorf("expected: %v got: %v", c.want, addr)
		}
	}

}
