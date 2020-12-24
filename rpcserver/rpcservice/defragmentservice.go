package rpcservice

import (
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/incognitokey"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/transaction"
)

func (txService TxService) BuildDeFragmentRawTransaction(
	params interface{},
	metadataParam metadata.Metadata,
) (*transaction.Tx, error) {
	arrayParams := common.InterfaceSlice(params)
	if len(arrayParams) < 4 {
		return nil, errors.New("Params is invalid")
	}

	senderKeyParam, ok := arrayParams[0].(string)
	if !ok {
		return nil, errors.New("senderKeyParam is invalid")
	}

	maxValTemp, ok := arrayParams[1].(int64)
	if !ok {
		return nil, errors.New("maxVal is invalid")
	}
	maxVal := uint64(maxValTemp)

	estimateFeeCoinPerKbtemp, ok := arrayParams[2].(int)
	if !ok {
		return nil, errors.New("estimateFeeCoinPerKb is invalid")
	}
	estimateFeeCoinPerKb := int64(estimateFeeCoinPerKbtemp)

	// param #4: hasPrivacyCoin flag: 1 or -1
	hasPrivacyCoinParam := arrayParams[3].(int)
	hasPrivacyCoin := hasPrivacyCoinParam > 0

	maxDefragmentQuantity := 32
	if len(arrayParams) >= 5 {
		maxDefragmentQuantityTemp, ok := arrayParams[4].(int64)
		if !ok {
			maxDefragmentQuantityTemp = 32
		}
		if maxDefragmentQuantityTemp > 32 || maxDefragmentQuantityTemp <= 0 {
			maxDefragmentQuantityTemp = 32
		}
		maxDefragmentQuantity = int(maxDefragmentQuantityTemp)
	}
	//end

	// get list outputcoins tx
	prvCoinID := &common.Hash{}
	err := prvCoinID.SetBytes(common.PRVCoinID[:])
	if err != nil {
		return nil, err
	}

	outCoins, err := rpcclient.GetUnspentOutputCoins(txService.RpcClient, txService.KeyWallet, prvCoinID)
	if err != nil {
		return nil, err
	}

	outCoins, amount := txService.calculateOutputCoinsByMinValue(outCoins, maxVal, maxDefragmentQuantity)

	if len(outCoins) == 0 {
		return nil, errors.New("outCoins is empty")
	}

	senderKeySet, shardIDSender, err := bean.GetKeySetFromPrivateKeyParams(senderKeyParam)
	if err != nil {
		return nil, err
	}

	// add more into output for estimate fee
	paymentInfo := &privacy.PaymentInfo{
		Amount:         amount,
		PaymentAddress: senderKeySet.PaymentAddress,
		Message:        []byte{},
	}

	paymentInfos := []*privacy.PaymentInfo{paymentInfo}

	realFee, _, _, err := txService.estimateFee(
		estimateFeeCoinPerKb,
		false,
		outCoins,
		paymentInfos,
		shardIDSender,
		0,
		hasPrivacyCoin,
		metadataParam,
		nil,
		0,
	)

	if err != nil {
		return nil, err
	}

	if len(outCoins) == 0 {
		realFee = 0
	}

	if amount < realFee {
		return nil, errors.New("amount must large ethan fee")
	}
	paymentInfo.Amount = amount - realFee

	// convert to inputcoins
	inputCoins := transaction.ConvertOutputCoinToInputCoin(outCoins)

	// init tx
	tx := transaction.Tx{}

	err = tx.Init(
		transaction.NewTxPrivacyInitParams(
			&senderKeySet.PrivateKey,
			paymentInfos,
			inputCoins,
			outCoins,
			realFee,
			hasPrivacyCoin,
			nil, // use for prv coin -> nil is valid
			metadataParam,
			nil,
		),
		txService.RpcClient,
		txService.KeyWallet,
	)

	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (txService TxService) calculateOutputCoinsByMinValue(outCoins []*privacy.OutputCoin, maxVal uint64, maxDefragmentQuantityTemp int) ([]*privacy.OutputCoin, uint64) {
	outCoinsTmp := make([]*privacy.OutputCoin, 0)
	amount := uint64(0)
	for _, outCoin := range outCoins {
		if outCoin.CoinDetails.GetValue() <= maxVal {
			outCoinsTmp = append(outCoinsTmp, outCoin)
			amount += outCoin.CoinDetails.GetValue()
			if len(outCoinsTmp) >= maxDefragmentQuantityTemp {
				break
			}
		}
	}
	return outCoinsTmp, amount
}

func (txService TxService) buildDefragmentPrivacyCustomTokenParam(
	tokenParamsRaw map[string]interface{},
	senderKeySet *incognitokey.KeySet,
	shardIDSender byte,
) (
	*transaction.CustomTokenPrivacyParamTx,
	map[common.Hash]transaction.TxCustomTokenPrivacy,
	error,
) {
	property, ok := tokenParamsRaw["TokenID"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid Token ID, Params %+v ", tokenParamsRaw)
	}

	tokenName, ok := tokenParamsRaw["TokenName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid Token Name, Params %+v ", tokenParamsRaw)
	}
	tokenSymbol, ok := tokenParamsRaw["TokenSymbol"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid Token Symbol, Params %+v ", tokenParamsRaw)
	}
	tokenTxType, ok := tokenParamsRaw["TokenTxType"].(int)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid Token Tx Type, Params %+v ", tokenParamsRaw)
	}
	tokenAmount, ok := tokenParamsRaw["TokenAmount"].(uint64)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid Token Amount, Params %+v ", tokenParamsRaw)
	}

	tokenFee, ok := tokenParamsRaw["TokenFee"].(uint64)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid Token Fee, Params %+v ", tokenParamsRaw)
	}

	tokenParams := &transaction.CustomTokenPrivacyParamTx{
		PropertyID:     property,
		PropertyName:   tokenName,
		PropertySymbol: tokenSymbol,
		TokenTxType:    tokenTxType,
		Amount:         tokenAmount,
		TokenInput:     nil,
		Fee:            tokenFee,
	}

	maxDefragmentQuantity := 32

	// get list custom token
	switch tokenParams.TokenTxType {
	case transaction.CustomTokenTransfer:
		tokenID, err := common.Hash{}.NewHashFromStr(tokenParams.PropertyID)
		if err != nil {
			return nil, nil, err
		}

		outputTokens, err := rpcclient.GetUnspentOutputCoins(txService.RpcClient, txService.KeyWallet, tokenID)
		if err != nil {
			return nil, nil, err
		}

		candidateOutputTokens, amount := txService.calculateOutputCoinsByMinValue(outputTokens, 10000*1e9, maxDefragmentQuantity)

		if len(candidateOutputTokens) == 0 {
			return nil, nil, errors.New("lis output coin is empty")
		}

		intputToken := transaction.ConvertOutputCoinToInputCoin(candidateOutputTokens)
		tokenParams.TokenInput = intputToken
		tokenParams.TokenOutput = candidateOutputTokens
		tokenParams.Receiver = []*privacy.PaymentInfo{{
			PaymentAddress: senderKeySet.PaymentAddress,
			Amount:         amount,
		}}
	}

	return tokenParams, nil, nil
}

func (txService TxService) buildDefragmentTokenParam(tokenParamsRaw map[string]interface{}, senderKeySet *incognitokey.KeySet, shardIDSender byte) (*transaction.CustomTokenPrivacyParamTx, error) {
	var privacyTokenParam *transaction.CustomTokenPrivacyParamTx
	var err error

	isPrivacy, ok := tokenParamsRaw["Privacy"].(bool)
	if !ok {
		return nil, errors.New("Invalid params")
	}
	if !isPrivacy {
		// Check normal custom token param
	} else {
		// Check privacy custom token param
		privacyTokenParam, _, err = txService.buildDefragmentPrivacyCustomTokenParam(tokenParamsRaw, senderKeySet, shardIDSender)
		if err != nil {
			return nil, err
		}
	}

	return privacyTokenParam, nil
}

func (txService TxService) BuildDeFragmentPTokenRawTransaction(
	params interface{},
	metaData metadata.Metadata,
) (*transaction.TxCustomTokenPrivacy, error) {
	txParam, errParam := bean.NewCreateRawPrivacyTokenTxParam(params)
	if errParam != nil {
		return nil, errParam
	}
	tokenParamsRaw := txParam.TokenParamsRaw
	var err error
	tokenParams, err := txService.buildDefragmentTokenParam(tokenParamsRaw, txParam.SenderKeySet, txParam.ShardIDSender)

	if err != nil {
		return nil, err
	}

	if tokenParams == nil {
		return nil, errors.New("can not build token params for request")
	}

	/******* START choose output native coins(PRV), which is used to create tx *****/
	var inputCoins []*privacy.InputCoin
	var outputPrvCoins []*privacy.OutputCoin
	realFeePRV := uint64(0)

	inputCoins, outputPrvCoins, realFeePRV, err = txService.chooseOutsCoinByKeyset(
		txParam.PaymentInfos,
		txParam.EstimateFeeCoinPerKb,
		0,
		txParam.SenderKeySet,
		txParam.ShardIDSender,
		txParam.HasPrivacyCoin,
		nil,
		tokenParams,
		txParam.IsGetPTokenFee,
		txParam.UnitPTokenFee,
	)

	if err != nil {
		return nil, err
	}

	if len(txParam.PaymentInfos) == 0 && realFeePRV == 0 {
		txParam.HasPrivacyCoin = false
	}

	/******* END GET output coins native coins(PRV), which is used to create tx *****/
	tx := &transaction.TxCustomTokenPrivacy{}
	err = tx.Init(
		transaction.NewTxPrivacyTokenInitParams(
			&txParam.SenderKeySet.PrivateKey,
			txParam.PaymentInfos,
			inputCoins,
			outputPrvCoins,
			realFeePRV,
			tokenParams,
			metaData,
			txParam.HasPrivacyCoin,
			txParam.HasPrivacyToken,
			txParam.ShardIDSender,
			txParam.Info,
		),
		txService.RpcClient,
		txService.KeyWallet,
	)

	if err != nil {
		return nil, err
	}

	return tx, nil
}
