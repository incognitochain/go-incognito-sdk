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

func handleCreateRawStopAutoStakingTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	// get component
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
	_ = funderPaymentAddress

	//Get data to create meta data
	data, ok := paramsArray[4].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid Params %+v", paramsArray[4])
	}

	//Get staking type
	stopAutoStakingType, ok := data["StopAutoStakingType"].(int)
	if !ok {
		return nil, fmt.Errorf("Invalid Staking Type For Staking Transaction %+v", data["StopAutoStakingType"])
	}

	//Get Candidate Payment Address
	candidatePaymentAddress, ok := data["CandidatePaymentAddress"].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid Producer Payment Address for Staking Transaction %+v", data["CandidatePaymentAddress"])
	}
	// Get private seed, a.k.a mining key
	privateSeed, ok := data["PrivateSeed"].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid Private Seed for Staking Transaction %+v", data["PrivateSeed"])
	}
	privateSeedBytes, ver, err := base58.Base58Check{}.Decode(privateSeed)
	if (err != nil) || (ver != common.ZeroByte) {
		return nil, errors.New("Decode privateseed failed!")
	}

	// Get candidate publickey
	candidateWallet, err := wallet.Base58CheckDeserialize(candidatePaymentAddress)
	if err != nil || candidateWallet == nil {
		return nil, errors.New("Base58CheckDeserialize candidate Payment Address failed")
	}
	pk := candidateWallet.KeySet.PaymentAddress.Pk

	committeePK, err := incognitokey.NewCommitteeKeyFromSeed(privateSeedBytes, pk)
	if err != nil {
		return nil, err
	}

	committeePKBytes, err := committeePK.Bytes()
	if err != nil {
		return nil, err
	}

	stakingMetadata, err := metadata.NewStopAutoStakingMetadata(int(stopAutoStakingType), base58.Base58Check{}.Encode(committeePKBytes, common.ZeroByte))
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

func CreateAndSendStopAutoStakingTransaction(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	var err error
	data, err := handleCreateRawStopAutoStakingTransaction(rpcClient, params)
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
