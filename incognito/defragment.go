package incognito

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/rpcservice"
)

func createRawDefragmentAccountTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	keyWallet, err := bean.GetPrivateKey(params)
	if err != nil {
		return nil, err
	}

	txService := &rpcservice.TxService{
		RpcClient: rpcClient,
		KeyWallet: keyWallet,
	}

	tx, err := txService.BuildDeFragmentRawTransaction(params, nil)
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
		ShardID:         txShardID,
	}
	return result, nil
}

func createRawDeFragmentPTokenAccountTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	keyWallet, err := bean.GetPrivateKey(params)
	if err != nil {
		return nil, err
	}

	txService := &rpcservice.TxService{
		RpcClient: rpcClient,
		KeyWallet: keyWallet,
	}

	tx, err := txService.BuildDeFragmentPTokenRawTransaction(params, nil)
	if err != nil {
		return nil, err
	}

	byteArrays, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	txShardID := common.GetShardIDFromLastByte(tx.GetSenderAddrLastByte())
	result := rpcclient.CreateTransactionTokenResult{
		TxID:            tx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
		ShardID:         txShardID,
		TokenID:         tx.TxPrivacyTokenData.PropertyID.String(),
		TokenName:       tx.TxPrivacyTokenData.PropertyName,
		TokenAmount:     tx.TxPrivacyTokenData.Amount,
	}
	return result, nil
}

func DeFragmentAccount(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	data, err := createRawDefragmentAccountTransaction(rpcClient, params)
	if err != nil {
		return nil, err
	}

	tx := data.(rpcclient.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData

	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	return newParam, nil

	//httpServer.handleSendRawPrivacyCustomTokenTransaction(newParam, closeChan)
}

func DeFragmentPTokenAccount(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	data, err := createRawDeFragmentPTokenAccountTransaction(rpcClient, params)
	if err != nil {
		return nil, err
	}

	tx := data.(rpcclient.CreateTransactionTokenResult)
	base58CheckData := tx.Base58CheckData

	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	return newParam, nil

	//httpServer.handleSendRawPrivacyCustomTokenTransaction(newParam, closeChan)
}
