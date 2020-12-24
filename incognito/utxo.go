package incognito

import (
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
)

func GetUTXO(rpcClient *rpcclient.HttpClient, privateKey string, tokenId string) ([]*privacy.InputCoin, error) {
	inputCoin, err := getUnspentOutputCoinsExceptSpendingUTXO(rpcClient, privateKey, tokenId)
	if err != nil {
		return nil, err
	}

	return inputCoin, nil
}