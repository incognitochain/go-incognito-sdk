package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestCreateAndSendWithDrawTransaction(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		nil,
		0,
		0,
		map[string]interface{}{
			"PaymentAddress": "12S3Cm7ZyzzheDNLrke2V4fpPuSvRZnMpWA1X99aXhKXa3VLNqAiQkNBWGTs6549JUrCSA9LjzsMmueqAWfcYQWqsC9WLoVgJ8fhEsL",
			"TokenID":        "375a76afdbf9aadb6e23a181e86de643cd53cf576d7b3fead0086105e550321f",
		},
	}

	data, err := CreateAndSendWithDrawTransaction(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create and send normal tx %v\n", err)
		return
	}

	fmt.Printf("Send tx successfully - Data %v !!!", data)
}
