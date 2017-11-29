package main

import (
	"net/http"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
)

/*
 * 1. Server sends hash to client
 * 2. Client checks if hashes match
 * 3. If they're the same we're done. Otherwise goto 4.
 * 4. If missmatch, client reqeuests the roots of the two subtrees.
 * 5. Server creates the necessary hashes and sends them back to the client.
 * 6. Repeat 4 and 5 until you've found the inconsistent data blocks(s).
 */

var tree *Tree = nil

type MerkleServer struct {
	Port int
	CertLocation string
	KeyLocation string
}

func GetNode(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var json_req ReadNodeJSONRequest
	e := json.Unmarshal(body, &json_req)
	if e != nil {
		log.Fatal(e)
	}

	node := tree.GetNodeAtPos(json_req.X, json_req.Y)

	json_node := ReadNodeJSONRequest{
		X: node.X,
		Y: node.Y,
	}

	j,e := json.Marshal(&json_node)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Fprintln(w, string(j))
}

func GetRootNode(w http.ResponseWriter, req *http.Request) {
	root := tree.GetRoot()
	node := NodeJSON{
		X: root.X,
		Y: root.Y,
		Hash: base64.StdEncoding.EncodeToString(root.Checksum),
		BlockLevel: (root.Data != nil),
	}
	j,e := json.Marshal(&node)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Fprintln(w, string(j))
}

func GetSubNodes(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var json_req SubNodesJSONRequest
	e := json.Unmarshal(body, &json_req)
	if e != nil {
		log.Fatal(e)
	}

	node := tree.GetNodeAtPos(json_req.X, json_req.Y)
	left := node.Left
	right := node.Right

	sub_nodes := []NodeJSON {
		{
			X: left.X,
			Y: left.Y,
			Hash: base64.StdEncoding.EncodeToString(left.Checksum),
			BlockLevel: (left.Data != nil),
		},
		{
			X: right.X,
			Y: right.Y,
			Hash: base64.StdEncoding.EncodeToString(right.Checksum),
			BlockLevel: (right.Data != nil),
		},
	}

	j,e := json.Marshal(&sub_nodes)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Fprintln(w, string(j))
}

func (server *MerkleServer) Start() {
	http.HandleFunc("/getroot", GetRootNode)
	http.HandleFunc("/getsub", GetSubNodes)
	http.HandleFunc("/getnode", GetNode)

	err := http.ListenAndServeTLS(
		fmt.Sprintf(":%d", server.Port),
		server.CertLocation,
		server.KeyLocation,
		nil)

	if err != nil {
		log.Fatal(err)
	}
}

func InitTree() {
	t := CreateTree(sha256.New())
	t.AddData([]byte{0x41, 0x42,})
	t.AddData([]byte{0x60, 0x61,})
	t.AddData([]byte{0x69, 0x69,})
	t.AddData([]byte{0x42, 0x42,})
	t.Generate()
	tree = &t
}

func main() {
	InitTree()
	server := MerkleServer{1337, "cert/server.crt", "cert/server.key"}
	server.Start()
}
