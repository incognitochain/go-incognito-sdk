package incognito

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/rpcservice"
)

func handleCreateRawTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
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

	tx, err := txService.BuildRawTransaction(createRawTxParam, nil)
	if err != nil {
		return nil, err
	}

	byteArrays, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	txShardID := common.GetShardIDFromLastByte(tx.GetSenderAddrLastByte())

	result := rpcclient.CreateTransactionResult{
		TxID:            tx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
		ShardID: txShardID,
	}
	return result, nil
}

func CreateAndSendTx(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	var err error
	data, err := handleCreateRawTransaction(rpcClient, params)
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
