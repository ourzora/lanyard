package merkle

import (
	"bytes"
	"sort"

	"github.com/ethereum/go-ethereum/crypto"
)

type Tree [][][]byte

func hashPair(a, b []byte) []byte {
	if bytes.Compare(a, b) == -1 {
		return crypto.Keccak256(append(a, b...))
	}
	return crypto.Keccak256(append(b, a...))
}

func hashMerge(level [][]byte) [][]byte {
	var newLevel [][]byte
	for i := 0; i < len(level); i += 2 {
		switch {
		case i+1 == len(level):
			newLevel = append(newLevel, level[i])
		default:
			newLevel = append(newLevel, hashPair(level[i], level[i+1]))
		}
	}
	return newLevel
}

func New(items [][]byte) Tree {
	var leaves [][]byte
	for i := range items {
		leaves = append(leaves, crypto.Keccak256(items[i]))
	}
	sort.Slice(leaves, func(i, j int) bool {
		return bytes.Compare(leaves[i], leaves[j]) == -1
	})

	var t Tree
	t = append(t, leaves)

	for {
		level := t[len(t)-1]
		if len(level) == 1 {
			break
		}
		t = append(t, hashMerge(level))
	}
	return t
}

func (t Tree) Root() []byte {
	return t[len(t)-1][0]
}

func (t Tree) Proof(target []byte) [][]byte {
	var (
		proof [][]byte
		index int
	)
	for i, h := range t[0] {
		if bytes.Equal(crypto.Keccak256(target), h) {
			index = i
		}
	}
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

func Valid(root []byte, proof [][]byte, target []byte) bool {
	target = crypto.Keccak256(target)
	for i := range proof {
		target = hashPair(target, proof[i])
	}
	return bytes.Compare(target, root) == 0
}
