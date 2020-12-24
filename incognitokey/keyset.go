package incognitokey

import (
	"errors"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/privacy"
)

// KeySet is real raw data of wallet account, which user can use to
// - spend and check double spend coin with private key
// - receive coin with payment address
// - read tx data with readonly key
type KeySet struct {
	PrivateKey     privacy.PrivateKey
	PaymentAddress privacy.PaymentAddress
	ReadonlyKey    privacy.ViewingKey
}

// GenerateKey generates key set from seed in byte array
func (keySet *KeySet) GenerateKey(seed []byte) *KeySet {
	keySet.PrivateKey = privacy.GeneratePrivateKey(seed)
	keySet.PaymentAddress = privacy.GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = privacy.GenerateViewingKey(keySet.PrivateKey[:])
	return keySet
}

// InitFromPrivateKeyByte receives private key in bytes array,
// and regenerates payment address and readonly key
// returns error if private key is invalid
func (keySet *KeySet) InitFromPrivateKeyByte(privateKey []byte) error {
	if len(privateKey) != common.PrivateKeySize {
		return errors.New("invalid size of private key")
	}

	keySet.PrivateKey = privateKey
	keySet.PaymentAddress = privacy.GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = privacy.GenerateViewingKey(keySet.PrivateKey[:])
	return nil
}

// InitFromPrivateKey receives private key in PrivateKey type,
// and regenerates payment address and readonly key
// returns error if private key is invalid
func (keySet *KeySet) InitFromPrivateKey(privateKey *privacy.PrivateKey) error {
	if privateKey == nil || len(*privateKey) != common.PrivateKeySize {
		return errors.New("invalid size of private key")
	}

	keySet.PrivateKey = *privateKey
	keySet.PaymentAddress = privacy.GeneratePaymentAddress(keySet.PrivateKey[:])
	keySet.ReadonlyKey = privacy.GenerateViewingKey(keySet.PrivateKey[:])

	return nil
}