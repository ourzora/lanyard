package merkle

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func ExampleTree() {
	var addrStrs = []string{
		"0xE124F06277b5AC791bA45B92853BA9A0ea93327D",
		"0x07d048f78B7C093B3Ef27D478B78026a70D9734e",
		"0x38976611f5f7bEAd7e79E752f5B80AE72dD3eFa7",
		"0x1Ab00ffedD724B930080aD30269083F1453cF34E",
		"0x860a6bC426C3bb1186b2E11Ac486ABa000C209B4",
		"0x0B3eC21fc53AD8b17AF4A80723c1496541fCb35f",
		"0x2D13F6CEe6dA8b30a84ee7954594925bd5E47Ab7",
		"0x3C64Cd43331beb5B6fAb76dbAb85226955c5CC3A",
		"0x238dA873f984188b4F4c7efF03B5580C65a49dcB",
		"0xbAfC038aDfd8BcF6E632C797175A057714416d04",
	}
	var addrs [][]byte
	for i := range addrStrs {
		addrs = append(addrs, common.HexToAddress(addrStrs[i]).Bytes())
	}
	var tr Tree
	tr = New(addrs)
	fmt.Println(common.Bytes2Hex(tr.Root()))

	tr = New(addrs, SortPairs)
	fmt.Println(common.Bytes2Hex(tr.Root()))

	tr = New(addrs, SortPairs, SortLeaves)
	fmt.Println(common.Bytes2Hex(tr.Root()))

	// Output:
	// cc32560f07705feff369dc58b2d65438e59395399f1e4c2b19758f3a752db050
	// ed40d49077a2cd13601cf79a512e6b92c7fd0f952e7dc9f4758d7134f9712bc4
	// 51c19eb453389ea5e10f221da238d321cde4f8979f65bde8445581320620bed5
}

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
		mt := New(tc.leaves, SortPairs, SortLeaves)
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
	}{
		{
			leaves: [][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
				[]byte("d"),
				[]byte("e"),
			},
		},
	}

	for _, tc := range cases {
		mt := New(tc.leaves, SortPairs, SortLeaves)
		for _, l := range tc.leaves {
			pf := mt.Proof(l)
			if !Valid(mt.Root(), pf, l) {
				t.Error("invalid proof")
			}
		}
	}
}
