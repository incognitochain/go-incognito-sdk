package repository

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/constant"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/entity"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/service"
	"github.com/pkg/errors"
)

type Block struct {
	Inc *service.IncogClient
}

func NewBlock(inc *service.IncogClient) *Block {
	return &Block{Inc: inc}
}

func (b *Block) GetBlockInfo(blockHeight int32, shardID int) (*entity.GetBlockInfo, error) {
	param := []interface{}{blockHeight, shardID, "2"}

	resp, _, err := b.Inc.PostAndReceiveInterface(constant.Retrieveblockbyheight, param)

	if err != nil {
		return nil, err
	}
	data := resp.(map[string]interface{})

	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response: resp: %+v", data)
	}
	var result entity.GetBlockInfo
	resultResp := data["Result"].([]interface{})
	resultRespStr, err := json.Marshal(resultResp[0])
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resultRespStr, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *Block) GetBlockChainInfo() (*entity.GetBlockChainInfoResult, error) {
	param := []interface{}{}
	resp, _, err := b.Inc.PostAndReceiveInterface(constant.Getblockchaininfo, param)

	if err != nil {
		return nil, err
	}
	data := resp.(map[string]interface{})
	resultResp := data["Result"]

	result := entity.GetBlockChainInfoResult{}

	if resultResp == nil {
		return nil, errors.New("Fail")
	}

	resultInBytes, err := json.MarshalIndent(resultResp, "", "\t")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resultInBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBestBlockHeight - get height of the highest block. it could be either shard or beacon.
// Note: it would return the highest beacon block height if shardID = -1
func (b *Block) GetBestBlockHeight(shardID int) (uint64, error) {
	param := []interface{}{shardID}
	resp, _, err := b.Inc.PostAndReceiveInterface(constant.GetBlockCount, param)
	if err != nil {
		return 0, errors.Wrapf(err, "b.blockchainAPI: param: %+v", param)
	}

	data, ok := resp.(map[string]interface{})

	if !ok {
		return 0, errors.Errorf("Wrong response format, it should be an object that contains Error or Result field")
	}
	if data["Error"] != nil {
		return 0, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return 0, nil
	}
	blockHeight, ok := data["Result"].(float64)
	if !ok {
		return 0, errors.Errorf("Wrong response format, the returned block height should be a number.")
	}
	return uint64(blockHeight), nil
}

func (b *Block) GetBeaconHeight() (int32, error) {
	blockChainInfo, err := b.GetBlockChainInfo()

	if err != nil {
		return 0, errors.Wrap(err, "b.GetBlockChainInfo")
	}

	beaconHeight := blockChainInfo.BestBlocks["-1"].Height

	return beaconHeight, nil
}

func (b *Block) GetBeaconBestStateDetail() (res *entity.BeaconBestStateResp, err error) {
	empty := []string{}
	body, err := b.Inc.Post(constant.GetBeaconBestStateDetail, empty)
	if err != nil {
		return nil, err
	}

	var resp entity.BeaconBestStateResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (b *Block) GetBurningAddress() (string, error) {
	beaconHeight, err := b.GetBeaconHeight()

	if err != nil {
		return "", errors.Wrap(err, "w.GetBeaconHeight")
	}

	paramArray := []interface{}{
		beaconHeight,
	}

	resp, _, err := b.Inc.PostAndReceiveInterface(constant.GetBurningAddress, paramArray)
	if err != nil {
		return "", errors.Wrapf(err, "w.blockchainAPI: param: %+v", paramArray)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: req: %+v, resp: %+v", paramArray, data)
	}
	result, ok := data["Result"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: param: %+v, resp: %+v", paramArray, data)
	}
	return result, nil
}
