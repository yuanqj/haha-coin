package merkle_test

import (
	"fmt"
	"github.com/yuanqj/haha-coin/merkle"
	"testing"
)

func walkLevels(root *merkle.Node) []*merkle.Node {
	nodes := make([]*merkle.Node, 0, 100)
	q := make(chan *merkle.Node, 100)
	q <- root
	for len(q) > 0 {
		node := <-q
		if node == nil {
			continue
		}
		q <- node.Left
		q <- node.Right
		nodes = append(nodes, node)
	}
	return nodes
}

func TestTree(t *testing.T) {
	data := [][]byte{
		[]byte("0"),
		[]byte("1"),
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
	}
	hashes := []string{
		"d1f0ede7376054cac733014de7fb04f0584e12b45c2f1a8df40ad594548a0c8c",

		"c478fead0c89b79540638f844c8819d9a4281763af9272c7f3968776b6052345",
		"72747334a2dbe8e874eb297eb50e0036153517dc43c5c1301523793f2e0f503f",

		"b9b10a1bc77d2a241d120324db7f3b81b2edb67eb8e9cf02af9c95d30329aef5",
		"a9f5b3ab61e28357cfcd14e2b42397f896aeea8d6998d19e6da85584e150d2b4",
		"5a9eab9148389395eff050ddf00220d722123ca8736c862bf200316389b3f611",

		"5feceb66ffc86f38d952786c6d696c79c2dbc239dd4e91b46729d73a27fb57e9",
		"6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
		"d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35",
		"4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce",
		"4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a",
	}

	tree := merkle.NewTree(data)
	nodes := walkLevels(tree.Root)
	if len(nodes) != 11 {
		t.Errorf("node count mismatches: expect=11, exact=%d", len(nodes))
		return
	}
	for i := range nodes {
		hash := fmt.Sprintf("%x", (nodes[i].Hash)[:])
		if hashes[i] != hash {
			t.Errorf("decode error: #%02d: expected=%s, exact=%s", i, hashes[i], hash)
		}
	}
}
