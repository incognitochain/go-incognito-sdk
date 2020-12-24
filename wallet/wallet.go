package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"

	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
)

type AccountWallet struct {
	Name       string
	Key        KeyWallet
	Child      []AccountWallet
	IsImported bool
}

type WalletConfig struct {
	DataDir        string
	DataFile       string
	DataPath       string
	IncrementalFee uint64
	ShardID        *byte //default is nil -> create account for any shard
}

type Wallet struct {
	Seed          []byte
	Entropy       []byte
	PassPhrase    string
	Mnemonic      string
	MasterAccount AccountWallet
	Name          string
	config        *WalletConfig
}

// GetConfig returns configuration of wallet
func (wallet Wallet) GetConfig() *WalletConfig {
	return wallet.config
}

func CreateNewAccount() (*AccountWallet, error) {
	mnemonicGen := MnemonicGenerator{}
	entropy, _ := mnemonicGen.newEntropy(128)
	mnemonic, _ := mnemonicGen.newMnemonic(entropy)
	seed := mnemonicGen.NewSeed(mnemonic, "")

	masterKey, err := NewMasterKey(seed)

	if err != nil {
		return &AccountWallet{}, err
	}

	masterAccount := AccountWallet{
		Key:   *masterKey,
		Child: make([]AccountWallet, 0),
		Name:  "master",
	}

	childKey, _ := masterAccount.Key.NewChildKey(0)
	account := AccountWallet{
		Key:   *childKey,
		Child: make([]AccountWallet, 0),
		Name:  fmt.Sprintf("Child%v", rand.Intn(10000000)),
	}

	lastByte := childKey.KeySet.PaymentAddress.Pk[len(childKey.KeySet.PaymentAddress.Pk)-1]
	shardId := common.GetShardIDFromLastByte(lastByte)

	fmt.Println(fmt.Sprintf("Generating wallet with shardId %v and Index %v", shardId, 0))

	masterAccount.Child = append(masterAccount.Child, account)

	return &account, nil
}

func CreateNewAccountByShardId(shardId int) (*AccountWallet, error) {
	mnemonicGen := MnemonicGenerator{}
	entropy, _ := mnemonicGen.newEntropy(128)
	mnemonic, _ := mnemonicGen.newMnemonic(entropy)
	seed := mnemonicGen.NewSeed(mnemonic, "")

	masterKey, err := NewMasterKey(seed)

	if err != nil {
		return &AccountWallet{}, err
	}

	masterAccount := AccountWallet{
		Key:   *masterKey,
		Child: make([]AccountWallet, 0),
		Name:  "master",
	}

	newIndex := uint64(0)
	shardIDByte := byte(shardId)

	// loop to get create a new child which can be equal shardID param
	var childKey *KeyWallet
	for true {
		childKey, _ = masterAccount.Key.NewChildKey(uint32(newIndex))
		lastByte := childKey.KeySet.PaymentAddress.Pk[len(childKey.KeySet.PaymentAddress.Pk)-1]
		if common.GetShardIDFromLastByte(lastByte) == shardIDByte {
			fmt.Println(fmt.Sprintf("Generating wallet with shardId %v and Index %v", shardId, newIndex))
			break
		}
		newIndex += 1
	}

	if childKey == nil {
		return nil, errors.New("ChildKey not found")
	}

	account := AccountWallet{
		Key:   *childKey,
		Child: make([]AccountWallet, 0),
		Name:  fmt.Sprintf("Child%v", rand.Intn(10000000)),
	}

	masterAccount.Child = append(masterAccount.Child, account)

	return &account, nil
}

func CreateImportMasterAccount(mnemonic, passPhrase string) (*Wallet, *KeySerializedData, error) {

	wallets := &Wallet{}

	mnemonicGen := MnemonicGenerator{}

	if len(mnemonic) == 0 {
		entropy, _ := mnemonicGen.newEntropy(128)
		mnemonic, _ = mnemonicGen.newMnemonic(entropy)
		wallets.Entropy = entropy
	}

	seed := mnemonicGen.NewSeed(mnemonic, passPhrase)

	masterKey, err := NewMasterKey(seed)

	if err != nil {
		return nil, nil, err
	}

	masterAccount := AccountWallet{
		Key:   *masterKey,
		Child: make([]AccountWallet, 0),
		Name:  "master",
	}

	childKey, _ := masterAccount.Key.NewChildKey(0)

	wallets.Mnemonic = mnemonic
	wallets.Seed = seed
	wallets.PassPhrase = passPhrase
	wallets.Name = "master account"

	lastByte := childKey.KeySet.PaymentAddress.Pk[len(childKey.KeySet.PaymentAddress.Pk)-1]
	shardId := common.GetShardIDFromLastByte(lastByte)

	key := &KeySerializedData{
		PaymentAddress: childKey.Base58CheckSerialize(PaymentAddressType),
		Pubkey:         hex.EncodeToString(childKey.KeySet.PaymentAddress.Pk),
		ReadonlyKey:    childKey.Base58CheckSerialize(ReadonlyKeyType),
		PrivateKey:     childKey.Base58CheckSerialize(PriKeyType),
		ValidatorKey:   base58.Base58Check{}.Encode(common.HashB(common.HashB(childKey.KeySet.PrivateKey)), common.ZeroByte),
		ShardId:        int(shardId),
	}

	return wallets, key, nil
}

func ImportAccount(privateKeyStr string, accountName string) (*AccountWallet, error) {
	keyWallet, err := Base58CheckDeserialize(privateKeyStr)
	if err != nil {
		return nil, err
	}

	err = keyWallet.KeySet.InitFromPrivateKey(&keyWallet.KeySet.PrivateKey)
	if err != nil {
		return nil, err
	}

	masterAccount := &AccountWallet{}

	account := AccountWallet{
		Key:        *keyWallet,
		Child:      make([]AccountWallet, 0),
		IsImported: true,
		Name:       accountName,
	}
	masterAccount.Child = append(masterAccount.Child, account)

	return &account, nil
}
