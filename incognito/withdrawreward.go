package incognito

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/rpcservice"
	"github.com/pkg/errors"
)

func handleCreateRawWithDrawTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	arrayParams := common.InterfaceSlice(params)
	if arrayParams == nil || len(arrayParams) < 5 {
		return nil, errors.New("param must be an array at least 5 elements")
	}
	arrayParams[1] = nil

	keyWallet, err := bean.GetPrivateKey(params)
	if err != nil {
		return nil, err
	}

	metaParam, ok := arrayParams[4].(map[string]interface{})
	if !ok {
		return nil, errors.New("metadata is invalid")
	}
	tokenIDParam, ok := metaParam["TokenID"]
	if !ok {
		return nil, errors.New("token ID is invalid")
	}

	//refactor param
	param := map[string]interface{}{}
	param["PaymentAddress"] = keyWallet.Base58CheckSerialize(1)
	param["TokenID"] = tokenIDParam
	param["Version"] = 1
	if version, ok := metaParam["Version"]; ok {
		param["Version"] = version
	}
	arrayParams[4] = interface{}(param)

	// param #5 get meta data param
	metaRaw, ok := arrayParams[4].(map[string]interface{})
	if !ok {
		return nil, errors.New("metadata param is invalid")
	}

	meta, errCons := metadata.NewWithDrawRewardRequestFromRPC(metaRaw)
	if errCons != nil {
		return nil, err
	}

	// create new param to build raw tx from param interface
	createRawTxParam, errNewParam := bean.NewCreateRawTxParam(params)
	if errNewParam != nil {
		return nil, errNewParam
	}

	txService := &rpcservice.TxService{
		RpcClient: rpcClient,
		KeyWallet: keyWallet,
	}

	tx, err := txService.BuildRawTransaction(createRawTxParam, meta)
	if err != nil {
		return nil, err
	}

	byteArrays, errMarshal := json.Marshal(tx)
	if errMarshal != nil {
		return nil, errMarshal
	}

	result := rpcclient.CreateTransactionResult{
		TxID:            tx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func CreateAndSendWithDrawTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	var err error
	data, err := handleCreateRawWithDrawTransaction(rpcClient, params)
	if err != nil {
		return nil, err
	}

	tx := data.(rpcclient.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	return newParam, nil
	//sendResult, err := httpServer.handleSendRawTransaction(newParam, closeChan)
}
