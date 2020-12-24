package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestSendpTokenPrivacyTx(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	var paymentAddresses = make(map[string]uint64)
	paymentAddresses["12RqaTLErSnN88pGgXaKmw1PSQEaG86FA4uJsm32RZetAy7e5yEncqjTC6QJcMRjMfTSc48tcWRTyy8FoB9VkCHu56Vd9b86gd8Pq8k"] = uint64(9800000)

	tokenData := map[string]interface{}{}
	tokenData["Privacy"] = true
	tokenData["TokenID"] = "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
	tokenData["TokenTxType"] = 1
	tokenData["TokenName"] = ""
	tokenData["TokenSymbol"] = ""
	tokenData["TokenReceivers"] = paymentAddresses
	tokenData["TokenAmount"] = uint64(0)
	tokenData["TokenFee"] = uint64(5)
	object := map[string]uint64{}
	nativeFee := -1

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		object,
		nativeFee,
		1,
		tokenData,
		"",
		1,
	}

	data, err := CreateAndSendPrivacyCustomTokenTransaction(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create and send normal tx %v\n", err)
		return
	}

	fmt.Printf("Send tx successfully - Data %v !!!", data)
}

