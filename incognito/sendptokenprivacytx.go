package incognito

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/rpcservice"
)

func handleCreateRawPrivacyCustomTokenTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	keyWallet, err := bean.GetPrivateKey(params)
	if err != nil {
		return nil, err
	}

	txService := &rpcservice.TxService{
		RpcClient: rpcClient,
		KeyWallet: keyWallet,
	}

	tx, err := txService.BuildRawPrivacyCustomTokenTransaction(params, nil)
	if err != nil {
		return nil, err
	}

	byteArrays, err := json.Marshal(tx)
	if err != nil {
		return nil,  err
	}

	result := rpcclient.CreateTransactionTokenResult{
		ShardID:         common.GetShardIDFromLastByte(tx.Tx.PubKeyLastByteSender),
		TxID:            tx.Hash().String(),
		TokenID:         tx.TxPrivacyTokenData.PropertyID.String(),
		TokenName:       tx.TxPrivacyTokenData.PropertyName,
		TokenAmount:     tx.TxPrivacyTokenData.Amount,
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func CreateAndSendPrivacyCustomTokenTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	var err error
	data, err := handleCreateRawPrivacyCustomTokenTransaction(rpcClient, params)
	if err != nil {
		return nil, err
	}

	tx := data.(rpcclient.CreateTransactionTokenResult)
	base58CheckData := tx.Base58CheckData
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	return newParam, nil
	//txId, err := httpServer.handleSendRawPrivacyCustomTokenTransaction(newParam, closeChan)
}

