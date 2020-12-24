package service

type Parameter struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      uint        `json:"id"`
}

type Payload struct {
	Method     string   `json:"method"`
	PrivateKey string   `json:"private_key"`
	Data       interface{} `json:"data"`
}

type Response struct {
	Result interface{} `json:"Result"`
	Error  string      `json:"Error"`
	Stack  string      `json:"Stack"`
	Method string      `json:"Method"`
}