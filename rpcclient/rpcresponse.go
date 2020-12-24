package rpcclient

type RPCError struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
	StackTrace string `json:"StackTrace"`
}

type RPCBaseRes struct {
	Id       int       `json:"Id"`
	RPCError *RPCError    `json:"Error"`
}

type IncognitoRPCRes struct {
	RPCBaseRes
	Result interface{}
}

type ListOutputCoinsRes struct {
	RPCBaseRes
	Result *ListOutputCoins
}

type HasSerialNumberRes struct {
	RPCBaseRes
	Result []bool
}

type HasSNDerivatorRes struct {
	RPCBaseRes
	Result []bool
}

type SendRawTxRes struct {
	RPCBaseRes
	Result *CreateTransactionResult
}

type EstimateFeeRes struct {
	RPCBaseRes
	Result *EstimateFeeResult
}

type RandomCommitmentRes struct {
	RPCBaseRes
	Result *RandomCommitmentResult
}