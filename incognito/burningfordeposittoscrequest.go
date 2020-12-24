package incognito

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/bean"
	"github.com/incognitochain/go-incognito-sdk/rpcserver/rpcservice"
	"github.com/incognitochain/go-incognito-sdk/transaction"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"github.com/pkg/errors"
)

func newBurningRequestMetadata(
	senderPrivateKeyStr string,
	tokenReceivers interface{},
	tokenID string,
	tokenName string,
	remoteAddress string,
	burningMetaType int,
) (*metadata.BurningRequest, error) {
	senderKey, err := wallet.Base58CheckDeserialize(senderPrivateKeyStr)
	if err != nil {
		return nil, err
	}
	err = senderKey.KeySet.InitFromPrivateKey(&senderKey.KeySet.PrivateKey)
	if err != nil {
		return nil, err
	}
	paymentAddr := senderKey.KeySet.PaymentAddress

	_, voutsAmount, err := transaction.CreateCustomTokenPrivacyReceiverArray(tokenReceivers)
	if err != nil {
		return nil, err
	}
	tokenIDHash, err := common.Hash{}.NewHashFromStr(tokenID)
	if err != nil {
		return nil, err
	}

	meta, err := metadata.NewBurningRequest(
		paymentAddr,
		uint64(voutsAmount),
		*tokenIDHash,
		tokenName,
		remoteAddress,
		burningMetaType,
	)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func handleCreateRawTxWithBurningForDepositToSCReq(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	arrayParams := common.InterfaceSlice(params)
	if arrayParams == nil || len(arrayParams) < 5 {
		return nil, errors.New("param must be an array at least 5 elements")
	}

	if len(arrayParams) > 6 {
		privacyTemp, ok := arrayParams[6].(int)
		if !ok {
			return nil, errors.New("The privacy mode must be valid ")
		}
		hasPrivacyToken := int(privacyTemp) > 0
		if hasPrivacyToken {
			return nil, errors.New("The privacy mode must be disabled")
		}
	}

	senderPrivateKeyParam, ok := arrayParams[0].(string)
	if !ok {
		return nil, errors.New("private key is invalid")
	}

	tokenParamsRaw, ok := arrayParams[4].(map[string]interface{})
	if !ok {
		return nil, errors.New("token param is invalid")
	}

	tokenReceivers, ok := tokenParamsRaw["TokenReceivers"].(interface{})
	if !ok {
		return nil, errors.New("token receivers is invalid")
	}

	tokenID, ok := tokenParamsRaw["TokenID"].(string)
	if !ok {
		return nil, errors.New("token ID is invalid")
	}

	tokenName, ok := tokenParamsRaw["TokenName"].(string)
	if !ok {
		return nil, errors.New("token name is invalid")
	}

	remoteAddress, ok := tokenParamsRaw["RemoteAddress"].(string)
	if !ok {
		return nil, errors.New("remote address is invalid")
	}

	meta, err := newBurningRequestMetadata(
		senderPrivateKeyParam,
		tokenReceivers,
		tokenID,
		tokenName,
		remoteAddress,
		metadata.BurningForDepositToSCRequestMeta,
	)
	if err != nil {
		return nil, err
	}

	keyWallet, err := bean.GetPrivateKey(params)
	if err != nil {
		return nil, err
	}

	txService := &rpcservice.TxService{
		RpcClient: rpcClient,
		KeyWallet: keyWallet,
	}

	customTokenTx, rpcErr := txService.BuildRawPrivacyCustomTokenTransaction(params, meta)
	if rpcErr != nil {
		return nil, rpcErr
	}

	byteArrays, err2 := json.Marshal(customTokenTx)
	if err2 != nil {
		return nil, err2
	}

	result := rpcclient.CreateTransactionResult{
		TxID:            customTokenTx.Hash().String(),
		Base58CheckData: base58.Base58Check{}.Encode(byteArrays, 0x00),
	}
	return result, nil
}

//decentralized
func CreateAndSendBurningForDepositToSCRequest(rpcClient *rpcclient.HttpClient, params interface{}) (interface{}, error) {
	var err error
	data, err := handleCreateRawTxWithBurningForDepositToSCReq(rpcClient, params)
	if err != nil {
		return nil, err
	}

	tx := data.(rpcclient.CreateTransactionResult)
	base58CheckData := tx.Base58CheckData
	newParam := make([]interface{}, 0)
	newParam = append(newParam, base58CheckData)
	return newParam, nil
	//txId, err := httpServer.handleSendRawPrivacyCustomTokenTransaction(newParam, closeChan)
}
