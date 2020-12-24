package repository

import (
	"github.com/incognitochain/go-incognito-sdk/incognito"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/wallet"
)

type IntegrationInterface interface {
	CreateAndSendConstantTransaction(param interface{}) (interface{}, error)
	SendPrivacyCustomTokenTransaction(params interface{}) (interface{}, error)
	CreateAndSendIssuingRequest(params interface{}) (interface{}, error)
	CreateAndSendTxWithIssuingEth(params interface{}) (interface{}, error)
	CreateAndSendBurningForDepositToSCRequest(params interface{}) (interface{}, error)
	CreateAndSendContractingRequest(params interface{}) (interface{}, error)
	CreateAndSendStakingTx(params interface{}) (interface{}, error)
	CreateAndSendStopAutoStakingTransaction(params interface{}) (interface{}, error)
	CreateAndSendWithDrawTransaction(params interface{}) (interface{}, error)
	DefragmentationPrv(params interface{}) (interface{}, error)
	DefragmentationPToken(params interface{}) (interface{}, error)
	GetBalance(privateKey string, tokenId string) (uint64, error)
	CreateWalletAddress() (*wallet.KeySerializedData, error)
	CreateNewWalletByShardId(shardId int) (*wallet.KeySerializedData, error)
	GetUTXO(privateKey string, tokenId string) ([]*privacy.InputCoin, error)
}

type IncChainIntegration struct {
	RpcClient *rpcclient.HttpClient
}

//prv, normal tx and privacy tx
func (i IncChainIntegration) CreateAndSendConstantTransaction(param interface{}) (interface{}, error) {
	return incognito.CreateAndSendTx(i.RpcClient, param)
}

//pETH, pBTC
func (i IncChainIntegration) SendPrivacyCustomTokenTransaction(params interface{}) (interface{}, error) {
	return incognito.CreateAndSendPrivacyCustomTokenTransaction(i.RpcClient, params)
}

func (i IncChainIntegration) CreateAndSendIssuingRequest(params interface{}) (interface{}, error) {
	return incognito.CreateAndSendIssuingRequest(i.RpcClient, params)
}

func (i IncChainIntegration) CreateAndSendTxWithIssuingEth(params interface{}) (interface{}, error) {
	return incognito.CreateAndSendTxWithIssuingETHReq(i.RpcClient, params)
}

func (i IncChainIntegration) CreateAndSendBurningForDepositToSCRequest(params interface{}) (interface{}, error) {
	return incognito.CreateAndSendBurningForDepositToSCRequest(i.RpcClient, params)
}

func (i IncChainIntegration) CreateAndSendContractingRequest(params interface{}) (interface{}, error) {
	return incognito.CreateAndSendContractingRequest(i.RpcClient, params)
}

func (i IncChainIntegration) GetBalance(privateKey string, tokenId string) (uint64, error) {
	return incognito.GetBalance(i.RpcClient, privateKey, tokenId)
}

func (i IncChainIntegration) CreateWalletAddress() (*wallet.KeySerializedData, error) {
	return incognito.CreateNewWallet()
}

func (i IncChainIntegration) CreateAndSendStakingTx(params interface{}) (interface{}, error) {
	return incognito.CreateAndSendStakingTx(i.RpcClient, params)
}

func (i IncChainIntegration) CreateAndSendStopAutoStakingTransaction(params interface{}) (interface{}, error) {
	return incognito.CreateAndSendStopAutoStakingTransaction(i.RpcClient, params)
}

func (i IncChainIntegration) CreateAndSendWithDrawTransaction(params interface{}) (interface{}, error) {
	return incognito.CreateAndSendWithDrawTransaction(i.RpcClient, params)
}

func (i IncChainIntegration) CreateNewWalletByShardId(shardId int) (*wallet.KeySerializedData, error) {
	return incognito.CreateNewWalletByShardId(shardId)
}

func (i IncChainIntegration) DefragmentationPrv(param interface{}) (interface{}, error) {
	return incognito.DeFragmentAccount(i.RpcClient, param)
}

func (i IncChainIntegration) DefragmentationPToken(param interface{}) (interface{}, error) {
	return incognito.DeFragmentPTokenAccount(i.RpcClient, param)
}

func (i IncChainIntegration) GetUTXO(privateKey string, tokenId string) ([]*privacy.InputCoin, error) {
	return incognito.GetUTXO(i.RpcClient, privateKey, tokenId)
}

func NewIncChainIntegration(rpcClient *rpcclient.HttpClient) *IncChainIntegration {
	return &IncChainIntegration{
		RpcClient: rpcClient,
	}
}
