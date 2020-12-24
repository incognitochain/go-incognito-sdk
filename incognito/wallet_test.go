package incognito

import (
	"encoding/hex"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
	result, err := CreateNewWallet()
	assert.Equal(t, nil, err)
	assert.NotEmpty(t, result.Pubkey)
	assert.NotEmpty(t, result.PaymentAddress)
	assert.NotEmpty(t, result.PrivateKey)
	assert.NotEmpty(t, result.ReadonlyKey)

	fmt.Println("PrivateKey:", result.PrivateKey)
	fmt.Println("ShardId:", result.ShardId)

	result1, err := CreateNewWalletByShardId(6)
	assert.Equal(t, nil, err)
	assert.NotEmpty(t, result1.Pubkey)
	assert.NotEmpty(t, result1.PaymentAddress)
	assert.NotEmpty(t, result1.PrivateKey)
	assert.NotEmpty(t, result1.ReadonlyKey)

	fmt.Println("PrivateKey:", result1.PrivateKey)
	fmt.Println("ShardId:", result1.ShardId)

	importWallet, err := ImportNewWallet(result.PrivateKey, "name")
	assert.Equal(t, nil, err)
	assert.NotEmpty(t, importWallet.Pubkey)
	assert.NotEmpty(t, importWallet.PaymentAddress)
	assert.NotEmpty(t, importWallet.PrivateKey)
	assert.NotEmpty(t, importWallet.ReadonlyKey)

	assert.Equal(t, result.Pubkey, importWallet.Pubkey)
	assert.Equal(t, result.PaymentAddress, importWallet.PaymentAddress)
	assert.Equal(t, result.PrivateKey, importWallet.PrivateKey)
	assert.Equal(t, result.ReadonlyKey, importWallet.ReadonlyKey)
}

func TestCreateNewWallet(t *testing.T) {
	var i int
	for i = 0; i <= 2; i++ {
		newAccount, err := wallet.CreateNewAccount()

		assert.Equal(t, nil, err)
		assert.Equal(t, false, newAccount.IsImported)
		assert.Equal(t, 0, len(newAccount.Child))
		assert.Equal(t, 4, len(newAccount.Key.ChildNumber))
		assert.Equal(t, 32, len(newAccount.Key.ChainCode))
		assert.Equal(t, common.PublicKeySize, len(newAccount.Key.KeySet.PaymentAddress.Pk))
		assert.Equal(t, common.TransmissionKeySize, len(newAccount.Key.KeySet.PaymentAddress.Tk))
		assert.Equal(t, common.PrivateKeySize, len(newAccount.Key.KeySet.PrivateKey))
		assert.Equal(t, common.ReceivingKeySize, len(newAccount.Key.KeySet.ReadonlyKey.Rk))

		paymentAddrSerialized := newAccount.Key.Base58CheckSerialize(wallet.PaymentAddressType)
		privateKeySerialized := newAccount.Key.Base58CheckSerialize(wallet.PriKeyType)

		fmt.Println("New name:", newAccount.Name)
		fmt.Println(paymentAddrSerialized)
		fmt.Println(privateKeySerialized)

		newAccountImport, err := wallet.ImportAccount(privateKeySerialized, "Name"+strconv.Itoa(i))
		fmt.Println("Import name: ", newAccountImport.Name)

		keyWallet, _ := wallet.Base58CheckDeserialize(privateKeySerialized)

		assert.Equal(t, nil, err)
		assert.Equal(t, 0, len(newAccountImport.Child))

		assert.Equal(t, "Name"+strconv.Itoa(i), newAccountImport.Name)
		assert.Equal(t, true, newAccountImport.IsImported)
		assert.Equal(t, 0, len(newAccountImport.Child))
		assert.Equal(t, 4, len(newAccountImport.Key.ChildNumber))
		assert.Equal(t, 32, len(newAccountImport.Key.ChainCode))
		assert.Equal(t, keyWallet.KeySet.PrivateKey, newAccountImport.Key.KeySet.PrivateKey)
	}
}

func TestDuplicateNewWallet(t *testing.T) {
	mapKey := make(map[string]string)
	var i int
	for i = 0; i <= 2; i++ {
		newAccount, err := wallet.CreateNewAccount()

		if err != nil {
			fmt.Println(err)
			panic("error")
		}

		paymentAddrSerialized := newAccount.Key.Base58CheckSerialize(wallet.PaymentAddressType)
		privateKeySerialized := newAccount.Key.Base58CheckSerialize(wallet.PriKeyType)

		if !utf8.Valid([]byte(paymentAddrSerialized)) {
			panic("Private key invalid")
		}

		if !utf8.Valid([]byte(privateKeySerialized)) {
			panic("Private key invalid")
		}

		//fmt.Println("New name:", newAccount.Name)
		//fmt.Println(paymentAddrSerialized)
		//fmt.Println(privateKeySerialized)

		if _, ok := mapKey[privateKeySerialized]; ok {
			panic("Private key duplicate!")
		}

		mapKey[privateKeySerialized] = paymentAddrSerialized
	}
}

func TestCreateNewWalletByShardId(t *testing.T) {
	var i int
	for i = 0; i <= 2; i++ {
		newAccount, err := wallet.CreateNewAccountByShardId(3)
		if err != nil {
			fmt.Println(err)
			panic("error")
		}

		assert.Equal(t, nil, err)
		assert.Equal(t, false, newAccount.IsImported)
		assert.Equal(t, 0, len(newAccount.Child))
		assert.Equal(t, 4, len(newAccount.Key.ChildNumber))
		assert.Equal(t, 32, len(newAccount.Key.ChainCode))
		assert.Equal(t, common.PublicKeySize, len(newAccount.Key.KeySet.PaymentAddress.Pk))
		assert.Equal(t, common.TransmissionKeySize, len(newAccount.Key.KeySet.PaymentAddress.Tk))
		assert.Equal(t, common.PrivateKeySize, len(newAccount.Key.KeySet.PrivateKey))
		assert.Equal(t, common.ReceivingKeySize, len(newAccount.Key.KeySet.ReadonlyKey.Rk))

		paymentAddrSerialized := newAccount.Key.Base58CheckSerialize(wallet.PaymentAddressType)
		privateKeySerialized := newAccount.Key.Base58CheckSerialize(wallet.PriKeyType)

		fmt.Println("New name:", newAccount.Name)
		fmt.Println(paymentAddrSerialized)
		fmt.Println(privateKeySerialized)

		newAccountImport, err := wallet.ImportAccount(privateKeySerialized, "Name"+strconv.Itoa(i))
		fmt.Println("Import name: ", newAccountImport.Name)

		keyWallet, _ := wallet.Base58CheckDeserialize(privateKeySerialized)

		assert.Equal(t, nil, err)
		assert.Equal(t, 0, len(newAccountImport.Child))

		assert.Equal(t, "Name"+strconv.Itoa(i), newAccountImport.Name)
		assert.Equal(t, true, newAccountImport.IsImported)
		assert.Equal(t, 0, len(newAccountImport.Child))
		assert.Equal(t, 4, len(newAccountImport.Key.ChildNumber))
		assert.Equal(t, 32, len(newAccountImport.Key.ChainCode))
		assert.Equal(t, keyWallet.KeySet.PrivateKey, newAccountImport.Key.KeySet.PrivateKey)
	}
}

func TestCreateImportMasterAccount(t *testing.T) {

	wallets, keys, err := wallet.CreateImportMasterAccount("", "")

	if err != nil {
		fmt.Println(err)
		panic("error")
	}
	input := strings.Fields(wallets.Mnemonic)
	assert.Equal(t, 12, len(input)) // 12 words

	mnemonic := wallets.Mnemonic
	privateKey := keys.PrivateKey
	fmt.Println("mnemonic create:", mnemonic)
	fmt.Println("privateKey create:", privateKey)

	walletsImport, keysImport, err := wallet.CreateImportMasterAccount(mnemonic, "")
	if err != nil {
		fmt.Println(err)
		panic("error")
	}
	fmt.Println("mnemonic import:", walletsImport.Mnemonic)
	fmt.Println("privateKey import:", keysImport.PrivateKey)

	assert.Equal(t, privateKey, keysImport.PrivateKey)

}

func TestImportWallet(t *testing.T) {
	importWallet, _ := ImportNewWallet("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", "name")

	fmt.Println(importWallet.Pubkey)

	b, _ := hex.DecodeString(importWallet.Pubkey)
	fmt.Println(base58.EncodeCheck(b))

	fmt.Println(importWallet.ShardId)
}
