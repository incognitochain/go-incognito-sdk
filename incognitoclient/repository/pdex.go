package repository

import (
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/constant"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/service"
	"github.com/pkg/errors"
	"strconv"
)

type Pdex struct {
	Inc        *service.IncogClient
	ConstantId string
	Block      *Block
}

func NewPdex(inc *service.IncogClient, constantId string, block *Block) *Pdex {
	return &Pdex{Inc: inc, ConstantId: constantId, Block: block}
}

func (p *Pdex) GetPDexState(beacon int32) (map[string]interface{}, error) {
	beaconParams := map[string]interface{}{"BeaconHeight": beacon}
	param := []interface{}{beaconParams}
	resp, _, err := p.Inc.PostAndReceiveInterface(constant.GetPdeState, param)
	if err != nil {
		return nil, errors.Wrap(err, "b.GetPdeState")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("couldn't get Result:  resp: %+v", data)
	}
	return result, nil
}

func (p *Pdex) TradePDex(privateKey string, buyTokenId string, tradingFee uint64, sellTokenId string, sellTokenAmount uint64, minimumAmount uint64, traderAddress string, networkFeeTokenID string, networkFee uint64) (string, error) {
	if sellTokenId == p.ConstantId {
		return p.SellPRVCrosspool(privateKey, buyTokenId, tradingFee, sellTokenAmount, minimumAmount, traderAddress)
	}

	return p.SellPTokenCrosspool(privateKey, buyTokenId, tradingFee, sellTokenId, sellTokenAmount, minimumAmount, traderAddress, networkFeeTokenID, networkFee)
}

func (p *Pdex) GetPDexTradeStatus(txId string) (constant.PDexTradeStatus, error) {
	param := map[string]interface{}{
		"TxRequestIDStr": txId,
	}

	paramArray := []interface{}{
		param,
	}

	resp, _, err := p.Inc.PostAndReceiveInterface(constant.GetPDETradeStatus, paramArray)
	if err != nil {
		return 0, errors.Wrapf(err, "p.blockchainAPI: param: %+v", paramArray)
	}
	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return 0, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return 0, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(float64)
	if !ok {
		return 0, errors.Errorf("couldn't get txID:  resp: %+v", data)
	}
	return constant.PDexTradeStatus(result), nil
}

func (p *Pdex) SellPTokenCrosspool(privateKey string, buyTokenId string, tradingFee uint64, sellTokenId string, sellTokenAmount uint64, minimumAmount uint64, traderAddress string, networkFeeTokenID string, networkFee uint64) (string, error) {
	var burningAddress string
	var err error
	burningAddress, err = p.Block.GetBurningAddress()
	var FeePerKb int
	var TokenFee uint64
	if err != nil {
		return "", errors.Wrap(err, "w.GetBurningAddress")
	}

	if networkFeeTokenID == p.ConstantId {
		FeePerKb = 5
		TokenFee = 0
	} else {
		FeePerKb = 0
		TokenFee = networkFee / constant.PDEX_TRADE_STEPS
	}

	metadata := map[string]interface{}{
		"Privacy":     true,
		"TokenID":     sellTokenId,
		"TokenTxType": 1,
		"TokenName":   "",
		"TokenSymbol": "",
		"TokenAmount": strconv.FormatUint(sellTokenAmount, 10),
		"TokenReceivers": map[string]interface{}{
			burningAddress: strconv.FormatUint(sellTokenAmount, 10),
		},
		"TokenFee":            strconv.FormatUint(TokenFee, 10),
		"TokenIDToBuyStr":     buyTokenId,
		"TokenIDToSellStr":    sellTokenId,
		"SellAmount":          strconv.FormatUint(sellTokenAmount, 10),
		"MinAcceptableAmount": strconv.FormatUint(minimumAmount, 10),
		"TradingFee":          strconv.FormatUint(tradingFee, 10),
		"TraderAddressStr":    traderAddress,
	}

	paramArray := []interface{}{
		privateKey,
		map[string]interface{}{
			burningAddress: strconv.FormatUint(tradingFee, 10),
		},
		FeePerKb,
		-1,
		metadata,
		"",
		0,
	}

	resp, _, err := p.Inc.PostAndReceiveInterface(constant.CreateAndSendTxWithPTokenCrosspolTradeReq, paramArray)
	if err != nil {
		return "", errors.Wrapf(err, "w.blockchainAPI: param: %+v", paramArray)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: req: %+v, resp: %+v", metadata, data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	if result["TxID"] == nil {
		return "", constant.ErrTxHashNotExists
	}

	return result["TxID"].(string), nil
}

func (p *Pdex) SellPRVCrosspool(privateKey string, buyTokenId string, tradingFee uint64, sellTokenAmount uint64, minimumAmount uint64, traderAddress string) (string, error) {
	var burningAddress string
	var err error
	burningAddress, err = p.Block.GetBurningAddress()

	if err != nil {
		return "", errors.Wrap(err, "w.GetBurningAddress")
	}

	metadata := map[string]interface{}{
		"TokenIDToBuyStr":     buyTokenId,
		"TokenIDToSellStr":    p.ConstantId,
		"SellAmount":          strconv.FormatUint(sellTokenAmount, 10),
		"MinAcceptableAmount": strconv.FormatUint(minimumAmount, 10),
		"TradingFee":          strconv.FormatUint(tradingFee, 10),
		"TraderAddressStr":    traderAddress,
	}

	paramArray := []interface{}{
		privateKey,
		map[string]interface{}{
			burningAddress: strconv.FormatUint(sellTokenAmount+tradingFee, 10),
		},
		1,
		-1,
		metadata,
	}

	resp, _, err := p.Inc.PostAndReceiveInterface(constant.CreateAndSendTxWithPRVCrosspollTradeReq, paramArray)
	if err != nil {
		return "", errors.Wrapf(err, "w.blockchainAPI: param: %+v", paramArray)
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: req: %+v, resp: %+v", metadata, data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get result: data: %+v", data["Result"])
	}
	if result["TxID"] == nil {
		return "", constant.ErrTxHashNotExists
	}
	return result["TxID"].(string), nil
}
