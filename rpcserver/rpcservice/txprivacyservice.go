package rpcservice

import (
	"errors"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/incognitokey"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/transaction"
)

func (txService TxService) BuildRawPrivacyCustomTokenTransaction(params interface{}, metaData metadata.Metadata) (*transaction.TxCustomTokenPrivacy, error) {
	txParam, errParam := bean.NewCreateRawPrivacyTokenTxParam(params)
	if errParam != nil {
		return nil, errParam
	}
	tokenParamsRaw := txParam.TokenParamsRaw
	var err error
	tokenParams, err := txService.buildTokenParam(tokenParamsRaw, txParam.SenderKeySet, txParam.ShardIDSender)

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

func (txService TxService) buildTokenParam(tokenParamsRaw map[string]interface{}, senderKeySet *incognitokey.KeySet, shardIDSender byte) (*transaction.CustomTokenPrivacyParamTx, error) {
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
		privacyTokenParam, _, err = txService.buildPrivacyCustomTokenParam(tokenParamsRaw, senderKeySet, shardIDSender)
		if err != nil {
			return nil, err
		}
	}

	return privacyTokenParam, nil
}

func (txService TxService) buildPrivacyCustomTokenParam(
	tokenParamsRaw map[string]interface{},
	senderKeySet *incognitokey.KeySet,
	shardIDSender byte,
) (*transaction.CustomTokenPrivacyParamTx, map[common.Hash]transaction.TxCustomTokenPrivacy, error) {
	property, ok := tokenParamsRaw["TokenID"].(string)
	if !ok {
		return nil, nil, errors.New("Invalid Token ID")
	}
	_, ok = tokenParamsRaw["TokenReceivers"]
	if !ok {
		return nil, nil, errors.New("Token Receiver is invalid")
	}
	tokenName, ok := tokenParamsRaw["TokenName"].(string)
	if !ok {
		return nil, nil, errors.New("Invalid Token Name")
	}
	tokenSymbol, ok := tokenParamsRaw["TokenSymbol"].(string)
	if !ok {
		return nil, nil, errors.New("Invalid Token Symbol")
	}
	tokenTxType, ok := tokenParamsRaw["TokenTxType"].(int)
	if !ok {
		return nil, nil, errors.New("Invalid Token Tx Type")
	}

	tokenAmount, ok := tokenParamsRaw["TokenAmount"].(uint64)
	if !ok {
		return nil, nil, errors.New("Invalid Token Amount")
	}

	tokenFee, ok := tokenParamsRaw["TokenFee"].(uint64)
	if !ok {
		return nil, nil, errors.New("Invalid Token Fee")
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
	voutsAmount := int64(0)
	var err1 error

	tokenParams.Receiver, voutsAmount, err1 = transaction.CreateCustomTokenPrivacyReceiverArray(tokenParamsRaw["TokenReceivers"])
	if err1 != nil {
		return nil, nil, err1
	}
	voutsAmount += int64(tokenFee)

	// get list custom token
	switch tokenParams.TokenTxType {
	case transaction.CustomTokenTransfer:
		{
			tokenID, err := common.Hash{}.NewHashFromStr(tokenParams.PropertyID)
			if err != nil {
				return nil, nil, err
			}

			outputTokens, err := rpcclient.GetUnspentOutputCoins(txService.RpcClient, txService.KeyWallet, tokenID)
			if err != nil {
				return nil, nil, err
			}

			if len(outputTokens) == 0 && voutsAmount > 0 {
				return nil, nil, errors.New("not enough output coin")
			}

			candidateOutputTokens, _, _, err := txService.chooseBestOutCoinsToSpent(outputTokens, uint64(voutsAmount))
			if err != nil {
				return nil, nil, err
			}

			intputToken := transaction.ConvertOutputCoinToInputCoin(candidateOutputTokens)
			tokenParams.TokenInput = intputToken
			tokenParams.TokenOutput = candidateOutputTokens
		}
	}

	return tokenParams, nil, nil
}