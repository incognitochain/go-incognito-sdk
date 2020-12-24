package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestCreateAndSendTxWithIssuingETHReq(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	meta := map[string]interface{}{
		"IncTokenID": "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
		"BlockHash":  "1234",
		"ProofStrs":  []string{"123"},
		"TxIndex":    uint(1),
	}

	transParams := map[string]uint64{"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA": 0}

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		transParams,
		5, //
		-1,
		meta,
	}

	data, err := CreateAndSendTxWithIssuingETHReq(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create raw data: %v\n", err)
		return
	}

	fmt.Printf("Create raw data successfully - Data %v !!!", data)
}

func TestCreateAndSendTxWithIssuingReq(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	meta := map[string]interface{}{
		"ReceiveAddress":  "12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA",
		"DepositedAmount": uint64(1000000000),
		"TokenID":         "a0a22d131bbfdc892938542f0dbe1a7f2f48e16bc46bf1c5404319335dc1f0df", // or ptoken.TokenID
		"TokenName":       "ETH",
		"TokenSymbol":     "pETH",
	}

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		nil,
		5, //
		-1,
		meta,
	}

	data, err := CreateAndSendIssuingRequest(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create raw data: %v\n", err)
		return
	}

	fmt.Printf("Create raw data successfully - Data %v !!!", data)
}
