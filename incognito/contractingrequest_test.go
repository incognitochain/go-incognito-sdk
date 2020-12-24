package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestCreateAndSendContractingRequest(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	var paymentAddresses = make(map[string]uint64)
	paymentAddresses["12S2YrSLXQVS2rckxZUQ1o72d9AN8xktAhJiCZDa7Lo3SLKCPuWHEv9uAE2WjoLVsYjXP1EqFM5xqMqzihQ33tuzpeU3joJsM3dUJpN"] = uint64(1000000000)

	metadata := map[string]interface{}{}
	metadata["Privacy"] = true
	metadata["TokenID"] = "a0a22d131bbfdc892938542f0dbe1a7f2f48e16bc46bf1c5404319335dc1f0df"
	metadata["TokenTxType"] = 1
	metadata["TokenName"] = "ZIL"
	metadata["TokenSymbol"] = "pZIL"
	metadata["TokenReceivers"] = paymentAddresses
	metadata["TokenAmount"] = uint64(1000000000)
	metadata["TokenFee"] = uint64(0)

	autoChargePRVFee := -1

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		nil,
		autoChargePRVFee, //-1
		-1,
		metadata,
		"",
		0,
	}

	data, err := CreateAndSendContractingRequest(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create and send normal tx %v\n", err)
		return
	}

	fmt.Printf("Send tx successfully - Data %v !!!", data)
}
