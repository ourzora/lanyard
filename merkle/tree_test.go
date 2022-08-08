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
		mt := New(tc.leaves, SortPairs, SortLeaves)
		pf := mt.Proof(tc.leaf)
		if !Valid(mt.Root(), pf, tc.leaf) {
			t.Error("invalid proof")
		}
	}
}

func hexStrArr2ByteArr(addrStrs []string) [][]byte {
	var addrs [][]byte
	for i := range addrStrs {
		addrs = append(addrs, common.FromHex(addrStrs[i]))
	}
	return addrs
}

func TestProofFiveAddresses(t *testing.T) {
	addrStrs := []string{
		"0x0000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000002",
		"0x0000000000000000000000000000000000000003",
		"0x0000000000000000000000000000000000000004",
		"0x0000000000000000000000000000000000000005",
	}
	addrs := hexStrArr2ByteArr(addrStrs)

	tree := New(addrs, SortPairs)

	proofs := [][]string{
		{
			"d52688a8f926c816ca1e079067caba944f158e764817b83fc43594370ca9cf62",
			"735c77c52a2b69afcd4e13c0a6ece7e4ccdf2b379d39417e21efe8cd10b5ff1b",
			"421df1fa259221d02aa4956eb0d35ace318ca24c0a33a64c1af96cf67cf245b6",
		},
		{
			"1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d",
			"735c77c52a2b69afcd4e13c0a6ece7e4ccdf2b379d39417e21efe8cd10b5ff1b",
			"421df1fa259221d02aa4956eb0d35ace318ca24c0a33a64c1af96cf67cf245b6",
		},
		{
			"a876da518a393dbd067dc72abfa08d475ed6447fca96d92ec3f9e7eba503ca61",
			"f95c14e6953c95195639e8266ab1a6850864d59a829da9f9b13602ee522f672b",
			"421df1fa259221d02aa4956eb0d35ace318ca24c0a33a64c1af96cf67cf245b6",
		},
		{
			"5b70e80538acdabd6137353b0f9d8d149f4dba91e8be2e7946e409bfdbe685b9",
			"f95c14e6953c95195639e8266ab1a6850864d59a829da9f9b13602ee522f672b",
			"421df1fa259221d02aa4956eb0d35ace318ca24c0a33a64c1af96cf67cf245b6",
		},
		{
			"5071e19149cc9b870c816e671bc5db717d1d99185c17b082af957a0a93888dd9",
		},
	}

	for i, addr := range addrs {
		actual := tree.Proof(addr)
		expected := hexStrArr2ByteArr(proofs[i])
		if len(actual) != len(expected) {
			t.Errorf("got: %d want: %d", len(actual), len(expected))
		}
		for j := range actual {
			if !bytes.Equal(actual[j], expected[j]) {
				t.Errorf("got: %s want: %s",
					common.Bytes2Hex(actual[j]),
					common.Bytes2Hex(expected[j]),
				)
			}
		}
	}
}
