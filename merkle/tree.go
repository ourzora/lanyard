package merkle

import (
	"bytes"
	"sort"

	"github.com/ethereum/go-ethereum/crypto"
)

type option int

const (
	sortLeaves option = iota
	sortPairs
)

func has(opts []option, o option) bool {
	for i := range opts {
		if opts[i] == o {
			return true
		}
	}
	return false
}

type Tree [][][]byte

func New(items [][]byte, opts ...option) Tree {
	var leaves [][]byte
	for i := range items {
		leaves = append(leaves, crypto.Keccak256(items[i]))
	}
	if has(opts, sortLeaves) {
		sort.Slice(leaves, func(i, j int) bool {
			return bytes.Compare(leaves[i], leaves[j]) == -1
		})
	}
	var t Tree
	t = append(t, leaves)

	for {
		level := t[len(t)-1]
		if len(level) == 1 {
			break
		}
		f := hashPairNoSort
		if has(opts, sortPairs) {
			f = hashPairSort
		}
		t = append(t, hashMerge(level, f))
	}
	return t
}

type hashPair func(a, b []byte) []byte

func hashPairNoSort(a, b []byte) []byte {
	return crypto.Keccak256(append(a, b...))
}

func hashPairSort(a, b []byte) []byte {
	if bytes.Compare(a, b) == -1 {
		return crypto.Keccak256(append(a, b...))
	}
	return crypto.Keccak256(append(b, a...))
}

func hashMerge(level [][]byte, f hashPair) [][]byte {
	var newLevel [][]byte
	for i := 0; i < len(level); i += 2 {
		switch {
		case i+1 == len(level):
			newLevel = append(newLevel, level[i])
		default:
			newLevel = append(newLevel, f(level[i], level[i+1]))
		}
	}
	return newLevel
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
		target = hashPairSort(target, proof[i])
	}
	return bytes.Compare(target, root) == 0
}
