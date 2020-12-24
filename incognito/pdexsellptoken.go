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

func handleCreateRawTxWithPTokenTradeReq(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	arrayParams := common.InterfaceSlice(params)

	if len(arrayParams) >= 7 {
		hasPrivacyToken := int(arrayParams[6].(int)) > 0
		if hasPrivacyToken {
			return nil, errors.New("The privacy mode must be disabled")
		}
	}
	tokenParamsRaw := arrayParams[4].(map[string]interface{})

	tokenIDToBuyStr, ok := tokenParamsRaw["TokenIDToBuyStr"].(string)
	if !ok {
		return nil, errors.New("TokenIDToBuyStr is invalid")
	}

	tokenIDToSellStr, ok := tokenParamsRaw["TokenIDToSellStr"].(string)
	if !ok {
		return nil, errors.New("TokenIDToSellStr is invalid")
	}

	sellAmountData, ok := tokenParamsRaw["SellAmount"].(uint64)
	if !ok {
		return nil, errors.New("SellAmount is invalid")
	}
	sellAmount := uint64(sellAmountData)

	traderAddressStr, ok := tokenParamsRaw["TraderAddressStr"].(string)
	if !ok {
		return nil, errors.New("TraderAddressStr is invalid")
	}

	minAcceptableAmountData, ok := tokenParamsRaw["MinAcceptableAmount"].(uint64)
	if !ok {
		return nil, errors.New("MinAcceptableAmount is invalid")
	}
	minAcceptableAmount := uint64(minAcceptableAmountData)

	tradingFeeData, ok := tokenParamsRaw["TradingFee"].(uint64)
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

	customTokenTx, rpcErr := txService.BuildRawPrivacyCustomTokenTransaction(params, meta)
	if rpcErr != nil {
		return nil, rpcErr
	}

	byteArrays, err2 := json.Marshal(customTokenTx)
	if err2 != nil {
		return nil, err2
	}

	result := rpcclient.CreateTransactionResult{
		TxID:            customTokenTx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}

	return result, nil
}

func CreateAndSendTxWithPTokenTradeReq(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	data, err := handleCreateRawTxWithPTokenTradeReq(rpcClient, params)
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