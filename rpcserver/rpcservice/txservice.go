package rpcservice

import (
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/incognitokey"
	"github.com/incognitochain/go-incognito-sdk/mempool"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/transaction"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"sort"
)

type TxService struct {
	RpcClient    *rpcclient.HttpClient
	Wallet       *wallet.Wallet
	KeyWallet    *wallet.KeyWallet
	FeeEstimator map[byte]*mempool.FeeEstimator
}

func (txService TxService) BuildRawTransaction(params *bean.CreateRawTxParam, meta metadata.Metadata) (*transaction.Tx, error) {
	// get output coins to spend and real fee
	inputCoins, outputCoin, realFee, err := txService.chooseOutsCoinByKeyset(
		params.PaymentInfos,
		params.EstimateFeeCoinPerKb,
		0,
		params.SenderKeySet,
		params.ShardIDSender,
		params.HasPrivacyCoin,
		meta,
		nil,
		false,
		int64(0),
	)

	if err != nil {
		return nil, err
	}

	// init tx
	tx := transaction.Tx{}

	err = tx.Init(
		transaction.NewTxPrivacyInitParams(
			&params.SenderKeySet.PrivateKey,
			params.PaymentInfos,
			inputCoins,
			outputCoin,
			realFee,
			params.HasPrivacyCoin,
			nil, // use for prv coin -> nil is valid
			meta,
			params.Info,
		),
		txService.RpcClient,
		txService.KeyWallet,
	)

	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (txService TxService) chooseOutsCoinByKeyset(
	paymentInfos []*privacy.PaymentInfo,
	unitFeeNativeToken int64,
	numBlock uint64,
	keySet *incognitokey.KeySet,
	shardIDSender byte,
	hasPrivacy bool,
	metadataParam metadata.Metadata,
	privacyCustomTokenParams *transaction.CustomTokenPrivacyParamTx,
	isGetFeePToken bool,
	unitFeePToken int64,
) ([]*privacy.InputCoin, []*privacy.OutputCoin, uint64, error) {
	// calculate total amount to send
	totalAmmount := uint64(0)
	for _, receiver := range paymentInfos {
		totalAmmount += receiver.Amount
	}

	// get list outputcoins tx
	prvCoinID := &common.Hash{}
	err := prvCoinID.SetBytes(common.PRVCoinID[:])
	if err != nil {
		return nil, nil, 0, err
	}

	outCoins, err := rpcclient.GetUnspentOutputCoins(txService.RpcClient, txService.KeyWallet, prvCoinID)
	if err != nil {
		return nil, nil, 0, err
	}

	if len(outCoins) == 0 && totalAmmount > 0 {
		return nil, nil, 0, errors.New("not enough output coin")
	}

	// Use Knapsack to get candiate output coin
	candidateOutputCoins, outCoins, candidateOutputCoinAmount, err := txService.chooseBestOutCoinsToSpent(outCoins, totalAmmount)
	if err != nil {
		return nil, nil, 0, err
	}

	//todo
	// refund out put for sender
	overBalanceAmount := candidateOutputCoinAmount - totalAmmount
	if overBalanceAmount > 0 {
		// add more into output for estimate fee
		paymentInfos = append(paymentInfos, &privacy.PaymentInfo{
			PaymentAddress: keySet.PaymentAddress,
			Amount:         overBalanceAmount,
		})
	}

	realFee, _, _, err := txService.estimateFee(
		unitFeeNativeToken,
		false,
		candidateOutputCoins,
		paymentInfos,
		shardIDSender,
		numBlock,
		hasPrivacy,
		metadataParam,
		privacyCustomTokenParams,
		0,
	)

	if err != nil {
		return nil, nil, 0, err
	}

	if totalAmmount == 0 && realFee == 0 {
		if metadataParam != nil {
			metadataType := metadataParam.GetType()
			switch metadataType {
			case metadata.WithDrawRewardRequestMeta:
				{
					return nil, nil, realFee, nil
				}
			}
			return nil, nil, realFee, fmt.Errorf("totalAmmount: %+v, realFee: %+v", totalAmmount, realFee)
		}

		if privacyCustomTokenParams != nil {
			// for privacy token
			return nil, nil, 0, nil
		}
	}

	needToPayFee := int64((totalAmmount + realFee) - candidateOutputCoinAmount)
	// if not enough to pay fee
	if needToPayFee > 0 {
		if len(outCoins) > 0 {
			candidateOutputCoinsForFee, _, _, err1 := txService.chooseBestOutCoinsToSpent(outCoins, uint64(needToPayFee))
			if err1 != nil {
				return nil, nil, 0, err1
			}
			candidateOutputCoins = append(candidateOutputCoins, candidateOutputCoinsForFee...)
		}
	}

	// convert to inputcoins
	inputCoins := transaction.ConvertOutputCoinToInputCoin(candidateOutputCoins)
	return inputCoins, candidateOutputCoins, realFee, nil
}

func (txService TxService) chooseBestOutCoinsToSpent(outCoins []*privacy.OutputCoin, amount uint64) (resultOutputCoins []*privacy.OutputCoin, remainOutputCoins []*privacy.OutputCoin, totalResultOutputCoinAmount uint64, err error) {
	resultOutputCoins = make([]*privacy.OutputCoin, 0)
	remainOutputCoins = make([]*privacy.OutputCoin, 0)
	totalResultOutputCoinAmount = uint64(0)

	// either take the smallest coins, or a single largest one
	var outCoinOverLimit *privacy.OutputCoin
	outCoinsUnderLimit := make([]*privacy.OutputCoin, 0)
	for _, outCoin := range outCoins {
		if outCoin.CoinDetails.GetValue() < amount {
			outCoinsUnderLimit = append(outCoinsUnderLimit, outCoin)
		} else if outCoinOverLimit == nil {
			outCoinOverLimit = outCoin
		} else if outCoinOverLimit.CoinDetails.GetValue() > outCoin.CoinDetails.GetValue() {
			remainOutputCoins = append(remainOutputCoins, outCoin)
		} else {
			remainOutputCoins = append(remainOutputCoins, outCoinOverLimit)
			outCoinOverLimit = outCoin
		}
	}
	sort.Slice(outCoinsUnderLimit, func(i, j int) bool {
		return outCoinsUnderLimit[i].CoinDetails.GetValue() < outCoinsUnderLimit[j].CoinDetails.GetValue()
	})
	for _, outCoin := range outCoinsUnderLimit {
		if totalResultOutputCoinAmount < amount {
			totalResultOutputCoinAmount += outCoin.CoinDetails.GetValue()
			resultOutputCoins = append(resultOutputCoins, outCoin)
		} else {
			remainOutputCoins = append(remainOutputCoins, outCoin)
		}
	}
	if outCoinOverLimit != nil && (outCoinOverLimit.CoinDetails.GetValue() > 2*amount || totalResultOutputCoinAmount < amount) {
		remainOutputCoins = append(remainOutputCoins, resultOutputCoins...)
		resultOutputCoins = []*privacy.OutputCoin{outCoinOverLimit}
		totalResultOutputCoinAmount = outCoinOverLimit.CoinDetails.GetValue()
	} else if outCoinOverLimit != nil {
		remainOutputCoins = append(remainOutputCoins, outCoinOverLimit)
	}

	if totalResultOutputCoinAmount < amount {
		return resultOutputCoins, remainOutputCoins, totalResultOutputCoinAmount, errors.New("Not enough coin")
	} else {
		return resultOutputCoins, remainOutputCoins, totalResultOutputCoinAmount, nil
	}
}

func (txService TxService) estimateFee(
	defaultFee int64,
	isGetPTokenFee bool,
	candidateOutputCoins []*privacy.OutputCoin,
	paymentInfos []*privacy.PaymentInfo,
	shardID byte,
	numBlock uint64,
	hasPrivacy bool,
	metadata metadata.Metadata,
	privacyCustomTokenParams *transaction.CustomTokenPrivacyParamTx,
	beaconHeight int64,
) (uint64, uint64, uint64, error) {
	// check real fee(nano PRV) per tx
	var realFee uint64
	estimateFeeCoinPerKb := uint64(0)
	estimateTxSizeInKb := uint64(0)

	tokenId := &common.Hash{}
	if isGetPTokenFee {
		if privacyCustomTokenParams != nil {
			tokenId, _ = common.Hash{}.NewHashFromStr(privacyCustomTokenParams.PropertyID)
		}
	} else {
		tokenId = nil
	}

	paymentAddrStr := txService.KeyWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	//payment address from private key
	estimateFeeCoinPerKb, err := rpcclient.GetEstimateFeeWithEstimator(txService.RpcClient, defaultFee, paymentAddrStr, tokenId)
	if err != nil {
		return 0, 0, 0, err
	}

	if txService.Wallet != nil {
		estimateFeeCoinPerKb += uint64(txService.Wallet.GetConfig().IncrementalFee)
	}

	limitFee := uint64(0)
	if feeEstimator, ok := txService.FeeEstimator[shardID]; ok {
		limitFee = feeEstimator.GetLimitFeeForNativeToken()
	}

	estimateTxSizeInKb = transaction.EstimateTxSize(transaction.NewEstimateTxSizeParam(len(candidateOutputCoins), len(paymentInfos), hasPrivacy, metadata, privacyCustomTokenParams, limitFee))
	realFee = uint64(estimateFeeCoinPerKb) * uint64(estimateTxSizeInKb)

	fmt.Println(fmt.Sprintf("default fee: %v, estimate fee: %v, estimate tx size (kb) %v, real fee %v", defaultFee, estimateFeeCoinPerKb, estimateTxSizeInKb, realFee))

	return realFee, estimateFeeCoinPerKb, estimateTxSizeInKb, nil
}
