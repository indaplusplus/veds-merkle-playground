package main

import (
	"hash"
	"fmt"
	"bytes"
	"log"
)

type Node struct {
	Checksum []byte
	Data []byte
	X,Y int
	Parent *Node
	Left *Node
	Right *Node
}

type Tree struct {
	Nodes [][]Node
	HashFunc hash.Hash
}

// Gets the root node on the top of the tree. 
func (tree *Tree) GetRoot() Node {
	return tree.Nodes[len(tree.Nodes) - 1][0]
}

func (tree *Tree) PrintTree() {
	for i, nodes := range tree.Nodes {
		for n, node := range nodes {
			fmt.Println(fmt.Sprintf("[%d, %d] %x", i, n, node.Checksum))
		}
	}
}

func Lookup(node *Node) (*Node, *Node) {
	return node.Left, node.Right
}

//compare a nodes hash with a given
func (node *Node) Compare(hash []byte) (bool) {
	// https://golang.org/pkg/bytes/#Compare
	return bytes.Compare(node.Checksum, hash) == 0
}

func (tree *Tree) GetNodeAtPos(x,y int) *Node {
	if ((len(tree.Nodes) - 1) < y ||
		(len(tree.Nodes[y]) - 1) < x) {
		log.Fatal("Out of bound read")
	}
	return &tree.Nodes[y][x]
}

// Generates parents nodes until root node
func (tree *Tree) Generate() {
	// bottom level nodes have to be a multiple of two.
	if (len(tree.Nodes[0]) % 2 != 0) {
		panic("Not a multiple of two!")
	}

	level := 0
	for len(tree.Nodes[level]) != 1 {
		// Make more space in the array per level
		tree.Nodes = append(tree.Nodes, []Node{})
		l := len(tree.Nodes[level]) - 1
		for i := 0; i < l; i++ {
			node_one := tree.Nodes[level][i]
			node_two := tree.Nodes[level][i + 1]
			tree.HashFunc.Write(append(node_one.Checksum,node_two.Checksum...))
			hash := tree.HashFunc.Sum(nil)
			tree.HashFunc.Reset()

			parent_node := Node{
				Checksum: hash,
				X: len(tree.Nodes[level+1]),
				Y: (level+1),
			}
			//setup pointers
			parent_node.Left = &node_one
			parent_node.Right = &node_two
			tree.Nodes[level+1] = append(tree.Nodes[level+1], parent_node)
			//set parent pointer
			node_one.Parent = &parent_node
			node_two.Parent = &parent_node
		}
		level++
	}
}

func (tree *Tree) AddData(data []byte) {
	tree.HashFunc.Write(data)
	hash := tree.HashFunc.Sum(nil)
	tree.HashFunc.Reset()

	node := Node{
		Checksum: hash,
		Data: data,
		X: len(tree.Nodes[0]),
		Y: 0,
	}
	tree.Nodes[0] = append(tree.Nodes[0], node)
}

func CreateTree(hashfunc hash.Hash) Tree {
	return Tree {
		Nodes: [][]Node {{}},
		HashFunc: hashfunc,
	}
}
