package merkle

import (
	"crypto/sha256"
)

type Tree struct {
	Root *Node
}

type Node struct {
	Left, Right *Node
	Hash        *[32]byte
}

func NewTree(data [][]byte) *Tree {
	cnt := len(data)
	if cnt <= 0 {
		return nil
	}
	nodes := make([]*Node, cnt)
	for i := range data {
		nodes[i] = NewNode(nil, nil, data[i])
	}
	nodes = buildLevel(nodes, make([]byte, 64))
	return &Tree{Root: nodes[0]}
}

func buildLevel(nodes []*Node, buff []byte) []*Node {
	cnt := len(nodes)
	if cnt <= 1 {
		return nodes
	}
	newCnt := cnt / 2
	for i := 0; i < newCnt; i++ {
		buff, left, right := buff[:0], nodes[i*2], nodes[i*2+1]
		buff = append(buff, left.Hash[:]...)
		buff = append(buff, right.Hash[:]...)
		nodes[i] = NewNode(left, right, buff)
	}
	if newCnt*2 < cnt {
		buff, left := buff[:0], nodes[cnt-1]
		buff = append(buff, left.Hash[:]...)
		nodes[newCnt] = NewNode(left, nil, buff)
		newCnt++
	}
	return buildLevel(nodes[:newCnt], buff)
}

func NewNode(left, right *Node, data []byte) *Node {
	hash := sha256.Sum256(data)
	return &Node{Left: left, Right: right, Hash: &hash}
}
