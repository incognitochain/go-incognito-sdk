package incognitoclient

import (
	"encoding/json"
	"fmt"
	"math/big"
)

func (t *IncognitoTestSuite) TestListPrivacyCustomToken() {
	result, _ := t.wallet.ListPrivacyCustomToken()

	t.NotEmpty(result)
	fmt.Println(result)
}

func (t *IncognitoTestSuite) TestSendPrvPrivacy() {
	tx, err := t.wallet.SendToken(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"12Rsf3wFnThr3T8dMafmaw4b3CzUatNao61dkj8KyoHfH5VWr4ravL32sunA2z9UhbNnyijzWFaVDvacJPSRFAq66HU7YBWjwfWR7Ff",
		t.client.GetPRVToken(),
		500000000000,
		5,
		"")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(tx)
}

func (t *IncognitoTestSuite) TestSendPTokenWithFeePrv() {
	tx, err := t.wallet.SendToken(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"12RqaTLErSnN88pGgXaKmw1PSQEaG86FA4uJsm32RZetAy7e5yEncqjTC6QJcMRjMfTSc48tcWRTyy8FoB9VkCHu56Vd9b86gd8Pq8k",
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
		9800000,
		5,
		"0000000000000000000000000000000000000000000000000000000000000004")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(tx)
}

func (t *IncognitoTestSuite) TestSendPTokenWithFeePtoken() {
	tx1, err1 := t.wallet.SendToken(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"12RqaTLErSnN88pGgXaKmw1PSQEaG86FA4uJsm32RZetAy7e5yEncqjTC6QJcMRjMfTSc48tcWRTyy8FoB9VkCHu56Vd9b86gd8Pq8k",
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854",
		9800000,
		100,
		"ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854")

	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}

	fmt.Println(tx1)
}

func (t *IncognitoTestSuite) TestGetBalance() {
	//native coin
	amountPrv, err := t.wallet.GetBalance("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", "0000000000000000000000000000000000000000000000000000000000000004") //prv
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(amountPrv)

	//pToken
	amountEth, err := t.wallet.GetBalance("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854") //eth
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(amountEth)
}

func (t *IncognitoTestSuite) TestTransactionByReceivers() {
	//native coin
	receiveTx, err := t.wallet.GetTransactionByReceiversAddress("12RtT2tTRLfSFZpLo5xjjcUnydDS7QL9zKdVNyRZnnh3q2mQavkb5P7G62hHeYTLrwv6wapc8f2MS6KGaWpG7ossDrnWvBd5gUMgPra", "13hW3CGNB8jqGHUpRA96aJTiuMFdKRKg6XqJdx69sjpHkQhowM5otRd9xgF3UzPeApQhkea4zY2VAJAbmKwPumHDhT4a4g41jDP4B5u") //prv
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resultInBytes, err := json.Marshal(receiveTx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(resultInBytes))

	t.NotEmpty(receiveTx)
}

func (t *IncognitoTestSuite) TestGetTransactionAmount() {
	amount, err := t.wallet.GetTransactionAmount("fd92a9a4e99945810b8429489e2b9f74a5e1d15d90b7a8987b4afadbb3d792b5", "12RtT2tTRLfSFZpLo5xjjcUnydDS7QL9zKdVNyRZnnh3q2mQavkb5P7G62hHeYTLrwv6wapc8f2MS6KGaWpG7ossDrnWvBd5gUMgPra", "13hW3CGNB8jqGHUpRA96aJTiuMFdKRKg6XqJdx69sjpHkQhowM5otRd9xgF3UzPeApQhkea4zY2VAJAbmKwPumHDhT4a4g41jDP4B5u")
	fmt.Println(amount)
	fmt.Println(err)

	t.Equal(amount, 446239992496)
}

func (t *IncognitoTestSuite) TestCreateWalletAddressByShardId() {
	fmt.Println(t.wallet.CreateWallet())
	fmt.Println(t.wallet.CreateWalletByShardId(6))
}

func (t *IncognitoTestSuite) TestDefragmentationPrv() {
	value, err := t.wallet.Defragmentation("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", int64(500000*1e9), PRVToken)

	fmt.Println("Tx: ", value)
	fmt.Println(err)
}

func (t *IncognitoTestSuite) TestDefragmentationPToken() {
	value, err := t.wallet.Defragmentation("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", int64(500000*1e9), "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854")

	fmt.Println("Tx: ", value)
	fmt.Println(err)
}

func (t *IncognitoTestSuite) TestGetUTXO() {
	value, err := t.wallet.GetUTXO("112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5", "0000000000000000000000000000000000000000000000000000000000000004")

	fmt.Println(value)
	fmt.Println(err)
}

func (t *IncognitoTestSuite) TestMintCentralizedToken() {
	privateKey := "112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5"
	receiveAddress := "12RwamF5njyL5cqpiMZ3SrqGHMqDaEDLyQexeaHYjYn2LDMzKZzgPZHnbQ75iLBKxm4md4kiyLxrPrFRNRNNktmAMjmfD4ktmcptgiX"

	depositedAmount := big.NewInt(2000000000)
	tokenId := "4584d5e9b2fc0337dfb17f4b5bb025e5b82c38cfa4f54e8a3d4fcdd03954ff82"
	tokenName := "BTC"

	tx, err := t.wallet.MintCentralizedToken(privateKey, receiveAddress, depositedAmount, tokenId, tokenName)

	fmt.Println(err)
	fmt.Println(tx)

	t.NotEmpty(tx)
}

func (t *IncognitoTestSuite) TestBurnCentralizedToken() {
	privateKey := "112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5"
	//-1: auto prv fee, 0: 0 prv fee, 1: ptoken fee
	feePRV := -1

	burnAddress, _ := t.block.GetBurningAddress()

	metadata := map[string]interface{}{
		"TokenID":     "4584d5e9b2fc0337dfb17f4b5bb025e5b82c38cfa4f54e8a3d4fcdd03954ff82",
		"Privacy":     true,
		"TokenTxType": 1,
		"TokenName":   "BTC",
		"TokenSymbol": "pBTC",
		"TokenAmount": uint64(2000000000),
		"TokenReceivers": map[string]uint64{
			burnAddress: uint64(2000000000),
		},
		"TokenFee": uint64(0),
	}

	tx, err := t.wallet.BurnCentralizedToken(privateKey, feePRV, metadata)

	fmt.Println(err)
	fmt.Println(tx)
	t.NotEmpty(tx)
}

func (t *IncognitoTestSuite) TestGetMintStatusCentralized() {
	status, err := t.wallet.GetMintStatusCentralized("36a9aa575e910754f49aa739f5b69fd78b81da58a610940dd08f223e95e5286d")
	fmt.Println(err)
	fmt.Println(status)

	t.Equal(2, status)
}

func (t *IncognitoTestSuite) TestGenerateTokenID() {
	result, err := t.wallet.GenerateTokenID("ABC", "pABC")
	fmt.Println(err)
	fmt.Println(result)

	t.NotEmpty(result)
}
