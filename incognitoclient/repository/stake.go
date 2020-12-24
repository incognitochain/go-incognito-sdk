package repository

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/constant"

	"github.com/incognitochain/go-incognito-sdk/incognitoclient/entity"
	"github.com/incognitochain/go-incognito-sdk/incognitoclient/service"
	"github.com/pkg/errors"
)

type Stake struct {
	Inc                 *service.IncogClient
	IncChainIntegration *IncChainIntegration
}

func NewStake(inc *service.IncogClient, incChainIntegration *IncChainIntegration) *Stake {
	return &Stake{Inc: inc, IncChainIntegration: incChainIntegration}
}

func (s *Stake) ListUnstake() ([]entity.Unstake, error) {
	param := []interface{}{}
	resp, _, err := s.Inc.PostAndReceiveInterface(constant.GetBeaconBestStateDetail, param)
	if err != nil {
		return nil, errors.Wrap(err, "b.ListUnstake")
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
	// autoStake, ok := result["AutoStaking"].(map[string]interface{})
	// if !ok {
	// 	return nil, errors.Errorf("couldn't get AutoStaking:  resp: %+v",  data)
	// }

	// fmt.Println(autoStake)
	var AutoStake []entity.Unstake
	r, err := json.Marshal(result["AutoStaking"])
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	if err := json.Unmarshal(r, &AutoStake); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return AutoStake, nil

	// return nil, nil
}
func (s *Stake) GetTotalStaker() (float64, error) {
	param := []interface{}{}
	resp, _, err := s.Inc.PostAndReceiveInterface(constant.GetTotalStaker, param)
	if err != nil {
		return 0, errors.Wrap(err, "b.GetTotalStaker")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return 0, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return 0, errors.Errorf("couldn't get result from response:  resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	fmt.Println(result)
	if !ok {
		return 0, errors.Errorf("couldn't get Result:  resp: %+v", data)
	}

	autoStake, ok := result["TotalStaker"].(float64)
	if !ok {
		return 0, errors.Errorf("couldn't get TotalStaker: result: %+v", result)
	}
	fmt.Println(autoStake)
	return autoStake, nil
}

func (b *Stake) Staking(receiveRewardAddress, privateKey, userPaymentAddress, userValidatorKey, burnTokenAddress string) (string, error) {

	amountToStake := uint64(1750000000000)

	param := []interface{}{
		privateKey,
		map[string]uint64{burnTokenAddress: amountToStake},
		5,
		0,
		map[string]interface{}{
			"StakingType":                  63,
			"CandidatePaymentAddress":      userPaymentAddress,
			"PrivateSeed":                  userValidatorKey,
			"RewardReceiverPaymentAddress": receiveRewardAddress,
			"AutoReStaking":                true,
		},
	}

	//rpc: CreateAndSendStakingTransaction
	rawData, err := b.IncChainIntegration.CreateAndSendStakingTx(param)
	if err != nil {
		return "", errors.Wrap(err, "w.CreateAndSendStakingTx")
	}

	fmt.Printf("raw data method CreateAndSendStakingTx: %v \n", rawData)

	resp, _, err := b.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)
	if err != nil {
		return "", errors.Wrap(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (b *Stake) Unstaking(privateKey, userPaymentAddress, userValidatorKey, burnTokenAddress string) (string, error) {
	param := []interface{}{
		privateKey,
		map[string]uint64{burnTokenAddress: 0},
		10, // fee 10 nano prv
		0,
		map[string]interface{}{
			"StopAutoStakingType":     127,
			"CandidatePaymentAddress": userPaymentAddress,
			"PrivateSeed":             userValidatorKey,
		},
	}

	//rpc: CreateAndSendUnStakingTransaction
	rawData, err := b.IncChainIntegration.CreateAndSendStopAutoStakingTransaction(param)
	if err != nil {
		return "", errors.Wrap(err, "w.CreateAndSendStopAutoStakingTransaction")
	}

	fmt.Printf("raw data method CreateAndSendStopAutoStakingTransaction: %v \n", rawData)

	resp, _, err := b.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)
	if err != nil {
		return "", errors.Wrap(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID: resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (b *Stake) WithDrawReward(privateKey, paymentAddress, tokenID string) (string, error) {
	param := []interface{}{
		privateKey,
		nil,
		0,
		0,
		map[string]interface{}{
			"PaymentAddress": paymentAddress,
			"TokenID":        tokenID,
		},
	}

	//rpc: WithDrawReward
	rawData, err := b.IncChainIntegration.CreateAndSendWithDrawTransaction(param)
	if err != nil {
		return "", errors.Wrap(err, "w.CreateAndSendWithDrawTransaction")
	}

	fmt.Printf("raw data method CreateAndSendWithDrawTransaction: %v \n", rawData)

	resp, _, err := b.Inc.PostAndReceiveInterface(constant.SendRawTransaction, rawData)
	if err != nil {
		return "", errors.Wrap(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return "", errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return "", errors.Errorf("couldn't get result from response: resp: %+v", data)
	}
	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return "", errors.Errorf("couldn't get txID: resp: %+v", data)
	}
	txID, ok := result["TxID"].(string)
	if !ok {
		return "", errors.Errorf("couldn't get txID: result: %+v", result)
	}
	return txID, nil
}

func (b *Stake) GetRewardAmount(paymentAddress string) ([]entity.RewardItems, error) {
	param := []interface{}{
		paymentAddress,
	}

	resp, _, err := b.Inc.PostAndReceiveInterface(constant.RewardAmount, param)
	if err != nil {
		return nil, errors.Wrapf(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return nil, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}
	if data["Result"] == nil {
		return nil, errors.Errorf("couldn't get result from response: resp: %+v", data)
	}

	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("couldn't get result from response: resp: %+v", data)
	}

	var rewards []entity.RewardItems
	for s, k := range result {
		amount := k.(float64)
		rewards = append(rewards, entity.RewardItems{s, amount})
	}

	return rewards, nil
}

func (b *Stake) GetNodeAvailable(validatorKey string) (float64, error) {
	param := []interface{}{
		validatorKey,
	}

	resp, _, err := b.Inc.PostAndReceiveInterface(constant.RoleByValidatorKey, param)
	if err != nil {
		return 0, errors.Wrapf(err, "b.blockchainAPI")
	}

	data := resp.(map[string]interface{})
	if data["Error"] != nil {
		return 0, errors.Errorf("couldn't get result from response data: %+v", data["Error"])
	}

	if data["Result"] == nil {
		return 0, errors.Errorf("couldn't get result from response: resp: %+v", data)
	}

	result, ok := data["Result"].(map[string]interface{})
	if !ok {
		return 0, errors.Errorf("couldn't get result from response: resp: %+v", data)
	}

	role, ok := result["Role"].(float64)
	if !ok {
		return 0, errors.Errorf("couldn't get txID: result: %+v", result)
	}

	return role, nil
}

// Dung.Dang: for PRV only :(
func (w *Stake) ListRewardAmounts() ([]entity.RewardAmount, error) {
	param := []interface{}{}
	resp, _, err := w.Inc.PostAndReceiveInterface(constant.ListRewardAmount, param)
	if err != nil {
		return nil, errors.Wrap(err, "w.ListRewardAmounts")
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
		return nil, errors.Errorf("couldn't get Result: resp: %+v", data)
	}
	var rewards []entity.RewardAmount

	for s, k := range result {
		reward := k.(map[string]interface{})
		if reward["0000000000000000000000000000000000000000000000000000000000000004"] == nil {
			continue
		}
		amount := reward["0000000000000000000000000000000000000000000000000000000000000004"].(float64)
		rewards = append(rewards, entity.RewardAmount{s, amount})
	}

	return rewards, nil
}
