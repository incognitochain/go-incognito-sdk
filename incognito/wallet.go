package incognito

import (
	"encoding/hex"

	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/wallet"
)

func CreateNewWallet() (*wallet.KeySerializedData, error) {
	newAccount, err := wallet.CreateNewAccount()

	if err != nil {
		return &wallet.KeySerializedData{}, err
	}

	lastByte := newAccount.Key.KeySet.PaymentAddress.Pk[len(newAccount.Key.KeySet.PaymentAddress.Pk)-1]
	shardId := common.GetShardIDFromLastByte(lastByte)

	key := &wallet.KeySerializedData{
		PaymentAddress: newAccount.Key.Base58CheckSerialize(wallet.PaymentAddressType),
		Pubkey:         hex.EncodeToString(newAccount.Key.KeySet.PaymentAddress.Pk),
		ReadonlyKey:    newAccount.Key.Base58CheckSerialize(wallet.ReadonlyKeyType),
		PrivateKey:     newAccount.Key.Base58CheckSerialize(wallet.PriKeyType),
		ValidatorKey:   base58.Base58Check{}.Encode(common.HashB(common.HashB(newAccount.Key.KeySet.PrivateKey)), common.ZeroByte),
		ShardId:        int(shardId),
	}

	return key, nil
}

func CreateNewWalletByShardId(shardId int) (*wallet.KeySerializedData, error) {
	newAccount, err := wallet.CreateNewAccountByShardId(shardId)

	if err != nil {
		return &wallet.KeySerializedData{}, err
	}

	key := &wallet.KeySerializedData{
		PaymentAddress: newAccount.Key.Base58CheckSerialize(wallet.PaymentAddressType),
		Pubkey:         hex.EncodeToString(newAccount.Key.KeySet.PaymentAddress.Pk),
		ReadonlyKey:    newAccount.Key.Base58CheckSerialize(wallet.ReadonlyKeyType),
		PrivateKey:     newAccount.Key.Base58CheckSerialize(wallet.PriKeyType),
		ValidatorKey:   base58.Base58Check{}.Encode(common.HashB(common.HashB(newAccount.Key.KeySet.PrivateKey)), common.ZeroByte),
		ShardId:        shardId,
	}

	return key, nil
}

func ImportNewWallet(privateKeyStr string, accountName string) (*wallet.KeySerializedData, error) {
	newAccount, err := wallet.ImportAccount(privateKeyStr, accountName)

	if err != nil {
		return &wallet.KeySerializedData{}, err
	}

	lastByte := newAccount.Key.KeySet.PaymentAddress.Pk[len(newAccount.Key.KeySet.PaymentAddress.Pk)-1]
	shardId := common.GetShardIDFromLastByte(lastByte)

	key := &wallet.KeySerializedData{
		PaymentAddress: newAccount.Key.Base58CheckSerialize(wallet.PaymentAddressType),
		Pubkey:         hex.EncodeToString(newAccount.Key.KeySet.PaymentAddress.Pk),
		ReadonlyKey:    newAccount.Key.Base58CheckSerialize(wallet.ReadonlyKeyType),
		PrivateKey:     newAccount.Key.Base58CheckSerialize(wallet.PriKeyType),
		ValidatorKey:   base58.Base58Check{}.Encode(common.HashB(common.HashB(newAccount.Key.KeySet.PrivateKey)), common.ZeroByte),
		ShardId:        int(shardId),
	}

	return key, err
}
