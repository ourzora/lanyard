//go:build integration

package lanyard

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	client      *Client
	basicMerkle = []hexutil.Bytes{
		hexutil.MustDecode("0x0000000000000000000000000000000000000001"),
		hexutil.MustDecode("0x0000000000000000000000000000000000000002"),
		hexutil.MustDecode("0x0000000000000000000000000000000000000003"),
		hexutil.MustDecode("0x0000000000000000000000000000000000000004"),
		hexutil.MustDecode("0x0000000000000000000000000000000000000005"),
	}
)

func init() {
	if os.Getenv("LANYARD_API_BASE_URL") == "" {
		client = New()
	} else {
		client = New(WithURL(os.Getenv("LANYARD_API_BASE_URL")))
	}
}

const (
	basicRoot = "0xa7a6b1cb6d12308ec4818baac3413fafa9e8b52cdcd79252fa9e29c9a2f8aff1"
	typedRoot = "0x6306f03ad6ae2ffeca080333a0a6828669192f5f8b61f70738bfe8ceb7e0a434"
)

func TestBasicMerkleTree(t *testing.T) {
	tree, err := client.CreateTree(context.Background(), basicMerkle)
	if err != nil {
		t.Fatal(err)
	}

	if tree.MerkleRoot.String() != basicRoot {
		t.Fatalf("expected %s, got %s", basicRoot, tree.MerkleRoot.String())
	}
}

func TestCreateTypedTree(t *testing.T) {
	tree, err := client.CreateTypedTree(
		context.Background(),
		[]hexutil.Bytes{
			hexutil.MustDecode("0x00000000000000000000000000000000000000010000000000000000000000000000000000000000000000008ac7230489e80000"),
			hexutil.MustDecode("0x0000000000000000000000000000000000000002000000000000000000000000000000000000000000000001e5b8fa8fe2ac0000"),
		},
		[]string{"address", "uint256"},
		true,
	)

	if err != nil {
		t.Fatal(err)
	}

	if tree.MerkleRoot.String() != typedRoot {
		t.Fatalf("expected %s, got %s", typedRoot, tree.MerkleRoot.String())
	}
}

func TestBasicMerkleProof(t *testing.T) {
	_, err := client.GetProofFromLeaf(context.Background(), hexutil.MustDecode(basicRoot), basicMerkle[0])
	if err != nil {
		t.Fatal(err)
	}
}

func TestBasicMerkleProof404(t *testing.T) {
	_, err := client.GetProofFromLeaf(context.Background(), []byte{0x01}, hexutil.MustDecode("0x0000000000000000000000000000000000000001"))
	if err != ErrNotFound {
		t.Fatal("expected custom 404 err type for invalid request, got %w", err)
	}
}

func TestGetProofFromAddr(t *testing.T) {
	_, err := client.GetProofFromAddr(context.Background(), hexutil.MustDecode(typedRoot), hexutil.MustDecode("0x0000000000000000000000000000000000000001"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRootFromProof(t *testing.T) {
	p, err := client.GetProofFromLeaf(context.Background(), hexutil.MustDecode(basicRoot), basicMerkle[0])

	if err != nil {
		t.Fatal(err)
	}

	root, err := client.GetRootFromProof(context.Background(), p.Proof)

	if err != nil {
		t.Fatal(err)
	}

	if root.Root.String() != basicRoot {
		t.Fatalf("expected %s, got %s", basicRoot, root.Root.String())
	}

}

func TestGetTree(t *testing.T) {
	tree, err := client.GetTreeFromRoot(context.Background(), hexutil.MustDecode(basicRoot))

	if err != nil {
		t.Fatal(err)
	}

	if tree.UnhashedLeaves[0].String() != basicMerkle[0].String() {
		t.Fatalf("expected %s, got %s", basicMerkle[0].String(), tree.UnhashedLeaves[0].String())
	}

	if tree.LeafCount != len(basicMerkle) {
		t.Fatalf("expected %d, got %d", len(basicMerkle), tree.LeafCount)
	}
}
