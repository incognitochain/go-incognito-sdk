package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestCreateAndSendTxWithPTokenTradeReq(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	var paymentAddresses = make(map[string]uint64)
	paymentAddresses["12RqaTLErSnN88pGgXaKmw1PSQEaG86FA4uJsm32RZetAy7e5yEncqjTC6QJcMRjMfTSc48tcWRTyy8FoB9VkCHu56Vd9b86gd8Pq8k"] = uint64(1000000)


	metadata := map[string]interface{}{
		"Privacy":     true,
		"TokenID":     "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
		"TokenTxType": 1,
		"TokenName":   "",
		"TokenSymbol": "",
		"TokenAmount": uint64(1000000),
		"TokenReceivers": paymentAddresses,
		"TokenFee":            uint64(0),
		"TokenIDToBuyStr":     "9fca0a0947f4393994145ef50eecd2da2aa15da2483b310c2c0650301c59b17d",
		"TokenIDToSellStr":    "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
		"SellAmount":          uint64(1000000),
		"MinAcceptableAmount": uint64(1000),
		"TradingFee":          uint64(1),
		"TraderAddressStr":    "12RqaTLErSnN88pGgXaKmw1PSQEaG86FA4uJsm32RZetAy7e5yEncqjTC6QJcMRjMfTSc48tcWRTyy8FoB9VkCHu56Vd9b86gd8Pq8k",
	}

	FeePerKb := 5

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		nil,
		FeePerKb,
		-1,
		metadata,
		"",
		0,
	}

	data, err := CreateAndSendTxWithPTokenTradeReq(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create and send normal tx %v\n", err)
		return
	}

	fmt.Printf("Send tx successfully - Data %v !!!", data)
}