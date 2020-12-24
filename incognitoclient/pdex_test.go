package incognitoclient

import "fmt"

func (t *IncognitoTestSuite) TestGetPDexState() {
	result, err := t.pdex.GetPDexState(1153596)
	fmt.Println(result)
	fmt.Println(err)

	t.NotEmpty(result)
}

func (t *IncognitoTestSuite) TestTradePdex() {
	result, err := t.pdex.TradePDex(
		"112t8s4Pdng512MhHmLVJNYqzoEJQ1TG4XZduvjfwYZFJhmuNtGPhUYRko4jSPFBFmeRg6bumKQuhAEMriQ72cpp5SKAkRuXfLCv5xeZx3f5",
		"c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0",
		uint64(100),
		t.client.GetPRVToken(),
		uint64(1000000000),
		uint64(0),
		"12RwamF5njyL5cqpiMZ3SrqGHMqDaEDLyQexeaHYjYn2LDMzKZzgPZHnbQ75iLBKxm4md4kiyLxrPrFRNRNNktmAMjmfD4ktmcptgiX",
		t.client.GetPRVToken(),
		uint64(100))

	fmt.Println(err)
	fmt.Println(result)

	t.NotEmpty(result)
}
