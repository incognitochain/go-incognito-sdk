package incognito

import (
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/transaction"
	"github.com/incognitochain/go-incognito-sdk/wallet"
)

func GetBalance(rpcClient *rpcclient.HttpClient, privateKey string, tokenId string) (uint64, error) {
	inputCoin, err := getUnspentOutputCoinsExceptSpendingUTXO(rpcClient, privateKey, tokenId)
	if err != nil {
		return 0, err
	}
	var accountBalance uint64

	for _, value := range inputCoin {
		accountBalance = accountBalance + value.CoinDetails.GetValue()
	}

	return accountBalance, nil
}


// GetUnspentOutputCoins return utxos of an account
func getUnspentOutputCoinsExceptSpendingUTXO(rpcClient *rpcclient.HttpClient, privateKey string, tokenId string) ([]*privacy.InputCoin, error) {
	keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return nil, fmt.Errorf("Can not deserialize priavte key %v\n", err)
	}
	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	if err != nil {
		return nil, errors.New("sender private key is invalid")
	}

	tokenID, err := common.Hash{}.NewHashFromStr(tokenId)
	if err != nil {
		return nil, err
	}

	// get unspent output coins from network
	utxos, err := rpcclient.GetUnspentOutputCoins(rpcClient, keyWallet, tokenID)
	if err != nil {
		return nil, err
	}

	inputCoins := transaction.ConvertOutputCoinToInputCoin(utxos)
	return inputCoins, nil
}