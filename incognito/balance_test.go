package incognito

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"testing"
)

func TestGetBalance(t *testing.T)  {
	rpcClient := rpcclient.NewHttpClient("https://testnet.incognito.org/fullnode", "https", "testnet.incognito.org/fullnode", 0)

	balance, err := GetBalance(rpcClient,
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"880ea0787f6c1555e59e3958a595086b7802fc7a38276bcd80d4525606557fbc")

	if err != nil {
		fmt.Printf("Error when create raw data: %v\n", err)
		return
	}

	fmt.Println("Balance:", balance)
}
