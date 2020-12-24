package incognito

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/incognitokey"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/rpcservice"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"github.com/pkg/errors"
)

func handleCreateRawStakingTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	paramsArray := common.InterfaceSlice(params)
	if paramsArray == nil || len(paramsArray) < 5 {
		return nil, errors.New("param must be an array at least 5 element")
	}

	createRawTxParam, errNewParam := bean.NewCreateRawTxParam(params)
	if errNewParam != nil {
		return nil, errNewParam
	}

	keyWallet := new(wallet.KeyWallet)
	keyWallet.KeySet = *createRawTxParam.SenderKeySet
	funderPaymentAddress := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)

	// prepare meta data
	data, ok := paramsArray[4].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid Data For Staking Transaction %+v", paramsArray[4])
	}

	stakingType, ok := data["StakingType"].(int)
	if !ok {
		return nil, fmt.Errorf("Invalid Staking Type For Staking Transaction %+v", data["StakingType"])
	}

	candidatePaymentAddress, ok := data["CandidatePaymentAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid Producer Payment Address for Staking Transaction %+v", data["CandidatePaymentAddress"])
	}

	// Get private seed, a.k.a mining key
	privateSeed, ok := data["PrivateSeed"].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid Private Seed For Staking Transaction %+v", data["PrivateSeed"])
	}

	privateSeedBytes, ver, errDecode := base58.Base58Check{}.Decode(privateSeed)
	if (errDecode != nil) || (ver != common.ZeroByte) {
		return nil, errors.New("Decode privateseed failed!")
	}

	//Get RewardReceiver Payment Address
	rewardReceiverPaymentAddress, ok := data["RewardReceiverPaymentAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid Reward Receiver Payment Address For Staking Transaction %+v", data["RewardReceiverPaymentAddress"])
	}

	//Get auto staking flag
	autoReStaking, ok := data["AutoReStaking"].(bool)
	if !ok {
		return nil, fmt.Errorf("Invalid auto restaking flag %+v", data["AutoReStaking"])
	}

	// Get candidate publickey
	candidateWallet, err := wallet.Base58CheckDeserialize(candidatePaymentAddress)
	if err != nil || candidateWallet == nil {
		return nil, errors.New("Base58CheckDeserialize candidate Payment Address failed")
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk

	committeePK, err := incognitokey.NewCommitteeKeyFromSeed(privateSeedBytes, pk)
	if err != nil {
		return nil, errors.New("Cannot get payment address")
	}

	committeePKBytes, err := committeePK.Bytes()
	if err != nil {
		return nil, errors.New("Cannot import key set")
	}

	stakingMetadata, err := metadata.NewStakingMetadata(
		stakingType,
		funderPaymentAddress,
		rewardReceiverPaymentAddress,
		1750000000000, //1750e9
		base58.Base58Check{}.Encode(committeePKBytes, common.ZeroByte),
		autoReStaking,
	)

	if err != nil {
		return nil, err
	}

	keyWallet1, err := bean.GetPrivateKey(params)
	if err != nil {
		return nil, err
	}

	txService := &rpcservice.TxService{
		RpcClient: rpcClient,
		KeyWallet: keyWallet1,
	}

	txID, err := txService.BuildRawTransaction(createRawTxParam, stakingMetadata)
	if err != nil {
		return nil, err
	}

	txBytes, err := json.Marshal(txID)
	if err != nil {
		// return hex for a new tx
		return nil, err
	}

	txShardID := common.GetShardIDFromLastByte(txID.GetSenderAddrLastByte())

	result := rpcclient.CreateTransactionResult{
		TxID:            txID.String(),
		Base58CheckData: base58.Base58Check{}.Encode(txBytes, common.ZeroByte),
		ShardID:         txShardID,
	}
	return result, nil
}

func CreateAndSendStakingTx(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	var err error
	data, err := handleCreateRawStakingTransaction(rpcClient, params)
	if err != nil {
		return nil, err
	}

	tx := data.(rpcclient.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	return newParam, nil
	//sendResult, err := httpServer.handleSendRawTransaction(newParam, closeChan)
}
