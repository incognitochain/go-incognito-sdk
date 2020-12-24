package metadata

import (
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"github.com/pkg/errors"
	"strconv"
)

type WithDrawRewardRequest struct {
	privacy.PaymentAddress
	MetadataBase
	TokenID common.Hash
	Version int
}

func (withDrawRewardRequest WithDrawRewardRequest) Hash() *common.Hash {
	if withDrawRewardRequest.Version == 1 {
		bArr := append(withDrawRewardRequest.PaymentAddress.Bytes(), withDrawRewardRequest.TokenID.GetBytes()...)
		txReqHash := common.HashH(bArr)
		return &txReqHash
	} else {
		record := strconv.Itoa(withDrawRewardRequest.Type)
		data := []byte(record)
		hash := common.HashH(data)
		return &hash
	}
}

func NewWithDrawRewardRequestFromRPC(data map[string]interface{}) (Metadata, error) {
	metadataBase := MetadataBase{
		Type: WithDrawRewardRequestMeta,
	}

	requesterPaymentStr, ok := data["PaymentAddress"].(string)
	if !ok {
		return nil, errors.New("Invalid payment address receiver")
	}

	requestTokenID, ok := data["TokenID"].(string)
	if !ok {
		return nil, errors.New("Invalid token Id")
	}

	tokenID, err := common.Hash{}.NewHashFromStr(requestTokenID)
	if err != nil {
		return nil, err
	}

	requesterPublicKeySet, err := wallet.Base58CheckDeserialize(requesterPaymentStr)
	if err != nil {
		return nil, err
	}

	result := &WithDrawRewardRequest{
		MetadataBase:   metadataBase,
		PaymentAddress: requesterPublicKeySet.KeySet.PaymentAddress,
		TokenID:        *tokenID,
	}

	versionFloat, ok := data["Version"].(float64)
	if ok {
		version := int(versionFloat)
		result.Version = version
	}

	if ok, err := common.SliceExists(AcceptedWithdrawRewardRequestVersion, result.Version); !ok || err != nil {
		return nil, errors.Errorf("Invalid version %d", result.Version)
	}

	return result, nil
}
