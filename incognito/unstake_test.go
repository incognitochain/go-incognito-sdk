package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestCreateAndSendStopAutoStakingTransaction(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		map[string]uint64{"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": 0},
		10,
		0,
		map[string]interface{}{
			"StopAutoStakingType": 127,
			"CandidatePaymentAddress":
			"12S3Cm7ZyzzheDNLrke2V4fpPuSvRZnMpWA1X99aXhKXa3VLNqAiQkNBWGTs6549JUrCSA9LjzsMmueqAWfcYQWqsC9WLoVgJ8fhEsL",
			"PrivateSeed": "12NWC4aCvgXZWT1SZZEBsZFrgovQhR9GjQ8Q1JhpiT3zsK47Y2t",
		},
	}

	data, err := CreateAndSendStopAutoStakingTransaction(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create and send normal tx %v\n", err)
		return
	}

	fmt.Printf("Send tx successfully - Data %v !!!", data)
}
