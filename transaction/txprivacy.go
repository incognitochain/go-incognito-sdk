package transaction

import (
	"encoding/base64"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk/common"
	"github.com/incognitochain/go-incognito-sdk/common/base58"
	"github.com/incognitochain/go-incognito-sdk/incognitokey"
	"github.com/incognitochain/go-incognito-sdk/metadata"
	"github.com/incognitochain/go-incognito-sdk/privacy"
	"github.com/incognitochain/go-incognito-sdk/privacy/zkp"
	"github.com/incognitochain/go-incognito-sdk/rpcclient"
	"github.com/incognitochain/go-incognito-sdk/wallet"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
	"time"
)

type Tx struct {
	// Basic data, required
	Version  int8   `json:"Version"`
	Type     string `json:"Type"` // Transaction type
	LockTime int64  `json:"LockTime"`
	Fee      uint64 `json:"Fee"` // Fee applies: always consant
	Info     []byte // 512 bytes
	// Sign and Privacy proof, required
	SigPubKey            []byte `json:"SigPubKey, omitempty"` // 33 bytes
	Sig                  []byte `json:"Sig, omitempty"`       //
	Proof                *zkp.PaymentProof
	PubKeyLastByteSender byte
	// Metadata, optional
	Metadata metadata.Metadata
	// private field, not use for json parser, only use as temp variable
	sigPrivKey       []byte       // is ALWAYS private property of struct, if privacy: 64 bytes, and otherwise, 32 bytes
	cachedHash       *common.Hash // cached hash data of tx
	cachedActualSize *uint64      // cached actualsize data for tx
}

type TxPrivacyInitParams struct {
	senderSK    *privacy.PrivateKey
	paymentInfo []*privacy.PaymentInfo
	inputCoins  []*privacy.InputCoin
	outputCoins []*privacy.OutputCoin
	fee         uint64
	hasPrivacy  bool
	tokenID     *common.Hash // default is nil -> use for prv coin
	metaData    metadata.Metadata
	info        []byte // 512 bytes
}

func NewTxPrivacyInitParams(
	senderSK *privacy.PrivateKey,
	paymentInfo []*privacy.PaymentInfo,
	inputCoins []*privacy.InputCoin,
	outputCoins []*privacy.OutputCoin,
	fee uint64,
	hasPrivacy bool,
	tokenID *common.Hash, // default is nil -> use for prv coin
	metaData metadata.Metadata,
	info []byte) *TxPrivacyInitParams {
	params := &TxPrivacyInitParams{
		tokenID:     tokenID,
		hasPrivacy:  hasPrivacy,
		inputCoins:  inputCoins,
		outputCoins: outputCoins,
		fee:         fee,
		metaData:    metaData,
		paymentInfo: paymentInfo,
		senderSK:    senderSK,
		info:        info,
	}
	return params
}

func (tx *Tx) Init(params *TxPrivacyInitParams, client *rpcclient.HttpClient, keyWallet *wallet.KeyWallet) error {
	tx.Version = txVersion
	var err error
	if len(params.inputCoins) > 255 {
		return errors.New("Input coin is very larger")
	}
	if len(params.paymentInfo) > 254 {
		return errors.New("Payment info is very larger")
	}
	limitFee := uint64(0)
	estimateTxSizeParam := NewEstimateTxSizeParam(
		len(params.inputCoins),
		len(params.paymentInfo),
		params.hasPrivacy,
		nil,
		nil,
		limitFee,
	)

	if txSize := EstimateTxSize(estimateTxSizeParam); txSize > common.MaxTxSize {
		return errors.New(fmt.Sprintf("Estimate tx size overload, maximum = %v", common.MaxTxSize))
	}

	if params.tokenID == nil {
		// using default PRV
		params.tokenID = &common.Hash{}
		err := params.tokenID.SetBytes(common.PRVCoinID[:])
		if err != nil {
			return errors.Wrap(err, "params.tokenID.SetBytes")
		}
	}

	// Calculate execution time
	start := time.Now()

	if tx.LockTime == 0 {
		tx.LockTime = time.Now().Unix()
	}

	// create sender's key set from sender's spending key
	senderFullKey := incognitokey.KeySet{}
	err = senderFullKey.InitFromPrivateKey(params.senderSK)
	if err != nil {
		return errors.Wrap(err, "senderFullKey.InitFromPrivateKey")
	}
	// get public key last byte of sender
	pkLastByteSender := senderFullKey.PaymentAddress.Pk[len(senderFullKey.PaymentAddress.Pk)-1]

	// init info of tx
	tx.Info = []byte{}
	lenTxInfo := len(params.info)

	if lenTxInfo > 0 {
		if lenTxInfo > MaxSizeInfo {
			return errors.New(fmt.Sprintf("Len Tx Info overload, maximum = %v", MaxSizeInfo))
		}

		tx.Info = params.info
	}

	// set metadata
	tx.Metadata = params.metaData

	// set tx type
	tx.Type = common.TxNormalType

	fmt.Println(fmt.Sprintf("==> Init Tx len(inputCoins) = %v fee = %v hasPrivacy = %v", len(params.inputCoins), params.fee, params.hasPrivacy))

	if len(params.inputCoins) == 0 && params.fee == 0 && !params.hasPrivacy {
		tx.Fee = params.fee
		tx.sigPrivKey = *params.senderSK
		tx.PubKeyLastByteSender = common.GetShardIDFromLastByte(pkLastByteSender)
		err := tx.signTx()
		if err != nil {
			return errors.Wrap(err, "tx.signTx")
		}
		return nil
	}

	var commitmentIndexs []uint64   // array index random of commitments in transactionStateDB
	var myCommitmentIndexs []uint64 // index in array index random of commitment in transactionStateDB
	var commitments []string

	if params.hasPrivacy {
		if len(params.inputCoins) == 0 {
			return errors.New("Input coins is empty")
		}

		paymentAddrStr := keyWallet.Base58CheckSerialize(wallet.PaymentAddressType)
		commitmentIndexs, myCommitmentIndexs, commitments, err = rpcclient.RandomCommitmentsProcess(client, params.outputCoins, paymentAddrStr, params.tokenID)

		// Check number of list of random commitments, list of random commitment indices
		if len(commitmentIndexs) != len(params.inputCoins)*privacy.CommitmentRingSize {
			return errors.New("Random commitments")
		}

		if len(myCommitmentIndexs) != len(params.inputCoins) {
			return errors.New("Number of list my commitment indices must be equal to number of input coins")
		}
	}

	// Calculate execution time for creating payment proof
	startPrivacy := time.Now()

	// Calculate sum of all output coins' value
	sumOutputValue := uint64(0)
	for _, p := range params.paymentInfo {
		sumOutputValue += p.Amount
	}

	// Calculate sum of all input coins' value
	sumInputValue := uint64(0)
	for _, coin := range params.inputCoins {
		sumInputValue += coin.CoinDetails.GetValue()
	}

	// Calculate over balance, it will be returned to sender
	overBalance := int64(sumInputValue - sumOutputValue - params.fee)

	// Check if sum of input coins' value is at least sum of output coins' value and tx fee
	if overBalance < 0 {
		return errors.New(fmt.Sprintf("Input value less than output value. sumInputValue=%d sumOutputValue=%d fee=%d", sumInputValue, sumOutputValue, params.fee))
	}

	// if overBalance > 0, create a new payment info with pk is sender's pk and amount is overBalance
	if overBalance > 0 {
		changePaymentInfo := new(privacy.PaymentInfo)
		changePaymentInfo.Amount = uint64(overBalance)
		changePaymentInfo.PaymentAddress = senderFullKey.PaymentAddress
		params.paymentInfo = append(params.paymentInfo, changePaymentInfo)
	}

	// create new output coins
	outputCoins := make([]*privacy.OutputCoin, len(params.paymentInfo))

	// create SNDs for output coins
	ok := true
	sndOuts := make([]*privacy.Scalar, 0)

	for ok {
		for i := 0; i < len(params.paymentInfo); i++ {
			sndOut := privacy.RandomScalar()
			for {
				keyWalletTmp := new(wallet.KeyWallet)
				keyWalletTmp.KeySet.PaymentAddress = params.paymentInfo[i].PaymentAddress
				paymentAddrTmpStr := keyWalletTmp.Base58CheckSerialize(wallet.PaymentAddressType)

				ok1, err := rpcclient.CheckSNDerivatorExistence(client, paymentAddrTmpStr, []*privacy.Scalar{sndOut})
				if err != nil {
					fmt.Println(errors.Wrap(err, "rpcclient.CheckSNDerivatorExistence").Error())
				}
				// if sndOut existed, then re-random it
				if ok1[0] {
					sndOut = privacy.RandomScalar()
				} else {
					//fmt.Println("break 3 RandomScalar")
					break
				}
			}

			sndOuts = append(sndOuts, sndOut)
		}

		// if sndOuts has two elements that have same value, then re-generates it
		ok = privacy.CheckDuplicateScalarArray(sndOuts)
		if ok {
			sndOuts = make([]*privacy.Scalar, 0)
		}
	}

	// create new output coins with info: Pk, value, last byte of pk, snd
	for i, pInfo := range params.paymentInfo {
		outputCoins[i] = new(privacy.OutputCoin)
		outputCoins[i].CoinDetails = new(privacy.Coin)
		outputCoins[i].CoinDetails.SetValue(pInfo.Amount)
		if len(pInfo.Message) > 0 {
			if len(pInfo.Message) > privacy.MaxSizeInfoCoin {
				return errors.New(fmt.Sprintf("Len pInfo.Message is overload, maximum = %v", privacy.MaxSizeInfoCoin))
			}
		}
		outputCoins[i].CoinDetails.SetInfo(pInfo.Message)

		PK, err := new(privacy.Point).FromBytesS(pInfo.PaymentAddress.Pk)
		if err != nil {
			return errors.Wrap(err, "DecompressPaymentAddress")
		}
		outputCoins[i].CoinDetails.SetPublicKey(PK)
		outputCoins[i].CoinDetails.SetSNDerivator(sndOuts[i])
	}

	// assign fee tx
	tx.Fee = params.fee

	// create zero knowledge proof of payment
	tx.Proof = &zkp.PaymentProof{}

	// get list of commitments for proving one-out-of-many from commitmentIndexs
	commitmentProving := make([]*privacy.Point, len(commitments))
	for i, cmIndex := range commitments {
		temp, _, err := base58.Base58Check{}.Decode(cmIndex)
		if err != nil {
			return errors.Wrap(err, "GetCommitment")
		}

		commitmentProving[i] = new(privacy.Point)
		commitmentProving[i], err = commitmentProving[i].FromBytesS(temp)
		if err != nil {
			return errors.Wrap(err, "GetCommitment")
		}
	}

	// prepare witness for proving
	witness := new(zkp.PaymentWitness)
	paymentWitnessParam := zkp.PaymentWitnessParam{
		HasPrivacy:              params.hasPrivacy,
		PrivateKey:              new(privacy.Scalar).FromBytesS(*params.senderSK),
		InputCoins:              params.inputCoins,
		OutputCoins:             outputCoins,
		PublicKeyLastByteSender: pkLastByteSender,
		Commitments:             commitmentProving,
		CommitmentIndices:       commitmentIndexs,
		MyCommitmentIndices:     myCommitmentIndexs,
		Fee:                     params.fee,
	}

	err = witness.Init(paymentWitnessParam)
	if err.(*privacy.PrivacyError) != nil {
		return errors.Wrap(err, "witness.Init")
	}

	tx.Proof, err = witness.Prove(params.hasPrivacy)
	if err.(*privacy.PrivacyError) != nil {
		return errors.Wrap(err, "witness.Prove")
	}

	// set private key for signing tx
	if params.hasPrivacy {
		randSK := witness.GetRandSecretKey()
		tx.sigPrivKey = append(*params.senderSK, randSK.ToBytesS()...)

		// encrypt coin details (Randomness)
		// hide information of output coins except coin commitments, public key, snDerivators
		for i := 0; i < len(tx.Proof.GetOutputCoins()); i++ {
			err = tx.Proof.GetOutputCoins()[i].Encrypt(params.paymentInfo[i].PaymentAddress.Tk)
			if err.(*privacy.PrivacyError) != nil {
				return errors.Wrap(err, "EncryptOutput")
			}
			tx.Proof.GetOutputCoins()[i].CoinDetails.SetSerialNumber(nil)
			tx.Proof.GetOutputCoins()[i].CoinDetails.SetValue(0)
			tx.Proof.GetOutputCoins()[i].CoinDetails.SetRandomness(nil)
		}

		// hide information of input coins except serial number of input coins
		for i := 0; i < len(tx.Proof.GetInputCoins()); i++ {
			tx.Proof.GetInputCoins()[i].CoinDetails.SetCoinCommitment(nil)
			tx.Proof.GetInputCoins()[i].CoinDetails.SetValue(0)
			tx.Proof.GetInputCoins()[i].CoinDetails.SetSNDerivator(nil)
			tx.Proof.GetInputCoins()[i].CoinDetails.SetPublicKey(nil)
			tx.Proof.GetInputCoins()[i].CoinDetails.SetRandomness(nil)
		}

	} else {
		tx.sigPrivKey = []byte{}
		randSK := big.NewInt(0)
		tx.sigPrivKey = append(*params.senderSK, randSK.Bytes()...)
	}

	// sign tx
	tx.PubKeyLastByteSender = common.GetShardIDFromLastByte(pkLastByteSender)
	err = tx.signTx()
	if err != nil {
		return errors.Wrap(err, "SignTx")
	}

	elapsedPrivacy := time.Since(startPrivacy)
	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("Creating payment proof time %s", elapsedPrivacy))
	fmt.Println(fmt.Sprintf("Successfully creating normal tx %+v in %s time", *tx.Hash(), elapsed))
	return nil
}

// signTx - signs tx
func (tx *Tx) signTx() error {
	//Check input transaction
	if tx.Sig != nil {
		return errors.New("input transaction must be an unsigned one")
	}

	/****** using Schnorr signature *******/
	// sign with sigPrivKey
	// prepare private key for Schnorr
	sk := new(privacy.Scalar).FromBytesS(tx.sigPrivKey[:common.BigIntSize])
	r := new(privacy.Scalar).FromBytesS(tx.sigPrivKey[common.BigIntSize:])
	sigKey := new(privacy.SchnorrPrivateKey)
	sigKey.Set(sk, r)

	// save public key for verification signature tx
	tx.SigPubKey = sigKey.GetPublicKey().GetPublicKey().ToBytesS()
	signature, err := sigKey.Sign(tx.Hash()[:])
	if err != nil {
		return err
	}

	// convert signature to byte array
	tx.Sig = signature.Bytes()

	return nil
}

func (tx *Tx) Hash() *common.Hash {
	if tx.cachedHash != nil {
		return tx.cachedHash
	}
	inBytes := []byte(tx.String())
	hash := common.HashH(inBytes)
	tx.cachedHash = &hash
	return &hash
}

func (tx Tx) String() string {
	record := strconv.Itoa(int(tx.Version))

	record += strconv.FormatInt(tx.LockTime, 10)
	record += strconv.FormatUint(tx.Fee, 10)
	if tx.Proof != nil {
		tmp := base64.StdEncoding.EncodeToString(tx.Proof.Bytes())
		//tmp := base58.Base58Check{}.Encode(tx.Proof.Bytes(), 0x00)
		record += tmp
		// fmt.Printf("Proof check base 58: %v\n",tmp)
	}
	if tx.Metadata != nil {
		metadataHash := tx.Metadata.Hash()
		//Logger.log.Debugf("\n\n\n\n test metadata after hashing: %v\n", metadataHash.GetBytes())
		metadataStr := metadataHash.String()
		record += metadataStr
	}

	//TODO: To be uncomment
	// record += string(tx.Info)
	return record
}

func (tx Tx) GetSenderAddrLastByte() byte {
	return tx.PubKeyLastByteSender
}
