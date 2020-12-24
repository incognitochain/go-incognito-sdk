package constant

const (
	// blockchain network
	Getblockchaininfo     = "getblockchaininfo"
	Retrieveblockbyheight = "retrieveblockbyheight"
	GetEncryptionFlag     = "getencryptionflag"

	DumpPrivKeyMethod        = "dumpprivkey"
	GetAccountAddress        = "getaccountaddress"
	GetBeaconBestStateDetail = "getbeaconbeststatedetail"
	GetTotalStaker           = "gettotalstaker"

	GetBurningAddress = "getburningaddress"

	// wallet methods
	ListAccountsMethod                     = "listaccounts"
	GetAccount                             = "getaccount"
	EncryptData                            = "encryptdata"
	GetBalanceByPrivateKeyMethod           = "getbalancebyprivatekey"
	GetBalanceByPaymentAddress             = "getbalancebypaymentaddress"
	GetCurrentSellingGOVTokens             = "getcurrentsellinggovtokens"
	CreateAndSendTxWithBuyGOVTokensRequest = "createandsendtxwithbuygovtokensrequest"

	GetPublickeyFromPaymentAddress = "getpublickeyfrompaymentaddress"
	GetShardFromPaymentAddress     = "getpublickeyfrompaymentaddress"

	// tx
	CreateAndSendTransaction                   = "createandsendtransaction"
	CreateAndSendCustomTokenTransaction        = "createandsendcustomtokentransaction"
	CreateAndSendPrivacyCustomTokenTransaction = "createandsendprivacycustomtokentransaction"
	GetTransactionByHash                       = "gettransactionbyhash"
	Createandsendloanresponse                  = "createandsendloanresponse"
	CreateAndSendLoanPayment                   = "createandsendloanpayment"
	CreateAndSendLoanWithdraw                  = "createandsendloanwithdraw"
	CheckBorrowApproved                        = "getloanresponseapproved"
	CheckBorrowRejected                        = "getloanresponserejected"
	GetBorrowPaymentInfo                       = "getloanpaymentinfo"
	DecryptOutputCoinByKeyOfTransaction        = "decryptoutputcoinbykeyoftransaction"
	GetTransactionByReceiver                   = "gettransactionbyreceiver"

	GetEstimateFee = "estimatefee"

	// custom token
	GetListCustomTokenBalance               = "getlistcustomtokenbalance"
	GetListprivacyPrivacyCustomTokenBalance = "getlistprivacycustomtokenbalance"
	GetAmountVoteToken                      = "getamountvotetoken"
	ListPrivacyCustomToken                  = "listprivacycustomtoken"

	CreateAndSendIssuingRequest     = "createandsendissuingrequest"
	GetIssuingStatus                = "getissuingstatus"
	CreateAndSendContractingRequest = "createandsendcontractingrequest"
	GetContractingStatus            = "getcontractingstatus"

	// deposit erc20:
	CreateAndSendTxWithIssuingEthReq = "createandsendtxwithissuingethreq"

	GetBridgeReqWithStatus = "getbridgereqwithstatus"
	GetBlockCount          = "getblockcount"
	GetBeaconSwapProof     = "getbeaconswapproof"
	GetBridgeSwapProof     = "getbridgeswapproof"

	// stake:
	CreateAndSendStakingTransaction   = "createandsendstakingtransaction"
	CreateAndSendUnStakingTransaction = "createandsendstopautostakingtransaction"

	GetPDETradeStatus = "getpdetradestatus"

	WithDrawReward     = "withdrawreward"
	RewardAmount       = "getrewardamount"
	RoleByValidatorKey = "getrolebyvalidatorkey"

	//get ptoken id:
	GenerateTokenID  = "generatetokenid"
	ListRewardAmount = "listrewardamount"

	GetPdeState                               = "getpdestate"
	CreateAndSendBurningForDepositToSCRequest = "createandsendburningfordeposittoscrequest"
	CreateAndSendTxWithPRVTradeReq            = "createandsendtxwithprvtradereq"
	CreateAndSendTxWithPTokenTradeReq         = "createandsendtxwithptokentradereq"

	CreateAndSendTxWithPRVCrosspollTradeReq   = "createandsendtxwithprvcrosspooltradereq"
	CreateAndSendTxWithPTokenCrosspolTradeReq = "createandsendtxwithptokencrosspooltradereq"
	SendRawTransaction                        = "sendtransaction"
	SendRawPrivacyCustomTokenTransaction      = "sendrawprivacycustomtokentransaction"
)
