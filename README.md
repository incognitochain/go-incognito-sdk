### Incognito SDK

This is SDK support integration with a chain. That mean you creating wallet, sending token, staking node, trading,
etc... with Incognito blockchain.

### Use

Incognito SDK use as library

Installation

Using go module

```
go get github.com/incognitochain/go-incognito-sdk@new-tag
```

Initialization

Init new PublicIncognito, setup endpoint url environment

Testnet: https://testnet.incognito.org/fullnode

```
client := &http.Client{}
publicIncognito := NewPublicIncognito(client, "https://testnet.incognito.org/fullnode")
```

To create new wallet init new Wallet Object

```
blockInfo := NewBlockInfo(publicIncognito)
wallet := NewWallet(publicIncognito, blockInfo)

//create new wallet
wallet.CreateAndSendConstantTransaction(....)
```

To staking a node

```
stake = NewStake(publicIncognito)
stake.ListUnstake()

```

All together

```
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
		
		//send Prv token
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
```

### How to works

Incognito SDK wrap all RPC of blockchain, build raw data at local before call rpc because can't send private key to
chain

The steps:

1. Build raw data
2. Call rpc to chain
3. Get result

UML Diagram

![Screenshot](UMLDiagram.png)

### Godoc

- With Docker

```
    docker build -t godoc-sdk .
    docker run -p 6060:6060  --name godoc-sdk  godoc-sdk
```

- Without Docker

```
godoc -http=:6060
```

Access endpoint

http://localhost:6060/pkg/github.com/incognitochain/go-incognito-sdk/incognitoclient
