package incognitoclient

import (
	"encoding/json"
	"fmt"
)

func (t *IncognitoTestSuite) TestGetBlockInfo() {
	result, _ := t.block.GetBlockInfo(1153596, 0)

	fmt.Println(result)
	t.NotEmpty(result)
}

func (t *IncognitoTestSuite) TestGetBlockChainInfo() {
	result, _ := t.block.GetChainInfo()

	fmt.Println(result)
	t.NotEmpty(result)
}

func (t *IncognitoTestSuite) TestGetBestBlockHeight() {
	result, _ := t.block.GetBestBlockHeight(0)

	fmt.Println(result)
	t.NotEmpty(result)
}

func (t *IncognitoTestSuite) TestGetBeaconHeight() {
	result, _ := t.block.GetBeaconHeight()

	fmt.Println(result)
	t.NotEmpty(result)
}

func (t *IncognitoTestSuite) TestGetBeaconBestStateDetail() {
	data, err := t.block.GetBeaconBestStateDetail()
	out, err := json.Marshal(data.Result.ShardCommittee)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
	t.NotEmpty(string(out))
}

func (t *IncognitoTestSuite) TestGetBurnToken() {
	tx, err := t.block.GetBurningAddress()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(tx)
}
