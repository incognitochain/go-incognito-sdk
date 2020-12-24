package incognitoclient

import (
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/constant"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/entity"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/repository"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/service"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"math/big"
	"net/http"
)

const PRVToken = "0000000000000000000000000000000000000000000000000000000000000004"

// Public  of incognito
type PublicIncognito struct {
	incClient      *service.IncogClient
	incIntegration *repository.IncChainIntegration
}

/*

Creat new PublicIncognito

Input:
	- c: init http client module (*http.Client)
	- endpointUri: endpoint chain url (string)

Output:
	- result: return an object (*PublicIncognito)

Example:
	client := &http.Client{}
	publicIncognito := NewPublicIncognito(client, "https://testnet.incognito.org/fullnode")

*/
func NewPublicIncognito(c *http.Client, endpointUri string) *PublicIncognito {
	inc := &service.IncogClient{
		Client:        c,
		ChainEndpoint: endpointUri,
	}

	rpcClient := rpcclient.NewHttpClient(endpointUri, "https", endpointUri, 0)
	incIntegration := repository.NewIncChainIntegration(rpcClient)

	return &PublicIncognito{incClient: inc, incIntegration: incIntegration}
}

/*
GetPRVToken return token id of PRV
*/
func (i *PublicIncognito) GetPRVToken() string {
	return PRVToken
}

type BlockInfo struct {
	public *PublicIncognito
	block  *repository.Block
}

/*

	Creat new BlockInfo

	Example:
		client := &http.Client{}
		publicIncognito := NewPublicIncognito(client, "https://testnet.incognito.org/fullnode")
		blockInfo := BlockInfo(publicIncognito)

*/
func NewBlockInfo(public *PublicIncognito) *BlockInfo {
	block := repository.NewBlock(public.incClient)
	return &BlockInfo{public: public, block: block}
}

/*
GetBlockInfo return info of that block
*/
func (b *BlockInfo) GetBlockInfo(blockHeight int32, shardID int) (*entity.GetBlockInfo, error) {
	return b.block.GetBlockInfo(blockHeight, shardID)
}

/*
GetChainInfo return info of Incognito Chain
*/
func (b *BlockInfo) GetChainInfo() (*entity.GetBlockChainInfoResult, error) {
	return b.block.GetBlockChainInfo()
}

/*
GetBestBlockHeight return block height current of any shard id
*/
func (b *BlockInfo) GetBestBlockHeight(shardID int) (uint64, error) {
	return b.block.GetBestBlockHeight(shardID)
}

/*
GetBeaconHeight return beacon height current
*/
func (b *BlockInfo) GetBeaconHeight() (int32, error) {
	return b.block.GetBeaconHeight()
}

/*
GetBeaconBestStateDetail return beacon stage detail current
*/
func (b *BlockInfo) GetBeaconBestStateDetail() (res *entity.BeaconBestStateResp, err error) {
	return b.block.GetBeaconBestStateDetail()
}

/*
GetBurningAddress return burn address of chain, burn address has burned token
*/
func (b *BlockInfo) GetBurningAddress() (string, error) {
	return b.block.GetBurningAddress()
}

type PDex struct {
	public *PublicIncognito
	pdex   *repository.Pdex
}

/*

	Creat new PDex

	Example:
		client := &http.Client{}
		publicIncognito := NewPublicIncognito(client, "https://testnet.incognito.org/fullnode")
		blockInfo := NewBlockInfo(publicIncognito)

		wallet := NewWallet(publicIncognito, blockInfo)
		pdex := NewPDex(publicIncognito, blockInfo)

*/
func NewPDex(public *PublicIncognito, block *BlockInfo) *PDex {
	pdex := repository.NewPdex(public.incClient, public.GetPRVToken(), block.block)
	return &PDex{public: public, pdex: pdex}
}

/*
GetPDexState return all pair pdex by beacon height
*/
func (b *PDex) GetPDexState(beaconHeight int32) (map[string]interface{}, error) {
	return b.pdex.GetPDexState(beaconHeight)
}

/*
TradePDex will trade pair token, sell this token and buy that token

Input:
	- privateKey: private key of trader  (string)
	- buyTokenId: token id of buy coin (string)
	- tradingFee:  amount trading fee to pay for trade if have (uint64)
	- sellTokenId: token id of sell coin (string)
	- sellTokenAmount: amount to sell (uint64)
	- minimumAmount: minimum amount can receive (uint64)
	- traderAddress: address of trader (string)
	- networkFeeTokenID: amount network fee to pay for trade  (string)
	- networkFee: amount of network fee (uint64)

Output:
	- result: tx hash (string)
	- err: err (error)

Example:

	//trade 1 PRV -> pDai
	tx, err := pdex.TradePDex(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0",
		uint64(100),
		t.client.GetPRVToken(),
		uint64(1000000000),
		uint64(0),
		"12RwamF5njyL5cqpiMZ3SrqGHMqDaEDLyQexeaHYjYn2LDMzKZzgPZHnbQ75iLBKxm4md4kiyLxrPrFRNRNNktmAMjmfD4ktmcptgiX",
		t.client.GetPRVToken(),
		uint64(100))

*/
func (b *PDex) TradePDex(privateKey string, buyTokenId string, tradingFee uint64, sellTokenId string, sellTokenAmount uint64, minimumAmount uint64, traderAddress string, networkFeeTokenID string, networkFee uint64) (string, error) {
	return b.pdex.TradePDex(privateKey, buyTokenId, tradingFee, sellTokenId, sellTokenAmount, minimumAmount, traderAddress, networkFeeTokenID, networkFee)
}

/*
GetPDexTradeStatus return status of trade tx
*/
func (b *PDex) GetPDexTradeStatus(txId string) (constant.PDexTradeStatus, error) {
	return b.pdex.GetPDexTradeStatus(txId)
}

type Stake struct {
	public *PublicIncognito
	stake  *repository.Stake
}

/*

	Creat new Stake

	Example:
		client := &http.Client{}
		publicIncognito := NewPublicIncognito(client, "https://testnet.incognito.org/fullnode")
		stake := NewStake(publicIncognito)

*/
func NewStake(public *PublicIncognito) *Stake {
	stake := repository.NewStake(public.incClient, public.incIntegration)
	return &Stake{public: public, stake: stake}
}

/*
ListUnstake return all node unstake
*/
func (b *Stake) ListUnstake() ([]entity.Unstake, error) {
	return b.stake.ListUnstake()
}

/*
Staking is action to stake a node validator

Input:
	- receiveRewardAddress: payment address of staker, address receive reward  (string)
	- privateKey: private key of staker (string)
	- userPaymentAddress:  payment address of node (string)
	- userValidatorKey: validator key of node (string)
	- burnTokenAddress: burn address (string)

Output:
	- result: tx hash (string)
	- err: err (error)

*/
func (b *Stake) Staking(receiveRewardAddress, privateKey, userPaymentAddress, userValidatorKey, burnTokenAddress string) (string, error) {
	return b.stake.Staking(receiveRewardAddress, privateKey, userPaymentAddress, userValidatorKey, burnTokenAddress)
}

/*
Unstaking is action to unstake node validator

Input:
	 - privateKey: private key of staker (string)
	 - userPaymentAddress:  payment address of node (string)
	 - userValidatorKey: validator key of node (string)
	 - burnTokenAddress: burn address (string)

Output:
	- result: tx hash (string)
	- err: err (error)

*/
func (b *Stake) Unstaking(privateKey, userPaymentAddress, userValidatorKey, burnTokenAddress string) (string, error) {
	return b.stake.Unstaking(privateKey, userPaymentAddress, userValidatorKey, burnTokenAddress)
}

/*
WithDrawReward is action to withdraw all reward earn of node validator

Input:
	 - privateKey: private key of staker (string)
	 - paymentAddress:  payment address of staker (string)
	 - tokenId: token id (string)

Output:
	- result: tx hash (string)
	- err: err (error)

*/
func (b *Stake) WithDrawReward(privateKey, paymentAddress, tokenId string) (string, error) {
	return b.stake.WithDrawReward(privateKey, paymentAddress, tokenId)
}

/*
GetRewardAmount return list amount reward each token

Input:
    - paymentAddress: payment address of staker (string)

Output:
	- result: List reward item ([]entity.RewardItems)
	- err: err

*/
func (b *Stake) GetRewardAmount(paymentAddress string) ([]entity.RewardItems, error) {
	return b.stake.GetRewardAmount(paymentAddress)
}

/*
ListRewardAmounts return list reward Prv amount all node validator

Output:
	- result: list reward item ([]entity.RewardAmount)
	- error: error (error)

*/

func (b *Stake) ListRewardAmounts() ([]entity.RewardAmount, error) {
	return b.stake.ListRewardAmounts()
}

/*
GetStatusNodeValidator return status of node validator

Input:
    - validatorKey: validator key of node (string)

Output:
	- result: status of node (RoleNodeStatusNotStake | RoleNodeStatusCandidate | RoleNodeStatusCommittee) (float)
	- err: err

*/
func (b *Stake) GetStatusNodeValidator(validatorKey string) (float64, error) {
	return b.stake.GetNodeAvailable(validatorKey)
}

/*
GetTotalStaker return total staker

Output:
	- result: Total staker (float)
	- err: err

*/
func (b *Stake) GetTotalStaker() (float64, error) {
	return b.stake.GetTotalStaker()
}

type Wallet struct {
	public *PublicIncognito
	wallet *repository.Wallet
}

/*

	Creat new Wallet

	Example:
		client := &http.Client{}
		publicIncognito := NewPublicIncognito(client, "https://testnet.incognito.org/fullnode")
		block = NewBlockInfo(publicIncognito)
		wallet := NewWallet(publicIncognito, block)

*/
func NewWallet(public *PublicIncognito, block *BlockInfo) *Wallet {
	wallet := repository.NewWallet(public.incClient, public.GetPRVToken(), block.block, public.incIntegration)
	return &Wallet{public: public, wallet: wallet}
}

/*
CreateWallet return all info of a wallet

Output:
	- paymentAddress: payment address  (string)
	- pubkey: public key (string)
	- readonlyKey: readonly key (string)
	- privateKey: private key (string)
	- validatorKey: validator key (string)
	- shardId: shard id (int)
	- err: err (error)

*/
func (b *Wallet) CreateWallet() (paymentAddress, pubkey, readonlyKey, privateKey, validatorKey string, shardId int, err error) {
	return b.wallet.CreateNodeWalletAddress(-1)
}

/*
CreateWallet return all info of a wallet

Input:
	- byShardId: shard id (int)

Output:
	- paymentAddress: payment address  (string)
	- pubkey: public key (string)
	- readonlyKey: readonly key (string)
	- privateKey: private key (string)
	- validatorKey: validator key (string)
	- shardId: shard id (int)
	- err: err (error)

*/
func (b *Wallet) CreateWalletByShardId(byShardId int) (paymentAddress, pubkey, readonlyKey, privateKey, validatorKey string, shardId int, err error) {
	return b.wallet.CreateNodeWalletAddress(byShardId)
}

/*
ListPrivacyCustomToken return all token info of chain as token id, logo, name, symbol

Output:
	- result: list reward item ([]entity.PCustomToken)
	- error: error (error)

*/
func (b *Wallet) ListPrivacyCustomToken() ([]entity.PCustomToken, error) {
	return b.wallet.ListPrivacyCustomToken()
}

/*
GetTransactionDetailByTxHash return info detail of a tx

Input:
	- txHash: transaction (string)

Output:
	- result: transaction detail (entity.TransactionDetail)
	- error: error (error)

*/
func (b *Wallet) GetTransactionDetailByTxHash(txHash string) (*entity.TransactionDetail, error) {
	return b.wallet.GetTxByHash(txHash)
}

/*
GetMintStatusCentralized return status tx of mint action

Input:
	- txHash: tx hash (string)

Output:
	- result: status tx (MintCentralizedStatusNotFound | MintCentralizedStatusPending | MintCentralizedStatusSuccess | MintCentralizedStatusReject) (int)
	- error: error (error)

*/
func (b *Wallet) GetMintStatusCentralized(txHash string) (int, error) {
	return b.wallet.GetBridgeReqWithStatus(txHash)
}

/*
MintCentralizedToken is action to mint a coin with anything amount in incognito chain. Your individual chain need setup master private key, it's is only key of chain

Input:
	- privateKey: master private key of chain (string)
	- receiveAddress: receiver address (string)
	- depositedAmount: deposit amount (*big.Int)
	- tokenId: token id (string)
	- tokenName: token name (string)

Output:
	- result: transaction id (entity.TransactionDetail)
	- error: error (error)

Example:

	privateKey := "112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5"
	receiveAddress := "12RwamF5njyL5cqpiMZ3SrqGHMqDaEDLyQexeaHYjYn2LDMzKZzgPZHnbQ75iLBKxm4md4kiyLxrPrFRNRNNktmAMjmfD4ktmcptgiX"
	depositedAmount := big.NewInt(2000000000)
	tokenId := "4584d5e9b2fc0337dfb17f4b5bb025e5b82c38cfa4f54e8a3d4fcdd03954ff82"
	tokenName := "BTC"

	//mint 2 BTC
	tx, err := t.wallet.MintCentralizedToken(privateKey, receiveAddress, depositedAmount, tokenId, tokenName)

*/

func (b *Wallet) MintCentralizedToken(privateKey, receiveAddress string, depositedAmount *big.Int, tokenId string, tokenName string) (string, error) {
	return b.wallet.CreateAndSendIssuingRequest(privateKey, receiveAddress, depositedAmount, tokenId, tokenName)
}

/*
BurnCentralizedToken is action to burn token

Input:
	- privateKey: private key hold token is burned (string)
	- autoChargePRVFee: method to calculate fee (ChargeFeeViaAutoCalculatePrvFee | ChargeFeeViaPrv | ChargeFeeViaToken) (int)
	- metadata: info to burn (BurnCentralizedTokenMetadata) (map[string]interface{})
		TokenFee: it is must more than 0 when charge fee via token else equal 0 with charge fee via Prv

Output:
	- result: tx hash (string)
	- error: error (error)

Example:
	masterPrivKey := "112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5"
	autoChargePRVFee := -1
	burnAddress, _ := block.GetBurningAddress()

	txID, err := wallet.BurnCentralizedToken(
		masterPrivKey,
		autoChargePRVFee,
		map[string]interface{}{
			"TokenID":     "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854,
			"Privacy":     true,
			"TokenTxType": 1,
			"TokenName":   "ETH",
			"TokenSymbol": "pETH,
			"TokenAmount": 1000000000,
			"TokenReceivers": map[string]uint64{
				burnAddress: 1000000000,
			},
			"TokenFee": 0,
	})

*/
func (b *Wallet) BurnCentralizedToken(privateKey string, autoChargePRVFee int, metadata map[string]interface{}) (string, error) {
	return b.wallet.CreateAndSendContractingRequestForPrivacyToken(privateKey, autoChargePRVFee, metadata)
}

/*
MintDecentralizedToken is action to mint decentralized token as ETH, Erc20. To done this action, first you must deposit coin to Ethereum chain after that you need get proof deposit which to mint token

Input:
	- privateKey: private key (string)
	- burnerAddress: burn address (string)
	- metadata: deposit amount (MintDecentralizedTokenMetadata)
		BlockHash, ProofStrs, TxIndex: this params got from Merkle proof in Ethereum
		Reference: https://blog.ethereum.org/2015/11/15/merkling-in-ethereum/

Output:
	- txHash: tx hash (string)
	- res: response data ([]byte)
	- error: error (error)

Example:
	privateKey := "112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5"
	burningAddress, _ := block.GetBurningAddress()

	txID, body, err := wallet.MintDecentralizedToken(
		privateKey,
		burningAddress,
		map[string]interface{}{
			"BlockHash":  "test",
			"IncTokenID": "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
			"ProofStrs":  "test",
			"TxIndex":    "test",
		})

*/
func (b *Wallet) MintDecentralizedToken(privateKey, burnerAddress string, metadata map[string]interface{}) (txHash string, res []byte, err error) {
	return b.wallet.CreateAndSendTxWithIssuingEth(privateKey, burnerAddress, metadata)
}

/*
BurnDecentralizedToken is action to burn eth, erc20 token. Advantage, after burn token you can get proof burn to deposit amount to smart contract

Input:
	- incPrivateKey: private key hold token is burned (string)
	- amount: burn amount (bigInt)
	- receiverAddress: receiver coin, format is erc20 address (string)
	- tokenId: token (string)

Output:
	- result: response burn token  (*entity.BurningForDepositToSCRes)
	- error: error (error)

Example:

	addStr := "0x15B9419e738393Dbc8448272b18CdE970a07864D"
	result, err := wallet.BurnDecentralizedToken(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		big.NewInt(1000000),
		addStr[2:],
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
	)
*/

func (b *Wallet) BurnDecentralizedToken(incPrivateKey string, amount *big.Int, receiverAddress string, tokenId string) (*entity.BurningForDepositToSCRes, error) {
	return b.wallet.CreateAndSendBurningForDepositToSCRequest(incPrivateKey, amount, receiverAddress, tokenId)
}

/*
GenerateTokenID return new token id, input token info defined by yourself. Note you shouldn't generate new token which is exist token

Input:
	- symbol: symbol of new token (string)
	- pSymbol: pSymbol of new token (string)

Output:
	- result: new token id  (string)
	- error: error (error)

Example:

	tokenId, err := wallet.GenerateTokenID("ABC", "pABC")

*/
func (b *Wallet) GenerateTokenID(symbol, pSymbol string) (string, error) {
	return b.wallet.GenerateTokenID(symbol, pSymbol)
}

/*
GetPublicKeyFromPaymentAddress return public key of payment address

Input:
	- paymentAddress: payment address (string)

Output:
	- result: public key (string)
	- error: error (error)
*/
func (b *Wallet) GetPublicKeyFromPaymentAddress(paymentAddress string) (string, error) {
	return b.wallet.GetPublickeyFromPaymentAddress(paymentAddress)
}

/*
GetTransactionByReceiversAddress return list transaction detail of payment address

Input:
	- paymentAddress: payment address (string)
	- readonlyKey: read only key (string)

Output:
	- result: transaction info (*entity.ReceivedTransactions)
	- error: error (error)

*/
func (b *Wallet) GetTransactionByReceiversAddress(paymentAddress, readonlyKey string) (*entity.ReceivedTransactions, error) {
	return b.wallet.GetTransactionByReceivers(paymentAddress, readonlyKey)
}

/*
GetBalance return current balance of wallet

Input:
	- privateKey: private key (string)
	- tokenId: token id (string)

Output:
	- result: amount (uint64)
	- Error: error (error)

Example:
	//get balance prv
	tx1, err1 := wallet.GetBalance("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", "0000000000000000000000000000000000000000000000000000000000000004")
	//get balance Eth
	tx2, err2 := wallet.GetBalance("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854")
*/

func (b *Wallet) GetBalance(privateKey string, tokenId string) (uint64, error) {
	return b.wallet.GetBalance(privateKey, tokenId)
}

/*
GetTransactionAmount return amount of transaction

Input:
	- txId: transaction (string)
	- walletAddress: wallet address (string)
	- readOnlyKey: read only key (string)

Output:
	- result: amount (uint64)
	- Error: error (error)
*/

func (b *Wallet) GetTransactionAmount(txId string, walletAddress string, readOnlyKey string) (uint64, error) {
	return b.wallet.GetTransactionAmount(txId, walletAddress, readOnlyKey)
}

/*
SendToken is action to send token or prv to receiver address

Input:
	- privateKey: incognito private key (string)
	- receiverAddress: address of receiver (string)
	- tokenId: token (string)
	- amount: amount to send (uint64)
	- fee: amount fee (uint64)
	- feeTokenId: token id of fee (string)

if you send Prv: you ignore fee and feeTokenId.

if you send Token:
	+ charge fee via Prv: fill prv token id and fee
	+ charge fee via Token: fill token id and fee

Output:
	- result: tx hash (string)
	- error: error (error)

Example:
	//send prv with fee is prv
	tx, err := wallet.SendToken(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"12Rsf3wFnThr3T8dMafmaw4b3CzUatNao61dkj8KyoHfH5VWr4ravL32sunA2z9UhbNnyijzWFaVDvacJPSRFAq66HU7YBWjwfWR7Ff",
		PRVToken,
		500000000000,
		0,
		",
	)

	//send ptoken with fee is prv
    tx1, err1 := wallet.SendToken(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"12RqaTLErSnN88pGgXaKmw1PSQEaG86FA4uJsm32RZetAy7e5yEncqjTC6QJcMRjMfTSc48tcWRTyy8FoB9VkCHu56Vd9b86gd8Pq8k",
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
		9800000,
		5,
		"0000000000000000000000000000000000000000000000000000000000000004")

	//send ptoken with fee is ptoken
	tx2, err2 := t.wallet.SendToken(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"12RqaTLErSnN88pGgXaKmw1PSQEaG86FA4uJsm32RZetAy7e5yEncqjTC6QJcMRjMfTSc48tcWRTyy8FoB9VkCHu56Vd9b86gd8Pq8k",
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
		9800000,
		100,
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854")
*/
func (b *Wallet) SendToken(privateKey string, receiverAddress string, tokenId string, amount uint64, fee uint64, feeTokenId string) (string, error) {
	return b.wallet.SendToken(privateKey, receiverAddress, tokenId, amount, fee, feeTokenId)
}

/*
Defragmentation is action to merge utxo of wallet

Input:
	- privateKey: incognito private key (string)
	- maxValue: Max value (int64)
		maxValue: it is useful with Prv token
	- tokenId: token (string)

Output:
	- result: tx hash (string)
	- error: error (error)

Example:
	//merge all utxo of Prv token
	tx, err := wallet.Defragmentation("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", int64(500000*1e9), PRVToken)

	//merge all utxo of token id
	tx1, err1 := t.wallet.Defragmentation("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", 0, "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854")


*/
func (b *Wallet) Defragmentation(privateKey string, maxValue int64, tokenId string) (string, error) {
	if tokenId == b.public.GetPRVToken() {
		return b.wallet.DefragmentationPrv(privateKey, maxValue)
	}

	return b.wallet.DefragmentationPToken(privateKey, tokenId)
}

/*
GetUTXO return all unspent output coin except spending of wallet

Input:
	- privateKey: incognito private key (string)
	- tokenId: token (string)

Output:
	- result: list utxo ([]*entity.Utxo)
	- error: error (error)
*/
func (b *Wallet) GetUTXO(privateKey string, tokenId string) ([]*entity.Utxo, error) {
	return b.wallet.GetUTXO(privateKey, tokenId)
}
