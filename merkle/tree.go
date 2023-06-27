// A merkle tree for [lanyard.org].
// merkle uses keccak256 hashing for leaves
// and intermediary nodes and therefore is vulnerable to a
// second preimage attack. This package does not duplicate or pad leaves in
// the case of odd cardinality and therefore may be unsafe for certain
// use cases. If you are curious about this type of bug,
// see the following [bitcoin issue].
//
// [bitcoin issue]: https://bitcointalk.org/index.php?topic=102395.0
// [lanyard.org]: https://lanyard.org
package merkle

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
)

// A the outer list represents levels in the tree. Each level is a list
// of nodes in the tree. Each node is a hash of its children.
type Tree [][][]byte

// Returns a complete Tree using items for the leaves.
// Intermediary nodes and items will be hashed using Keccak256.
func New(items [][]byte) Tree {
	var leaves [][]byte
	for i := range items {
		leaves = append(leaves, crypto.Keccak256(items[i]))
	}
	var t Tree
	t = append(t, leaves)

	for {
		level := t[len(t)-1]
		if len(level) == 1 { //root node
			break
		}
		t = append(t, hashMerge(level))
	}
	return t
}

func hashPair(a, b []byte) []byte {
	if bytes.Compare(a, b) == -1 { // a < b
		return crypto.Keccak256(a, b)
	}
	return crypto.Keccak256(b, a)
}

// Iterates through the level pairwise merging each
// pair with a hash function creating a new level that
// is half the size of the level.
func hashMerge(level [][]byte) [][]byte {
	var newLevel [][]byte
	for i := 0; i < len(level); i += 2 {
		switch {
		case i+1 == len(level):
			// In the case of a level with an odd number of nodes
			// we leave the parent with a single child.
			// Some merkle tree designs allow for the parent
			// to duplicate the child so that it has both children
			// thus leaving the level with an even number of nodes.
			// We don't have that requirement yet and if one day we do
			// this is the spot to change:
			newLevel = append(newLevel, level[i])
		default:
			newLevel = append(newLevel, hashPair(level[i], level[i+1]))
		}
	}
	return newLevel
}

func (t Tree) Root() []byte {
	return t[len(t)-1][0]
}

// Returns a list of hashes such that
// cumulatively hashing the list pairwise
// will yield the root hash of the tree. Example:
//
//	[abcde]
//	[abcd, e]
//	[ab, cd, e]
//	[a, b, c, d, e]
//
// If the target is 'c' Proof returns:
// [d, ab, e]
//
// If the target is 'e' Proof returns:
// [cd]
//
// The result of this func will be used in [Valid]
func (t Tree) Proof(target []byte) [][]byte {
	var (
		ht    = crypto.Keccak256(target)
		index int
	)
	for i, h := range t[0] {
		if bytes.Equal(ht, h) {
			index = i
			break
		}
	}

	return t.proofForEdge(index)
}

func (t Tree) proofForEdge(index int) [][]byte {
	var proof [][]byte
	for _, level := range t {
		var i int
		switch {
		case index%2 == 0:
			i = index + 1
		case index%2 == 1:
			i = index - 1
		}
		if i < len(level) {
			proof = append(proof, level[i])
		}
		index = index / 2
	}
	return proof
}

// Returns proofs for all edges in the tree.
// For details on how an individual proof is calculated, see [Tree.Proof].
func (t Tree) Proofs() [][][]byte {
	var (
		proofs = make([][][]byte, 0, len(t[0]))
	)

	for i := range t[0] {
		proofs = append(proofs, t.proofForEdge(i))
	}

	return proofs
}

// Cumulatively hashes the list pairwise starting with
// (target, proof[0]). Finally, the cumulative hash is compared with the root.
func Valid(root []byte, proof [][]byte, target []byte) bool {
	target = crypto.Keccak256(target)
	for i := range proof {
		target = hashPair(target, proof[i])
	}
	return bytes.Equal(target, root)
}
