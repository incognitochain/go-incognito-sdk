package bean

import (
	"errors"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/incognitokey"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/wallet"
)

type CreateRawTxParam struct {
	SenderKeySet         *incognitokey.KeySet
	ShardIDSender        byte
	PaymentInfos         []*privacy.PaymentInfo
	EstimateFeeCoinPerKb int64
	HasPrivacyCoin       bool
	Info                 []byte
}

func GetKeySetFromPrivateKeyParams(privateKeyWalletStr string) (*incognitokey.KeySet, byte, error) {
	// deserialize to crate keywallet object which contain private key
	keyWallet, err := wallet.Base58CheckDeserialize(privateKeyWalletStr)
	if err != nil {
		return nil, byte(0), err
	}

	// fill paymentaddress and readonly key with privatekey
	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	if err != nil {
		return nil, byte(0), err
	}

	if len(keyWallet.KeySet.PaymentAddress.Pk) == 0 {
		return nil, byte(0), errors.New("private key is not valid")
	}

	// calculate shard ID
	lastByte := keyWallet.KeySet.PaymentAddress.Pk[len(keyWallet.KeySet.PaymentAddress.Pk)-1]
	shardID := common.GetShardIDFromLastByte(lastByte)

	return &keyWallet.KeySet, shardID, nil
}

func GetPrivateKey(params interface{}) (*wallet.KeyWallet, error) {
	arrayParams := common.InterfaceSlice(params)
	if len(arrayParams) < 3 {
		return nil, errors.New("not enough param")
	}

	// param #1: private key of sender
	senderKeyParam, ok := arrayParams[0].(string)
	if !ok {
		return nil, errors.New("sender private key is invalid")
	}

	//generate key wallet
	keyWallet, err := wallet.Base58CheckDeserialize(senderKeyParam)
	if err != nil {
		return nil, err
	}

	// fill paymentaddress and readonly key with privatekey
	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	if err != nil {
		return nil, err
	}

	return keyWallet, nil
}

func NewCreateRawTxParam(params interface{}) (*CreateRawTxParam, error) {
	arrayParams := common.InterfaceSlice(params)
	if len(arrayParams) < 3 {
		return nil, errors.New("not enough param")
	}

	// param #1: private key of sender
	senderKeyParam, ok := arrayParams[0].(string)
	if !ok {
		return nil, errors.New("sender private key is invalid")
	}

	//generate key wallet
	keyWallet, err := wallet.Base58CheckDeserialize(senderKeyParam)
	if err != nil {
		return nil, err
	}

	// fill paymentaddress and readonly key with privatekey
	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	if err != nil {
		return nil, err
	}
	//end

	senderKeySet, shardIDSender, err := GetKeySetFromPrivateKeyParams(senderKeyParam)
	if err != nil {
		return nil, err
	}

	// param #2: list receivers
	receivers := make(map[string]uint64)
	if arrayParams[1] != nil {
		receivers, ok = arrayParams[1].(map[string]uint64)
		if !ok {
			return nil, errors.New("receivers param is invalid")
		}
	}

	paymentInfos := make([]*privacy.PaymentInfo, 0)
	for paymentAddressStr, amount := range receivers {
		keyWalletReceiver, err := wallet.Base58CheckDeserialize(paymentAddressStr)
		if err != nil {
			return nil, err
		}
		if len(keyWalletReceiver.KeySet.PaymentAddress.Pk) == 0 {
			return nil, fmt.Errorf("payment info %+v is invalid", paymentAddressStr)
		}

		paymentInfo := &privacy.PaymentInfo{
			Amount:         amount,
			PaymentAddress: keyWalletReceiver.KeySet.PaymentAddress,
		}
		paymentInfos = append(paymentInfos, paymentInfo)
	}

	// param #3: estimation fee nano P per kb
	estimateFeeCoinPerKb, ok := arrayParams[2].(int)
	if !ok {
		return nil, errors.New("estimate fee coin per kb is invalid")
	}

	// param #4: hasPrivacyCoin flag: 1 or -1
	// default: -1 (has no privacy) (if missing this param)
	hasPrivacyCoinParam := int(-1)
	if len(arrayParams) > 3 {
		hasPrivacyCoinParam, ok = arrayParams[3].(int)
		if !ok {
			return nil, errors.New("has privacy for tx is invalid")
		}
	}
	hasPrivacyCoin := int(hasPrivacyCoinParam) > 0

	// param #5: meta data (optional)
	// don't do anything

	// param#6: info (optional)
	info := []byte{}
	if len(arrayParams) > 5 {
		if arrayParams[5] != nil {
			infoStr, ok := arrayParams[5].(string)
			if !ok {
				return nil, errors.New("info is invalid")
			}
			info = []byte(infoStr)
		}

	}

	return &CreateRawTxParam{
		SenderKeySet:         senderKeySet,
		ShardIDSender:        shardIDSender,
		PaymentInfos:         paymentInfos,
		EstimateFeeCoinPerKb: int64(estimateFeeCoinPerKb),
		HasPrivacyCoin:       hasPrivacyCoin,
		Info:                 info,
	}, nil
}
