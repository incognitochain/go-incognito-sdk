package rpcclient

import (
	"errors"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"strconv"
)

func GetEstimateFeeWithEstimator(rpcClient *HttpClient, defaultFee int64, paymentAddrSerialize string, tokenIDStr *common.Hash) (uint64, error) {
	var estimateFees EstimateFeeRes

	params := []interface{}{
		defaultFee,
		paymentAddrSerialize,
		8,
		tokenIDStr,
	}

	err := rpcClient.RPCCall("estimatefeewithestimator", params, &estimateFees)
	if err != nil {
		return 0, err
	}

	if estimateFees.RPCError != nil {
		return 0, errors.New(estimateFees.RPCError.StackTrace)
	}

	return estimateFees.Result.EstimateFeeCoinPerKb, nil
}

// GetUnspentOutputCoins return utxos of an account
func GetUnspentOutputCoins(rpcClient *HttpClient, keyWallet *wallet.KeyWallet, tokenId *common.Hash) ([]*privacy.OutputCoin, error) {
	privateKey := &keyWallet.KeySet.PrivateKey
	paymentAddressStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)
	viewingKeyStr := keyWallet.Base58CheckSerialize(wallet.ReadonlyKeyType)

	outputCoins, err := getListOutputCoins(rpcClient, paymentAddressStr, viewingKeyStr, tokenId)
	if err != nil {
		return nil, err
	}

	serialNumbers, err := deriveSerialNumbers(privateKey, outputCoins)
	if err != nil {
		return nil, err
	}

	isExisted, err := checkExistenceSerialNumber(rpcClient, paymentAddressStr, serialNumbers, tokenId)
	if err != nil {
		return nil, err
	}

	utxos := make([]*privacy.OutputCoin, 0)
	for i, out := range outputCoins {
		if !isExisted[i] {
			utxos = append(utxos, out)
		}
	}

	return utxos, nil
}

// GetListOutputCoins calls Incognito RPC to get all output coins of the account
func getListOutputCoins(rpcClient *HttpClient, paymentAddress string, viewingKey string, tokenId *common.Hash) ([]*privacy.OutputCoin, error) {
	var outputCoinsRes ListOutputCoinsRes
	params := []interface{}{
		0,
		999999,
		[]map[string]string{
			{
				"PaymentAddress": paymentAddress,
				"ReadonlyKey":    viewingKey,
			},
		},
	}

	if len(tokenId.String()) > 0 {
		params = append(params, tokenId.String())
	}

	err := rpcClient.RPCCall("listoutputcoins", params, &outputCoinsRes)
	if err != nil {
		return nil, err
	}

	if outputCoinsRes.RPCError != nil {
		return nil, errors.New(outputCoinsRes.RPCError.StackTrace)
	}

	outputCoins, err := newOutputCoinsFromResponse(outputCoinsRes.Result.Outputs[viewingKey])
	if err != nil {
		return nil, err
	}
	return outputCoins, nil
}

func deriveSerialNumbers(privateKey *privacy.PrivateKey, outputCoins []*privacy.OutputCoin) ([]*privacy.Point, error) {
	serialNumbers := make([]*privacy.Point, len(outputCoins))
	for i, coin := range outputCoins {
		coin.CoinDetails.SetSerialNumber(
			new(privacy.Point).Derive(
				privacy.PedCom.G[privacy.PedersenPrivateKeyIndex],
				new(privacy.Scalar).FromBytesS(*privateKey),
				coin.CoinDetails.GetSNDerivator()))
		serialNumbers[i] = coin.CoinDetails.GetSerialNumber()
	}

	return serialNumbers, nil
}

// CheckExistenceSerialNumber calls Incognito RPC to check existence serial number on network
// to check output coins is spent or unspent
func checkExistenceSerialNumber(rpcClient *HttpClient, paymentAddressStr string, sns []*privacy.Point, tokenId *common.Hash) ([]bool, error) {
	var hasSerialNumberRes HasSerialNumberRes
	result := make([]bool, 0)
	snStrs := make([]interface{}, len(sns))
	for i, sn := range sns {
		snStrs[i] = base58.Base58Check{}.Encode(sn.ToBytesS(), common.Base58Version)
	}

	params := []interface{}{
		paymentAddressStr,
		snStrs,
		tokenId.String(),
	}
	err := rpcClient.RPCCall("hasserialnumbers", params, &hasSerialNumberRes)
	if err != nil {
		return nil, err
	}

	if hasSerialNumberRes.RPCError != nil {
		return nil, errors.New(hasSerialNumberRes.RPCError.StackTrace)
	}

	result = hasSerialNumberRes.Result
	return result, nil
}

func newOutputCoinsFromResponse(outCoins []OutCoin) ([]*privacy.OutputCoin, error) {
	outputCoins := make([]*privacy.OutputCoin, len(outCoins))
	for i, outCoin := range outCoins {
		outputCoins[i] = new(privacy.OutputCoin).Init()
		publicKey, _, _ := base58.Base58Check{}.Decode(outCoin.PublicKey)
		publicKeyPoint, _ := new(privacy.Point).FromBytesS(publicKey)
		outputCoins[i].CoinDetails.SetPublicKey(publicKeyPoint)

		cmBytes, _, _ := base58.Base58Check{}.Decode(outCoin.CoinCommitment)
		cmPoint, _ := new(privacy.Point).FromBytesS(cmBytes)
		outputCoins[i].CoinDetails.SetCoinCommitment(cmPoint)

		sndBytes, _, _ := base58.Base58Check{}.Decode(outCoin.SNDerivator)
		sndScalar := new(privacy.Scalar).FromBytesS(sndBytes)
		outputCoins[i].CoinDetails.SetSNDerivator(sndScalar)

		randomnessBytes, _, _ := base58.Base58Check{}.Decode(outCoin.Randomness)
		randomnessScalar := new(privacy.Scalar).FromBytesS(randomnessBytes)
		outputCoins[i].CoinDetails.SetRandomness(randomnessScalar)

		value, _ := strconv.Atoi(outCoin.Value)
		outputCoins[i].CoinDetails.SetValue(uint64(value))
	}

	return outputCoins, nil
}


func newOutCoin(outCoin *privacy.OutputCoin) OutCoin {
	serialNumber := ""

	if outCoin.CoinDetails.GetSerialNumber() != nil && !outCoin.CoinDetails.GetSerialNumber().IsIdentity() {
		serialNumber = base58.Base58Check{}.Encode(outCoin.CoinDetails.GetSerialNumber().ToBytesS(), common.ZeroByte)
	}

	result := OutCoin{
		PublicKey:      base58.Base58Check{}.Encode(outCoin.CoinDetails.GetPublicKey().ToBytesS(), common.ZeroByte),
		Value:          strconv.FormatUint(outCoin.CoinDetails.GetValue(), 10),
		Info:           base58.Base58Check{}.Encode(outCoin.CoinDetails.GetInfo()[:], common.ZeroByte),
		CoinCommitment: base58.Base58Check{}.Encode(outCoin.CoinDetails.GetCoinCommitment().ToBytesS(), common.ZeroByte),
		SNDerivator:    base58.Base58Check{}.Encode(outCoin.CoinDetails.GetSNDerivator().ToBytesS(), common.ZeroByte),
		SerialNumber:   serialNumber,
	}

	if outCoin.CoinDetails.GetRandomness() != nil {
		result.Randomness = base58.Base58Check{}.Encode(outCoin.CoinDetails.GetRandomness().ToBytesS(), common.ZeroByte)
	}
	// return more data of CoinDetailsEncrypted
	if outCoin.CoinDetailsEncrypted != nil {
		result.CoinDetailsEncrypted = base58.Base58Check{}.Encode(outCoin.CoinDetailsEncrypted.Bytes(), common.ZeroByte)
	}

	return result
}

func RandomCommitmentsProcess(rpcClient *HttpClient, outputCoins []*privacy.OutputCoin, paymentAddrStr string, tokenID *common.Hash) ([]uint64, []uint64, []string, error) {
	var randomCommitmentRes RandomCommitmentRes

	item := make([]OutCoin, 0)
	for _, outCoin := range outputCoins {
		item = append(item, newOutCoin(outCoin))
	}

	params := []interface{}{
		paymentAddrStr,
		item,
		tokenID.String(),
	}

	err := rpcClient.RPCCall("randomcommitments", params, &randomCommitmentRes)
	if err != nil {
		return nil, nil, nil, err
	}

	if randomCommitmentRes.RPCError != nil {
		return nil, nil, nil, errors.New(randomCommitmentRes.RPCError.StackTrace)
	}

	return randomCommitmentRes.Result.CommitmentIndices, randomCommitmentRes.Result.MyCommitmentIndexs, randomCommitmentRes.Result.Commitments, nil
}

func CheckSNDerivatorExistence(rpcClient *HttpClient, paymentAddressStr string, sndOut []*privacy.Scalar) ([]bool, error) {
	var hasSNDerivatorRes HasSNDerivatorRes
	sndStrs := make([]interface{}, len(sndOut))
	for i, sn := range sndOut {
		sndStrs[i] = base58.Base58Check{}.Encode(sn.ToBytesS(), common.Base58Version)
	}
	params := []interface{}{
		paymentAddressStr,
		sndStrs,
	}
	err := rpcClient.RPCCall("hassnderivators", params, &hasSNDerivatorRes)
	if err != nil {
		return nil, err
	}

	if hasSNDerivatorRes.RPCError != nil {
		return nil, errors.New(hasSNDerivatorRes.RPCError.StackTrace)
	}

	return hasSNDerivatorRes.Result, nil
}