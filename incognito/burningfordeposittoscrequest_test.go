package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestCreateAndSendBurningRequest(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	var paymentAddresses = make(map[string]uint64)
	paymentAddresses["12S2YrSLXQVS2rckxZUQ1o72d9AN8xktAhJiCZDa7Lo3SLKCPuWHEv9uAE2WjoLVsYjXP1EqFM5xqMqzihQ33tuzpeU3joJsM3dUJpN"] = uint64(1000000)

	incPrivateKey := "112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5"
	addStr := "0x15B9419e738393Dbc8448272b18CdE970a07864D"

	metadata := map[string]interface{}{}
	metadata["Privacy"] = true

	//token centralize
	metadata["TokenID"] = "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"

	metadata["TokenTxType"] = 1
	metadata["TokenName"] = ""
	metadata["TokenSymbol"] = ""
	metadata["TokenReceivers"] = paymentAddresses
	metadata["TokenAmount"] = uint64(1000000)
	metadata["TokenFee"] = uint64(0)
	metadata["RemoteAddress"] = addStr[2:]

	params := []interface{}{
		incPrivateKey,
		nil,
		5,
		-1,
		metadata,
		"",
		0,
	}

	data, err := CreateAndSendBurningForDepositToSCRequest(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create and send normal tx %v\n", err)
		return
	}

	fmt.Printf("Send tx successfully - Data %v !!!", data)
}
