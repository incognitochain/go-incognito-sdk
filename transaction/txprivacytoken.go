package transaction

import (
	"encoding/json"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"github.com/pkg/errors"
)

// TxCustomTokenPrivacy is class tx which is inherited from P tx(supporting privacy) for fee
// and contain data(with supporting privacy format) to support issuing and transfer a custom token(token from end-user, look like erc-20)
// Dev or end-user can use this class tx to create an token type which use personal purpose
// TxCustomTokenPrivacy is an advance format of TxNormalToken
// so that user need to spend a lot fee to create this class tx
type TxCustomTokenPrivacy struct {
	Tx                                    // inherit from normal tx of P(supporting privacy) with a high fee to ensure that tx could contain a big data of privacy for token
	TxPrivacyTokenData TxPrivacyTokenData `json:"TxTokenPrivacyData"` // supporting privacy format
	// private field, not use for json parser, only use as temp variable
	cachedHash *common.Hash // cached hash data of tx
}

type TxPrivacyTokenInitParams struct {
	senderKey       *privacy.PrivateKey
	paymentInfo     []*privacy.PaymentInfo
	inputCoin       []*privacy.InputCoin
	outputCoin      []*privacy.OutputCoin
	feeNativeCoin   uint64
	tokenParams     *CustomTokenPrivacyParamTx
	metaData        metadata.Metadata
	hasPrivacyCoin  bool
	hasPrivacyToken bool
	shardID         byte
	info            []byte
}

func NewTxPrivacyTokenInitParams(
	senderKey *privacy.PrivateKey,
	paymentInfo []*privacy.PaymentInfo,
	inputCoin []*privacy.InputCoin,
	outputCoin []*privacy.OutputCoin,
	feeNativeCoin uint64,
	tokenParams *CustomTokenPrivacyParamTx,
	metaData metadata.Metadata,
	hasPrivacyCoin bool,
	hasPrivacyToken bool,
	shardID byte,
	info []byte) *TxPrivacyTokenInitParams {
	params := &TxPrivacyTokenInitParams{
		shardID:         shardID,
		paymentInfo:     paymentInfo,
		metaData:        metaData,
		feeNativeCoin:   feeNativeCoin,
		hasPrivacyCoin:  hasPrivacyCoin,
		hasPrivacyToken: hasPrivacyToken,
		inputCoin:       inputCoin,
		outputCoin:      outputCoin,
		senderKey:       senderKey,
		tokenParams:     tokenParams,
		info:            info,
	}
	return params
}

func (txCustomTokenPrivacy *TxCustomTokenPrivacy) UnmarshalJSON(data []byte) error {
	tx := Tx{}
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return err
	}
	temp := &struct {
		TxTokenPrivacyData TxPrivacyTokenData
	}{}
	err = json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	TxTokenPrivacyDataJson, err := json.MarshalIndent(temp.TxTokenPrivacyData, "", "\t")
	if err != nil {
		return err
	}

	err = json.Unmarshal(TxTokenPrivacyDataJson, &txCustomTokenPrivacy.TxPrivacyTokenData)
	if err != nil {
		return err
	}

	txCustomTokenPrivacy.Tx = tx
	return nil
}

func (txCustomTokenPrivacy TxCustomTokenPrivacy) String() string {
	// get hash of tx
	record := txCustomTokenPrivacy.Tx.Hash().String()
	// add more hash of tx custom token data privacy
	tokenPrivacyDataHash, _ := txCustomTokenPrivacy.TxPrivacyTokenData.Hash()
	record += tokenPrivacyDataHash.String()
	if txCustomTokenPrivacy.Metadata != nil {
		record += string(txCustomTokenPrivacy.Metadata.Hash()[:])
	}
	return record
}

func (txCustomTokenPrivacy TxCustomTokenPrivacy) JSONString() (string, error) {
	data, err := json.MarshalIndent(txCustomTokenPrivacy, "", "\t")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Hash returns the hash of all fields of the transaction
func (txCustomTokenPrivacy *TxCustomTokenPrivacy) Hash() *common.Hash {
	if txCustomTokenPrivacy.cachedHash != nil {
		return txCustomTokenPrivacy.cachedHash
	}
	// final hash
	hash := common.HashH([]byte(txCustomTokenPrivacy.String()))
	return &hash
}

// Init -  build normal tx component and privacy custom token data
func (txCustomTokenPrivacy *TxCustomTokenPrivacy) Init(params *TxPrivacyTokenInitParams, client *rpcclient.HttpClient, keyWallet *wallet.KeyWallet) error {
	var err error
	// init data for tx PRV for fee
	normalTx := Tx{}
	err = normalTx.Init(
		NewTxPrivacyInitParams(
			params.senderKey,
			params.paymentInfo,
			params.inputCoin,
			params.outputCoin,
			params.feeNativeCoin,
			params.hasPrivacyCoin,
			nil,
			params.metaData,
			params.info,
		),
		client,
		keyWallet,
	)

	if err != nil {
		return err
	}

	// override TxCustomTokenPrivacyType type
	normalTx.Type = common.TxCustomTokenPrivacyType
	txCustomTokenPrivacy.Tx = normalTx

	// check tx size
	limitFee := uint64(0)
	estimateTxSizeParam := NewEstimateTxSizeParam(len(params.inputCoin), len(params.paymentInfo),
		params.hasPrivacyCoin, nil, params.tokenParams, limitFee)

	if txSize := EstimateTxSize(estimateTxSizeParam); txSize > common.MaxTxSize {
		return errors.New("Size of tx info exceed max size of tx")
	}

	// check action type and create privacy custom toke data
	var handled = false
	// Add token data component
	switch params.tokenParams.TokenTxType {
	case CustomTokenTransfer:
		{
			handled = true
			// make a transfering for privacy custom token
			// fee always 0 and reuse function of normal tx for custom token ID
			temp := Tx{}
			propertyID, _ := common.Hash{}.NewHashFromStr(params.tokenParams.PropertyID)

			txCustomTokenPrivacy.TxPrivacyTokenData = TxPrivacyTokenData{
				Type:           params.tokenParams.TokenTxType,
				PropertyName:   params.tokenParams.PropertyName,
				PropertySymbol: params.tokenParams.PropertySymbol,
				PropertyID:     *propertyID,
				Mintable:       params.tokenParams.Mintable,
			}

			err := temp.Init(
				NewTxPrivacyInitParams(
					params.senderKey,
					params.tokenParams.Receiver,
					params.tokenParams.TokenInput,
					params.tokenParams.TokenOutput,
					params.tokenParams.Fee,
					params.hasPrivacyToken,
					propertyID,
					nil,
					nil,
				),
				client,
				keyWallet,
			)

			if err != nil {
				return err
			}

			txCustomTokenPrivacy.TxPrivacyTokenData.TxNormal = temp
		}
	}

	if !handled {
		return errors.New("can't handle this TokenTxType")
	}

	return nil
}
