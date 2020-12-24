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

func handleCreateRawTxWithPRVTradeReq(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	arrayParams := common.InterfaceSlice(params)

	// get meta data from params
	data, ok := arrayParams[4].(map[string]interface{})
	if !ok {
		return nil, errors.New("metadata is invalid")
	}
	tokenIDToBuyStr, ok := data["TokenIDToBuyStr"].(string)
	if !ok {
		return nil, errors.New("TokenIDToBuyStr is invalid")
	}
	tokenIDToSellStr, ok := data["TokenIDToSellStr"].(string)
	if !ok {
		return nil,errors.New("TokenIDToSellStr is invalid")
	}
	sellAmount := uint64(data["SellAmount"].(uint64))
	traderAddressStr, ok := data["TraderAddressStr"].(string)
	if !ok {
		return nil, errors.New("TraderAddressStr is invalid")
	}
	minAcceptableAmountData, ok := data["MinAcceptableAmount"].(uint64)
	if !ok {
		return nil, errors.New("MinAcceptableAmount is invalid")
	}
	minAcceptableAmount := uint64(minAcceptableAmountData)
	tradingFeeData, ok := data["TradingFee"].(uint64)
	if !ok {
		return nil, errors.New("TradingFee is invalid")
	}
	tradingFee := uint64(tradingFeeData)
	meta, _ := metadata.NewPDETradeRequest(
		tokenIDToBuyStr,
		tokenIDToSellStr,
		sellAmount,
		minAcceptableAmount,
		tradingFee,
		traderAddressStr,
		metadata.PDETradeRequestMeta,
	)

	keyWallet, err := bean.GetPrivateKey(params)
	if err != nil {
		return nil, err
	}

	txService := &rpcservice.TxService{
		RpcClient: rpcClient,
		KeyWallet: keyWallet,
	}

	// create new param to build raw tx from param interface
	createRawTxParam, errNewParam := bean.NewCreateRawTxParam(params)
	if errNewParam != nil {
		return nil, errNewParam
	}

	tx, err1 := txService.BuildRawTransaction(createRawTxParam, meta)
	if err1 != nil {
		return nil, err1
	}

	byteArrays, err2 := json.Marshal(tx)
	if err2 != nil {
		return nil, err2
	}
	result := rpcclient.CreateTransactionResult{
		TxID:            tx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

func CreateAndSendTxWithPRVTradeReq(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	data, err := handleCreateRawTxWithPRVTradeReq(rpcClient, params)
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