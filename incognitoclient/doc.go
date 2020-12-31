/*
Package go-incognito-sdk is a tutorial that integrate with Incognito Chain.

Notice that this doc is written in godoc itself as package documentation.
The defined types are just for making the table of contents at the
head of the page; they have no meanings as types.

If you have any suggestion or comment, please feel free to open an issue on
this tutorial's GitHub page!

By Incognito.

Installation

To download the specific tagged release, run:

	go get github.com/incognitochain/go-incognito-sdk@V0.0.1

It requires Go 1.13 or later due to usage of Go Modules.

Usage

Initialize object Blockchain and play rock

Example:

	package main

	import (
		"fmt"
		"github.com/incognitochain/go-incognito-sdk/incognitoclient"
		"github.com/incognitochain/go-incognito-sdk/incognitoclient/entity"
		"net/http"
	)

	func main() {
		client := &http.Client{}
		publicIncognito := NewPublicIncognito(client, "https://testnet.incognito.org/fullnode")

		blockInfo := NewBlockInfo(publicIncognito)
		wallet := NewWallet(publicIncognito, blockInfo)

		//create new a wallet
        paymentAddress, pubkey, readonlyKey, privateKey, validatorKey, shardId , _ := wallet.CreateWallet()
        fmt.Println("payment adresss", paymentAddress)
        fmt.Println("public key", pubkey)
        fmt.Println("readonly key", readonlyKey)
        fmt.Println("private key", privateKey)
        fmt.Println("validator key", validatorKey)
        fmt.Println("shard id", shardId)

        //send prv
		listPaymentAddresses := entity.WalletSend{
			Type: 0,
			PaymentAddresses: map[string]uint64{
				"12Rsf3wFnThr3T8dMafmaw4b3CzUatNao61dkj8KyoHfH5VWr4ravL32sunA2z9UhbNnyijzWFaVDvacJPSRFAq66HU7YBWjwfWR7Ff": 500000000000,
			},
		}

		tx, err := wallet.CreateAndSendConstantTransaction(
			"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
			listPaymentAddresses,
		)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(tx)
	}
*/
package incognitoclient

// Metadata input for BurnCentralizedToken
var BurnCentralizedTokenMetadata = map[string]interface{}{
	"TokenID":     "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
	"Privacy":     true,
	"TokenTxType": 1,
	"TokenName":   "ETH",
	"TokenSymbol": "pETH",
	"TokenAmount": uint64(1000000000),
	"TokenReceivers": map[string]uint64{
		"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": uint64(1000000000),
	},
	"TokenFee": uint64(0),
}

// Metadata input for MintDecentralizedToken
var MintDecentralizedTokenMetadata = map[string]interface{}{
	"BlockHash":  "block hash of deposit coin",
	"IncTokenID": "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
	"ProofStrs":  "proof of deposit coin",
	"TxIndex":    "tx index of of deposit coin",
}

const (
	ChargeFeeViaAutoCalculatePrvFee = -1
	ChargeFeeViaPrv                 = 0
	ChargeFeeViaToken               = 1
)

//Status of node validator
const (
	RoleNodeStatusNotStake  = -1
	RoleNodeStatusCandidate = 0
	RoleNodeStatusCommittee = 1
)

//Status of mint centralized token
const (
	MintCentralizedStatusNotFound = 0
	MintCentralizedStatusPending  = 1
	MintCentralizedStatusSuccess  = 2
	MintCentralizedStatusReject   = 3
)
