package repository

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/constant"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/entity"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/service"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"github.com/pkg/errors"
)

type Wallet struct {
	Inc                 *service.IncogClient
	ConstantID          string
	Block               *Block
	IncChainIntegration *IncChainIntegration
}

func NewWallet(inc *service.IncogClient, constantID string, block *Block, incChainIntegration *IncChainIntegration) *Wallet {
	return &Wallet{Inc: inc, ConstantID: constantID, Block: block, IncChainIntegration: incChainIntegration}
}

func (w *Wallet) CreateWalletAddress() (paymentAddress, pubkey, readonlyKey, privateKey string, err error) {
	wallet, err := w.IncChainIntegration.CreateWalletAddress()

	if err != nil {
		return "", "", "", "", errors.Wrap(err, "w.IncChainIntegration.CreateWalletAddress")
	}

	paymentAddress = wallet.PaymentAddress
	readonlyKey = wallet.ReadonlyKey
	pubkey = wallet.Pubkey
	privateKey = wallet.PrivateKey
	return
}

func (w *Wallet) CreateNodeWalletAddress(byShardId int) (paymentAddress, pubkey, readonlyKey, privateKey, validatorKey string, shardId int, err error) {
	var newWallet *wallet.KeySerializedData
	var newErrors error

	if byShardId > -1 {
		newWallet, newErrors = w.IncChainIntegration.CreateNewWalletByShardId(byShardId)
	} else {
		newWallet, newErrors = w.IncChainIntegration.CreateWalletAddress()
	}

	if newErrors != nil {
		return "", "", "", "", "", -1, errors.Wrap(newErrors, "w.IncChainIntegration.CreateWalletAddress")
	}

	paymentAddress = newWallet.PaymentAddress
	readonlyKey = newWallet.ReadonlyKey
	pubkey = newWallet.Pubkey
	privateKey = newWallet.PrivateKey
	validatorKey = newWallet.ValidatorKey
	shardId = newWallet.ShardId

	return
}

func (w *Wallet) CreateWalletAddressByShardId(byShardId int) (paymentAddress, pubkey, readonlyKey, privateKey string, shardId int, err error) {
	wallet, err := w.IncChainIntegration.CreateNewWalletByShardId(byShardId)

	if err != nil {
		return "", "", "", "", 0, errors.Wrap(err, "w.IncChainIntegration.CreateWalletAddress")
	}

	paymentAddress = wallet.PaymentAddress
	readonlyKey = wallet.ReadonlyKey
	pubkey = wallet.Pubkey
	privateKey = wallet.PrivateKey
	shardId = wallet.ShardId
	return
}

func (w *Wallet) ListRewardAmountAll() ([]entity.RewardData, error) {
	param := []interface{}{}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.ListRewardAmount, param)
	if err != nil {
		return nil, errors.Wrap(err, "w.ListRewardAmounts")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}

	result, ok := data["Result"].(map[string]interface{})

	if !ok {
		return nil, errors.Errorf("couldn't get Result: resp: %+v", data)
	}

	// var rewards []entity.RewardItems
	var rewardDatas []entity.RewardData

	for k, v := range result {

		rewardData := entity.RewardData{
			PublicKey: k,
		}

		rewardItemMap := v.(map[string]interface{})

		var rewardDataItems []entity.RewardDataItem
		for k2, v2 := range rewardItemMap {

			rewardDataItems = append(rewardDataItems, entity.RewardDataItem{
				TokenId: k2,
				Reward:  uint64(v2.(float64)),
			})
			rewardData.RewardItems = rewardDataItems

		}
		rewardDatas = append(rewardDatas, rewardData)
	}

	return rewardDatas, nil
}

func (w *Wallet) GetBalanceByPrivateKey(privateKey string) (uint64, error) {
	//rpc: GetBalanceByPrivateKeyMethod
	amount, err := w.IncChainIntegration.GetBalance(privateKey, w.ConstantID)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func (w *Wallet) GetBalanceByPaymentAddress(paymentAddress string) (uint64, error) {
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetBalanceByPaymentAddress, []interface{}{paymentAddress})
	if err != nil {
		return 0, err
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return 0, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return 0, nil
	}
	if v, ok := data["Result"].(float64); ok {
		return uint64(v), nil
	}
	return 0, nil
}

func (w *Wallet) GetListCustomTokenBalance(paymentAddress string) (*entity.ListCustomTokenBalance, error) {
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetListCustomTokenBalance, []interface{}{paymentAddress})
	if err != nil {
		return nil, err
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response: req: %+v, resp: %+v", paymentAddress, data)
	}
	resultResp := data["Result"].(map[string]interface{})
	resultRespStr, err := json.Marshal(resultResp)
	if err != nil {
		return nil, err
	}
	var result entity.ListCustomTokenBalance
	err = json.Unmarshal(resultRespStr, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (w *Wallet) GetListPrivacyCustomTokenBalanceByID(privateKey, tokenID string) (*big.Int, error) {
	//rpc: GetBalanceByPrivateKeyMethod
	amount, err := w.IncChainIntegration.GetBalance(privateKey, tokenID)
	if err != nil {
		return nil, err
	}

	v := new(big.Int)
	v.SetUint64(amount)
	return v, nil
}

// dont use
func (w *Wallet) GetAmountVoteToken(paymentAddress string) (*entity.ListCustomTokenBalance, error) {
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetAmountVoteToken, []interface{}{paymentAddress, 0})
	if err != nil {
		return nil, errors.Wrap(err, "w.blockchainAPI")
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response: req: %+v, resp: %+v", paymentAddress, data)
	}
	var result entity.ListCustomTokenBalance
	r, err := json.Marshal(data["Result"])
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	if err := json.Unmarshal(r, &result); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return &result, nil
}

func (w *Wallet) createAndSendConstantPrivacyTransaction(privateKey string, req entity.WalletSend) (string, error) {
	param := []interface{}{privateKey, req.PaymentAddresses, constant.EstimateFee, 1}

	//rpc: CreateAndSendTransaction
	rawData, err := w.IncChainIntegration.CreateAndSendConstantTransaction(param)
	if err != nil {
		return "", errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method CreateAndSendConstantPrivacyTransaction: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)

	if err != nil {
		return "", errors.Wrap(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (w *Wallet) EstimatePRVFee(privateKey, toAddress string, amountToSend uint64) (int, int, error) {
	param := []interface{}{privateKey, map[string]uint64{toAddress: amountToSend}, -1, 0}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetEstimateFee, param)
	if err != nil {
		return 0, 0, errors.Wrap(err, "w.blockchainAPI")
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return 0, 0, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return 0, 0, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return 0, 0, errors.Errorf("couldn't get result:  resp: %+v", data)
	}

	estimateFeeCoinPerKb, ok := result["EstimateFeeCoinPerKb"].(float64)
	if !ok {
		return 0, 0, errors.Errorf("couldn't get estimateFeeCoinPerKb: result: %+v", result)
	}
	estimateTxSizeInKb, ok := result["EstimateTxSizeInKb"].(float64)
	if !ok {
		return 0, 0, errors.Errorf("couldn't get estimateTxSizeInKb: result: %+v", result)
	}
	return int(estimateFeeCoinPerKb), int(estimateTxSizeInKb), nil
}

// send max prv:
func (w *Wallet) CreateAndSendMaxPRVTransaction(privateKey, toAddress string) (string, error) {

	prvBalance, err := w.GetBalanceByPrivateKey(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "w.GetBalanceByPrivateKey")
	}
	fmt.Println("max amount:", prvBalance)

	// est fee:
	estimateFeeCoinPerKb, estimateTxSizeInKb, err := w.EstimatePRVFee(privateKey, toAddress, prvBalance)

	if err != nil {
		return "", errors.Wrap(err, "w.EstimateFee")
	}

	estimateFee := estimateFeeCoinPerKb * estimateTxSizeInKb

	maxAmount := prvBalance - uint64(estimateFee)

	param := []interface{}{privateKey, map[string]uint64{toAddress: maxAmount}, estimateFee, 1}

	//rpc: CreateAndSendTransaction
	rawData, err := w.IncChainIntegration.CreateAndSendConstantTransaction(param)
	if err != nil {
		return "", errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method CreateAndSendMaxPRVTransaction: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)

	if err != nil {
		return "", errors.Wrap(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (w *Wallet) sendPrivacyCustomTokenTransaction(privateKey string, req entity.WalletSend) (map[string]interface{}, error) {
	tokenData := map[string]interface{}{}
	tokenData["Privacy"] = true
	tokenData["TokenID"] = req.TokenID
	tokenData["TokenTxType"] = req.Type
	tokenData["TokenName"] = req.TokenName
	tokenData["TokenSymbol"] = req.TokenSymbol
	tokenData["TokenReceivers"] = req.PaymentAddresses
	tokenData["TokenAmount"] = req.TokenAmount
	tokenData["TokenFee"] = req.TokenFee
	object := map[string]uint64{}

	nativeFee := -1
	if req.TokenFee > 0 {
		nativeFee = 0
	}

	param := []interface{}{privateKey, object, nativeFee, 1, tokenData, "", 1}

	//rpc: CreateAndSendPrivacyCustomTokenTransaction
	rawData, err := w.IncChainIntegration.SendPrivacyCustomTokenTransaction(param)
	if err != nil {
		return nil, errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method SendPrivacyCustomTokenTransaction: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawPrivacyCustomTokenTransaction, rawData)

	if err != nil {
		return nil, errors.Wrap(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	return result, nil
}

func (w *Wallet) ListPrivacyCustomToken() ([]entity.PCustomToken, error) {
	param := []interface{}{}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.ListPrivacyCustomToken, param)
	if err != nil {
		return nil, errors.Wrap(err, "w.ListPrivacyCustomToken")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	// jsonData, err := json.Marshal(result["ListCustomToken"])

	jsonData, err := json.MarshalIndent(result["ListCustomToken"], "", "\t")

	fmt.Println("err", err)

	// fmt.Println("jsonData", jsonData)

	var finalData *[]entity.PCustomToken

	if err := json.Unmarshal(jsonData, &finalData); err != nil {
		fmt.Println("err Unmarshal", err)
		return nil, errors.Errorf("invalid data")
	}
	return *finalData, nil
}

func (w *Wallet) GetTxByHash(txHash string) (*entity.TransactionDetail, error) {
	param := []string{txHash}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetTransactionByHash, param)
	if err != nil {
		return nil, errors.Wrapf(err, "w.blockchainAPI: param: %+v", param)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	if result["Hash"] == nil {
		return nil, constant.ErrTxHashNotExists
	}

	sb, err := json.Marshal(result)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	var tx entity.TransactionDetail
	if err := json.Unmarshal(sb, &tx); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal: sb: %+v", string(sb))
	}
	return &tx, nil
}

func (w *Wallet) GetDecryptOutputCoinByKeyOfTransaction(txHash, paymentAddress, readonlyKey string) (*entity.DecrypTransactionPRV, error) {

	param := []interface{}{txHash, map[string]string{"PaymentAddress": paymentAddress, "ReadonlyKey": readonlyKey}}

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.DecryptOutputCoinByKeyOfTransaction, param)
	if err != nil {
		return nil, errors.Wrapf(err, "w.blockchainAPI: param: %+v", param)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}

	sb, err := json.Marshal(result)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	var decrypTransactionPRV entity.DecrypTransactionPRV
	if err := json.Unmarshal(sb, &decrypTransactionPRV); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal: sb: %+v", string(sb))
	}
	fmt.Println(decrypTransactionPRV)
	return &decrypTransactionPRV, nil
}

func (w *Wallet) GetDecryptOutputCoinByKeyOfTrans(txHash, paymentAddress, readonlyKey string) (map[string]interface{}, error) {
	var results map[string]interface{}
	param := []interface{}{txHash, map[string]string{"PaymentAddress": paymentAddress, "ReadonlyKey": readonlyKey}}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.DecryptOutputCoinByKeyOfTransaction, param)
	if err != nil {
		return results, errors.Wrapf(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return results, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response: resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return results, errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	sb, err := json.Marshal(result)
	if err != nil {
		return results, errors.Wrap(err, "json.Marshal")
	}
	if err := json.Unmarshal(sb, &results); err != nil {
		return results, errors.Wrapf(err, "json.Unmarshal: sb: %+v", string(sb))
	}
	return results, nil
}

//ProofDetail
// get amount by hash public:
func (w *Wallet) GetAmountByHashFromReceiveAddressAndToAddress(txHash, fromAddress, toAddress string) (*big.Int, error) {

	// convert address to PublicKey
	fromPublicKey, err := w.GetPublickeyFromPaymentAddress(fromAddress)

	if err != nil {
		return nil, errors.Wrapf(err, "w.blockchainAPI: param: %+v", fromAddress)
	}
	toPublicKey, err := w.GetPublickeyFromPaymentAddress(toAddress)
	fmt.Println("toPublicKey, err", toPublicKey, err)
	if err != nil {
		return nil, errors.Wrapf(err, "w.blockchainAPI: param: %+v", fromAddress)
	}

	param := []string{txHash}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetTransactionByHash, param)
	if err != nil {
		return nil, errors.Wrapf(err, "w.blockchainAPI: param: %+v", param)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	if result["Hash"] == nil {
		return nil, constant.ErrTxHashNotExists
	}

	sb, err := json.Marshal(result)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	var tx entity.TransactionDetail
	if err := json.Unmarshal(sb, &tx); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal: sb: %+v", string(sb))
	}
	// get amount:
	var sendAmount uint64 = 0
	var receiveFromPublicKey = false
	var sentToPublicKey = false

	if tx.ProofDetail != nil {
		// get amount to send to toPublicKey
		if tx.ProofDetail.OutputCoins != nil {
			if len(tx.ProofDetail.OutputCoins) > 0 {
				for _, coinDetail := range tx.ProofDetail.OutputCoins {
					if coinDetail.CoinDetails != nil {
						// check public key mathed to toPublicKey
						if toPublicKey == coinDetail.CoinDetails.PublicKey {
							sentToPublicKey = true
							sendAmount = coinDetail.CoinDetails.Value
							break
						}
					}
				}
			}
		}
		// check from address mathed to fromPublicKey??
		if tx.ProofDetail.InputCoins != nil {
			if len(tx.ProofDetail.InputCoins) > 0 {
				for _, coinDetail := range tx.ProofDetail.InputCoins {
					if coinDetail.CoinDetails != nil {
						// check public key mathed to romPublicKey
						if fromPublicKey == coinDetail.CoinDetails.PublicKey {
							receiveFromPublicKey = true
							break
						}
					}
				}
			}
		}
	}

	if receiveFromPublicKey == false {
		return nil, constant.ErrTxHashInvalidFromAddress
	}
	if sentToPublicKey == false {
		return nil, constant.ErrTxHashInvalidToAddress
	}

	fmt.Println("sendAmount:", sendAmount)
	v := new(big.Int)
	v.SetUint64(sendAmount)
	return v, nil
}

func (w *Wallet) CreateAndSendIssuingRequest(privateKey, receiveAddress string, depositedAmount *big.Int, tokenId string, tokenName string) (string, error) {
	if depositedAmount == nil {
		return "", errors.New("depositedAmount is nil")
	}

	depositedReq := map[string]interface{}{}
	depositedReq["ReceiveAddress"] = receiveAddress
	depositedReq["DepositedAmount"] = depositedAmount.Uint64()
	depositedReq["TokenID"] = tokenId
	depositedReq["TokenName"] = tokenName

	param := []interface{}{privateKey, nil, constant.EstimateFee, -1, depositedReq}

	//rpc: CreateAndSendIssuingRequest
	rawData, err := w.IncChainIntegration.CreateAndSendIssuingRequest(param)
	if err != nil {
		return "", errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method CreateAndSendIssuingRequest: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)
	//resp, _, err := w.Inc.PostAndReceiveInterface(CreateAndSendIssuingRequest, param)
	if err != nil {
		return "", errors.Wrap(err, "w.blockchainAPI")
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (w *Wallet) GetIssuingStatus(txHash string) (string, uint64, error) {
	param := []string{txHash}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetIssuingStatus, param)
	if err != nil {
		return "", 0, errors.Wrapf(err, "w.blockchainAPI: param: %+v", param)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", 0, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", 0, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", 0, errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	if result["Status"] == nil {
		return "", 0, errors.Errorf("bad result: data: %+v", data["Result"])
	}
	if result["Amount"] == nil {
		return "", 0, errors.Errorf("bad result: data: %+v", data["Result"])
	}
	return result["Status"].(string), uint64(result["Amount"].(float64)), nil
}

func (w *Wallet) GetContractingStatus(txHash string) (string, *big.Int, error) {
	param := []string{txHash}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetContractingStatus, param)
	if err != nil {
		return "", nil, errors.Wrapf(err, "w.blockchainAPI: param: %+v", param)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", nil, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", nil, errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	if result["Status"] == nil {
		return "", nil, errors.Errorf("bad result: data: %+v", data["Result"])
	}
	if result["Redeem"] == nil {
		return "", nil, errors.Errorf("bad result: data: %+v", data["Result"])
	}
	redeem, ok := new(big.Int).SetString(result["Redeem"].(string), 10)
	return result["Status"].(string), redeem, nil
}

func (w *Wallet) CreateAndSendIssuingRequestForPrivacyToken(privateKey string, metadata map[string]interface{}) (string, error) {
	param := []interface{}{privateKey, nil, constant.EstimateFee, -1, metadata}

	//rpc: CreateAndSendIssuingRequest
	rawData, err := w.IncChainIntegration.CreateAndSendIssuingRequest(param)
	if err != nil {
		return "", errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method CreateAndSendIssuingRequestForPrivacyToken: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)
	//resp, _, err := w.Inc.PostAndReceiveInterface(CreateAndSendIssuingRequest, param)
	if err != nil {
		return "", errors.Wrap(err, "w.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (w *Wallet) CreateAndSendContractingRequestForPrivacyToken(privateKey string, autoChargePRVFee int, metadata map[string]interface{}) (string, error) {
	// autoChargePRVFee: -1: auto prv fee, 0: 0 prv fee -> get ptoken fee
	param := []interface{}{
		privateKey,
		nil,
		autoChargePRVFee,
		-1,
		metadata,
		"",
		0,
	}

	//rpc: CreateAndSendContractingRequest
	rawData, err := w.IncChainIntegration.CreateAndSendContractingRequest(param)
	if err != nil {
		return "", errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method CreateAndSendContractingRequest: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawPrivacyCustomTokenTransaction, rawData)

	if err != nil {
		return "", errors.Wrap(err, "w.blockchainAPI")
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (w *Wallet) CreateAndSendTxWithIssuingEth(privateKey, burnerAddress string, metadata map[string]interface{}) (string, []byte, error) {
	transParams := map[string]uint64{burnerAddress: 0}
	param := []interface{}{privateKey, transParams, constant.EstimateFee, -1, metadata}

	//rpc: CreateAndSendTxWithIssuingEthReq
	rawData, err := w.IncChainIntegration.CreateAndSendTxWithIssuingEth(param)
	if err != nil {
		return "", nil, errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method CreateAndSendTxWithIssuingEth: %v \n", rawData)

	resp, body, err := w.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)

	if err != nil {
		return "", body, errors.Wrap(err, "w.blockchainAPI")
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", body, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", body, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", body, errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", body, errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, body, nil
}

func (w *Wallet) GetBridgeReqWithStatus(TxReqID string) (int, error) {

	txParams := map[string]interface{}{"TxReqID": TxReqID}

	param := []interface{}{txParams}

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetBridgeReqWithStatus, param)
	if err != nil {
		return -1, errors.Wrapf(err, "w.blockchainAPI: param: %+v", param)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return -1, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return -1, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}

	return int(data["Result"].(float64)), nil
}

// WithDrawReward
func (w *Wallet) CreateWithDrawReward(privateKey, tokenID string) (string, error) {

	param := []interface{}{
		privateKey,
		0,
		0,
		0,
		map[string]interface{}{"TokenID": tokenID},
	}

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.WithDrawReward, param)
	if err != nil {
		return "", errors.Wrapf(err, "w.blockchainAPI: param: %+v", param)
	}
	// todo: update code get response ...
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

// gen tokenid
func (w *Wallet) GenerateTokenID(symbol, pSymbol string) (string, error) {
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GenerateTokenID, []interface{}{symbol, pSymbol})
	if err != nil {
		return "", err
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", nil
	}
	tokenID, ok := data["Result"].(string)
	if !ok {
		return "", nil
	}
	return tokenID, nil
}

func (w *Wallet) GetPublickeyFromPaymentAddress(paymentAddress string) (string, error) {
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetPublickeyFromPaymentAddress, []interface{}{paymentAddress})

	if err != nil {
		return "", err
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", nil
	}
	result, ok := data["Result"].(map[string]interface{})

	if !ok {
		return "", errors.Errorf("couldn't get public ID: resp: %+v", data)
	}
	publicKeyInBase58Check, ok := result["PublicKeyInBase58Check"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get publicKeyInBase58Check: result: %+v", result)
	}
	return publicKeyInBase58Check, nil
}

func (w *Wallet) GetShardFromPaymentAddress(paymentAddress string) (int, error) {

	body, err := w.Inc.Post(constant.GetPublickeyFromPaymentAddress, []interface{}{paymentAddress})

	if err != nil {
		return -1, errors.Wrapf(err, "w.post: param: %+v", []interface{}{paymentAddress})
	}

	Message := entity.GetShardFromPaymentAddressObject{}

	if err := json.Unmarshal(body, &Message); err != nil {
		return -1, errors.Wrapf(err, "w.respone: param: %+v", body)
	}
	fmt.Println(Message.Result.PublicKeyInBytes[31])
	return (Message.Result.PublicKeyInBytes[31] % 8), nil
}

func (w *Wallet) getBurningAddressFromChain() (string, error) {

	param := []interface{}{}

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetBurningAddress, param)
	if err != nil {
		return "", errors.Wrapf(err, "w.blockchainAPI: param: %+v", param)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: resp: %+v", data)
	}

	return data["Result"].(string), nil
}

func (w *Wallet) GetTransactionByReceivers(PaymentAddress, ReadonlyKey string) (res *entity.ReceivedTransactions, err error) {
	param := []interface{}{map[string]string{"PaymentAddress": PaymentAddress, "ReadonlyKey": ReadonlyKey}}

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.GetTransactionByReceiver, param)

	if err != nil {
		return nil, err
	}
	data := resp.(map[string]interface{})
	resultResp := data["Result"]

	if resultResp == nil {
		return nil, errors.New("Fail")
	}

	result := entity.ReceivedTransactions{}
	resultInBytes, err := json.MarshalIndent(resultResp, "", "\t")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resultInBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (w *Wallet) GetAmountByMemo(listTrans []entity.ReceivedTransaction, memo, tokenID string) (string, uint64) {

	for _, receivedTransaction := range listTrans {
		memoDecode, _, err := base58.Base58Check{}.Decode(receivedTransaction.Info)

		if err != nil {
			fmt.Println("err", err)
			return "", 0
		}

		memoDecodeString := string(memoDecode)
		if strings.Contains(strings.ToLower(memoDecodeString), strings.ToLower(memo)) {
			// get amount for this memo:
			for token, ReceivedAmount := range receivedTransaction.ReceivedAmounts {
				if strings.ToLower(token) == strings.ToLower(tokenID) {
					return receivedTransaction.Hash, ReceivedAmount.CoinDetails.Value
				}
			}
		}
	}

	return "", 0
}

func (w *Wallet) CreateAndSendBurningForDepositToSCRequest(
	privateKey string,
	amount *big.Int,
	remoteAddrStr string,
	incTokenId string,
) (*entity.BurningForDepositToSCRes, error) {
	burningAddress, err := w.getBurningAddressFromChain()
	if err != nil {
		return nil, errors.Wrapf(err, "w.blockchainAPI: method %+v, Get burn address", constant.CreateAndSendBurningForDepositToSCRequest)
	}

	param := []interface{}{
		privateKey,
		nil,
		5,
		-1,
		map[string]interface{}{
			"TokenID":     incTokenId,
			"TokenTxType": 1,
			"TokenName":   "",
			"TokenSymbol": "",
			"TokenAmount": amount.Uint64(),
			"TokenReceivers": map[string]uint64{
				burningAddress: amount.Uint64(), //return pETH, burn token for nobody can not used
			},
			"RemoteAddress": remoteAddrStr, //receive ETH
			"Privacy":       true,
			"TokenFee":      uint64(0),
		},
		"",
		0,
	}

	//rpc: CreateAndSendBurningForDepositToSCRequest
	rawData, err := w.IncChainIntegration.CreateAndSendBurningForDepositToSCRequest(param)
	if err != nil {
		return nil, errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method CreateAndSendBurningForDepositToSCRequest: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawPrivacyCustomTokenTransaction, rawData)

	if err != nil {
		return nil, errors.Wrapf(err, "w.blockchainAPI: method %+v", constant.CreateAndSendBurningForDepositToSCRequest)
	}

	data := resp.(map[string]interface{})
	resultResp := data["Result"]

	if resultResp == nil {
		return nil, errors.Errorf("couldn't get txID: resp: %+v", data)
	}

	resultInBytes, err := json.MarshalIndent(resultResp, "", "\t")
	if err != nil {
		return nil, err
	}

	result := entity.BurningForDepositToSCRes{}
	err = json.Unmarshal(resultInBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (w *Wallet) GetBalance(privateKey string, tokenId string) (uint64, error) {
	var balance uint64
	var bigBalance *big.Int
	var err error

	if tokenId == w.ConstantID {
		balance, err = w.GetBalanceByPrivateKey(privateKey)
	} else {
		bigBalance, err = w.GetListPrivacyCustomTokenBalanceByID(privateKey, tokenId)
		if err == nil {
			balance = bigBalance.Uint64()
		}
	}

	return balance, err
}

func (w *Wallet) SellPRV(privateKey string, buyTokenId string, tradingFee uint64, sellTokenAmount uint64, minimumAmount uint64, traderAddress string) (string, error) {
	var burningAddress string
	var err error
	burningAddress, err = w.Block.GetBurningAddress()

	if err != nil {
		return "", errors.Wrap(err, "w.GetBurningAddress")
	}

	metadata := map[string]interface{}{
		"TokenIDToBuyStr":     buyTokenId,
		"TokenIDToSellStr":    w.ConstantID,
		"SellAmount":          sellTokenAmount,
		"MinAcceptableAmount": minimumAmount,
		"TradingFee":          tradingFee,
		"TraderAddressStr":    traderAddress,
	}

	paramArray := []interface{}{
		privateKey,
		map[string]interface{}{
			burningAddress: sellTokenAmount + tradingFee,
		},
		1,
		-1,
		metadata,
	}

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.CreateAndSendTxWithPRVTradeReq, paramArray)
	if err != nil {
		return "", errors.Wrapf(err, "w.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: req: %+v, resp: %+v", metadata, data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	if result["TxID"] == nil {
		return "", constant.ErrTxHashNotExists
	}
	return result["TxID"].(string), nil
}

func (w *Wallet) SellPToken(privateKey string, buyTokenId string, tradingFee uint64, sellTokenId string, sellTokenAmount uint64, minimumAmount uint64, traderAddress string, networkFeeTokenID string, networkFee uint64) (string, error) {
	var burningAddress string
	var err error
	burningAddress, err = w.Block.GetBurningAddress()
	var FeePerKb int
	var TokenFee uint64
	if err != nil {
		return "", errors.Wrap(err, "w.GetBurningAddress")
	}

	if networkFeeTokenID == w.ConstantID {
		FeePerKb = 5
		TokenFee = 0
	} else {
		FeePerKb = 0
		TokenFee = networkFee / constant.PDEX_TRADE_STEPS
	}

	metadata := map[string]interface{}{
		"Privacy":     true,
		"TokenID":     sellTokenId,
		"TokenTxType": 1,
		"TokenName":   "",
		"TokenSymbol": "",
		"TokenAmount": sellTokenAmount,
		"TokenReceivers": map[string]interface{}{
			burningAddress: sellTokenAmount + tradingFee,
		},
		"TokenFee":            TokenFee,
		"TokenIDToBuyStr":     buyTokenId,
		"TokenIDToSellStr":    sellTokenId,
		"SellAmount":          sellTokenAmount,
		"MinAcceptableAmount": minimumAmount,
		"TradingFee":          tradingFee,
		"TraderAddressStr":    traderAddress,
	}

	paramArray := []interface{}{
		privateKey,
		nil,
		FeePerKb,
		-1,
		metadata,
		"",
		0,
	}

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.CreateAndSendTxWithPTokenTradeReq, paramArray)
	if err != nil {
		return "", errors.Wrapf(err, "w.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: req: %+v, resp: %+v", metadata, data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	if result["TxID"] == nil {
		return "", constant.ErrTxHashNotExists
	}

	return result["TxID"].(string), nil
}

func (w *Wallet) GetTransactionAmount(txId string, walletAddress string, readOnlyKey string) (uint64, error) {
	receiveDetail, err := w.GetDecryptOutputCoinByKeyOfTransaction(txId, walletAddress, readOnlyKey)
	if err != nil {
		return 0, errors.Wrap(err, "w.GetDecryptOutputCoinByKeyOfTransaction")
	}

	return receiveDetail.AmountPRV, nil
}

func (w *Wallet) SendToken(privateKey string, receiverAddress string, tokenId string, amount uint64, fee uint64, feeTokenId string) (string, error) {
	if tokenId == w.ConstantID {
		var listPaymentAddresses = make(map[string]uint64)
		listPaymentAddresses[receiverAddress] = amount
		return w.createAndSendConstantPrivacyTransaction(privateKey, entity.WalletSend{
			Type:             0,
			PaymentAddresses: listPaymentAddresses,
		})
	}

	param := entity.WalletSend{
		TokenID:     tokenId,
		Type:        1,
		TokenName:   "",
		TokenSymbol: "",
		PaymentAddresses: map[string]uint64{
			receiverAddress: amount,
		},
	}

	//fee pay by token
	if feeTokenId == tokenId {
		param.TokenFee = fee
	}

	tx, err := w.sendPrivacyCustomTokenTransaction(privateKey, param)

	if err != nil {
		return "", errors.Wrap(err, "p.SendPrivacyCustomTokenTransaction")
	}

	txID, ok := tx["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", tx)
	}
	return txID, nil
}

func (w *Wallet) DefragmentationPrv(privateKey string, maxValue int64) (string, error) {
	param := []interface{}{
		privateKey,
		maxValue,
		constant.EstimateFee,
		0,
	}

	//rpc: defragmentaccount
	rawData, err := w.IncChainIntegration.DefragmentationPrv(param)
	if err != nil {
		return "", errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method DefragmentationPrv: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)

	if err != nil {
		return "", errors.Wrap(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (w *Wallet) DefragmentationPToken(privateKey string, tokenId string) (string, error) {
	tokenData := map[string]interface{}{}
	tokenData["Privacy"] = true
	tokenData["TokenID"] = tokenId
	tokenData["TokenName"] = ""
	tokenData["TokenSymbol"] = ""
	tokenData["TokenTxType"] = 1
	tokenData["TokenReceivers"] = map[string]uint64{}
	tokenData["TokenAmount"] = uint64(0)
	tokenData["TokenFee"] = uint64(0)

	object := map[string]uint64{}
	nativeFee := -1

	params := []interface{}{
		privateKey,
		object,
		nativeFee,
		1,
		tokenData,
		"",
		1,
	}

	//rpc: defragmentaccounttoken
	rawData, err := w.IncChainIntegration.DefragmentationPToken(params)
	if err != nil {
		return "", errors.Wrap(err, "w.IncChainIntegration")
	}

	fmt.Printf("raw data method DefragmentationPToken: %v \n", rawData)

	resp, _, err := w.Inc.PostAndReceiveInterface(constant.SendRawPrivacyCustomTokenTransaction, rawData)

	if err != nil {
		return "", errors.Wrap(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (w *Wallet) GetUTXO(privateKey string, tokenId string) ([]*entity.Utxo, error) {
	var input []*entity.Utxo

	inputCoin, err := w.IncChainIntegration.GetUTXO(privateKey, tokenId)
	if err != nil {
		return nil, err
	}

	for _, value := range inputCoin {
		input = append(input, &entity.Utxo{
			Value:        value.CoinDetails.GetValue(),
			SerialNumber: string(value.CoinDetails.GetSerialNumber().MarshalText()),
			SnDerivator:  value.CoinDetails.GetSNDerivator().String(),
		})
	}

	return input, nil
}
