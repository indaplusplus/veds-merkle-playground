package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"crypto/sha256"
	//"crypto/sha1"
	"bytes"
	"encoding/base64"
	"crypto/x509"
	"crypto/tls"
	"encoding/json"
	"log"
)

var tree *Tree
func InitTree() {
	t := CreateTree(sha256.New())
	t.AddData([]byte{0x41, 0x42,})
	t.AddData([]byte{0x60, 0x61,})
	t.AddData([]byte{0x69, 0x69,})
	t.AddData([]byte{0x42, 0x43,})
	t.Generate()
	tree = &t
}

var client *MerkleClient

type MerkleClient struct {
	CertLocation string
}

func JSONRequest(URL string, JSON []byte) []byte {
	roots := x509.NewCertPool()
	cert_data,_ := ioutil.ReadFile(client.CertLocation)
	roots.AppendCertsFromPEM(cert_data)

	tls_conf := &tls.Config{RootCAs: roots}
	tr := &http.Transport{TLSClientConfig: tls_conf}
	http_client := &http.Client{Transport: tr}
	r, err := http_client.Post(URL, "application/json", bytes.NewBuffer(JSON))
	defer r.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	data,_ := ioutil.ReadAll(r.Body)
	return data
}

func GetRemoteRootNode() *NodeJSON{
	var root NodeJSON
	req := JSONRequest("https://localhost:1337/getroot", []byte{})
	_ = json.Unmarshal(req, &root)
	return &root
}

func GetRemoteSubNodes(parent_node *NodeJSON) *[]NodeJSON {
	var sub_nodes []NodeJSON

	tmp := SubNodesJSONRequest {
		X: parent_node.X,
		Y: parent_node.Y}

	sub_json_req,e := json.Marshal(tmp)

	if e != nil {
		log.Fatal(e)
	}

	req := JSONRequest("https://localhost:1337/getsub", sub_json_req)
	json.Unmarshal(req, &sub_nodes)
	return &sub_nodes
}

func GetRemoteNode(x, y int) *NodeJSON{
	tmp := ReadNodeJSONRequest{
		X: x,
		Y: y,
	}
	read_json_req,e := json.Marshal(tmp)
	if e != nil {
		log.Fatal(e)
	}

	req := JSONRequest("https://localhost:1337/getnode", read_json_req)
	var node NodeJSON
	json.Unmarshal(req, &node)
	return &node
}

//recursivly follow the tree until we find the missmatching block/blocks
func CompareTrees(node *NodeJSON) {
	hash,_ := base64.StdEncoding.DecodeString(node.Hash)
	local_node := tree.GetNodeAtPos(node.X, node.Y)
	if bytes.Compare(hash,local_node.Checksum) == 0 {
		fmt.Println(fmt.Sprintf("Hashes match! \n%x \n%x" ,hash, local_node.Checksum))
	} else {
		//fmt.Println("Hashes do not match")
		// https://github.com/shomali11/util looks pretty cool
		sub_nodes := GetRemoteSubNodes(node)
		tmp := *sub_nodes
		for _,l := range tmp {
			if l.BlockLevel {
				tmp := tree.GetNodeAtPos(l.X, l.Y)
				tmp_hash,_ := base64.StdEncoding.DecodeString(l.Hash)
				if tmp.Compare(tmp_hash) {
				//if bytes.Compare(tmp_hash, tmp.Checksum) != 0 {
					fmt.Println(fmt.Sprintf("Invalid block at: [%d, %d]", l.X, l.Y))
				}
			} else {
				CompareTrees(&l)
			}
		}
	}
}

func main() {
	client = &MerkleClient{"cert/server.crt"}
	InitTree()
	root := GetRemoteRootNode()
	CompareTrees(root)
}
