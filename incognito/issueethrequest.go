package incognito

import (
	"encoding/json"
	"errors"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/rpcservice"
)

func handleCreateRawTxWithIssuingETHReq(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	arrayParams := common.InterfaceSlice(params)
	if arrayParams == nil || len(arrayParams) < 5 {
		return nil, errors.New("param must be an array at least 5 elements")
	}

	// get meta data from params
	data, ok := arrayParams[4].(map[string]interface{})
	if !ok {
		return nil, errors.New("metadata is invalid")
	}

	meta, err := metadata.NewIssuingETHRequestFromMap(data)
	if err != nil {
		return nil, err
	}

	// create new param to build raw tx from param interface
	createRawTxParam, errNewParam := bean.NewCreateRawTxParam(params)
	if errNewParam != nil {
		return nil, errNewParam
	}

	keyWallet, err := bean.GetPrivateKey(params)
	if err != nil {
		return nil, err
	}

	txService := &rpcservice.TxService{
		RpcClient: rpcClient,
		KeyWallet: keyWallet,
	}

	tx, err := txService.BuildRawTransaction(createRawTxParam, meta)
	if err != nil {
		return nil, err
	}

	byteArrays, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	result := rpcclient.CreateTransactionResult{
		TxID:            tx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func CreateAndSendTxWithIssuingETHReq(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	data, err := handleCreateRawTxWithIssuingETHReq(rpcClient, params)
	if err != nil {
		return nil, err
	}

	tx := data.(rpcclient.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData

	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	return newParam, nil

	//httpServer.handleSendRawTransaction(newParam, closeChan)
}
