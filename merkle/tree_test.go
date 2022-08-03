package merkle

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestRoot(t *testing.T) {
	cases := []struct {
		desc     string
		leaves   [][]byte
		wantRoot []byte
	}{
		{
			leaves: [][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
				[]byte("d"),
				[]byte("e"),
				[]byte("f"),
			},
			wantRoot: common.Hex2Bytes("1b404f199ea828ec5771fb30139c222d8417a82175fefad5cd42bc3a189bd8d5"),
		},
	}

	for _, tc := range cases {
		mt := New(tc.leaves)
		if !bytes.Equal(mt.Root(), tc.wantRoot) {
			t.Errorf("got: %s want: %s",
				common.Bytes2Hex(mt.Root()),
				common.Bytes2Hex(tc.wantRoot),
			)
		}
	}
}

func TestProof(t *testing.T) {
	cases := []struct {
		leaves [][]byte
		leaf   []byte
	}{
		{
			leaves: [][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
			},
			leaf: []byte("a"),
		},
	}

	for _, tc := range cases {
		mt := New(tc.leaves)
		pf := mt.Proof(tc.leaf)
		if !Valid(mt.Root(), pf, tc.leaf) {
			t.Error("invalid proof")
		}
	}
}
