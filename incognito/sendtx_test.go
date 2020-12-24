package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestSendTx(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	var listPaymentAddresses = make(map[string]uint64)
	listPaymentAddresses["12S2rj2zV2cEGQyNt5Xgzcvg7W6dUEs8cfvsqT66wUQVGVkiFnf5YweRmCFQSLGRoSrKC34CFiZHMT9ABhH5FdSWeXkvrkHDf9Kbtjc"] = uint64(1000000)

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		listPaymentAddresses,
		5,
		0,
	}

	data, err := CreateAndSendTx(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create and send normal tx %v\n", err)
		return
	}

	fmt.Printf("Send tx successfully - Data %v !!!", data)
}

func TestSendTxPrivacy(t *testing.T) {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	var listPaymentAddresses = make(map[string]uint64)
	listPaymentAddresses["12RuEdPjq4yxivzm8xPxRVHmkL74t4eAdUKPdKKhMEnpxPH3k8GEyULbwq4hjwHWmHQr7MmGBJsMpdCHsYAqNE18jipWQwciBf9yqvQ"] = uint64(1000000)

	params := []interface{}{
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		listPaymentAddresses,
		5,
		1,//privacy
	}

	data, err := CreateAndSendTx(rpcClient, params)

	if err != nil {
		fmt.Printf("Error when create and send normal tx %v\n", err)
		return
	}

	fmt.Printf("Send tx successfully - Data %v !!!", data)
}
