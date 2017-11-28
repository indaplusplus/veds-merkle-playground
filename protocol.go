package main

type NodeJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	BlockLevel bool `json:"block"`
	//base64 the checksum
	Hash string `json:"hash"`
}

type ReadNodeJSONRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type SubNodesJSONRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

