package entity

import "encoding/json"

type WalletSend struct {
	Type             int // 0: constant, 1: token
	TokenID          string
	TokenName        string
	TokenSymbol      string
	PaymentAddresses map[string]uint64
	TokenAmount      uint64
	TokenFee         uint64
}

type ReportPdex struct {
	Day         string  `json:"Day"`
	Total       int     `json:"Total"`
	Month       int     `json:"Month"`
	Year        int     `json:"Year"`
	TotalVolume float64 `json:"TotalVolume"`
}

type RewardAmount struct {
	PublicKey string
	Reward    float64
}

type GetShardFromPaymentAddressObject struct {
	ID     int `json:"Id"`
	Result struct {
		PublicKeyInBase58Check string `json:"PublicKeyInBase58Check"`
		PublicKeyInBytes       []int  `json:"PublicKeyInBytes"`
		PublicKeyInHex         string `json:"PublicKeyInHex"`
	} `json:"Result"`
	Error   interface{} `json:"Error"`
	Params  []string    `json:"Params"`
	Method  string      `json:"Method"`
	Jsonrpc string      `json:"Jsonrpc"`
}

type StabilityInfo struct {
	Result struct {
		GOVConstitution struct {
			ConstitutionIndex int
			GOVParams         struct {
				OracleNetwork struct {
					OraclePubKeys []string
				}
			}
		}
		DCBConstitution struct {
			ConstitutionIndex int
		}
		GOVGovernor struct {
			BoardPaymentAddress []struct {
				Pk string
				Tk string
			}
		}
		Oracle struct {
			DCBToken uint64
			GOVToken uint64
			Constant uint64
			ETH      uint64
			BTC      uint64
		}
	}
}

type ResponsePdex struct {
	Result struct {
		Data []*ReportPdex `json:"Data"`
	} `json:"Result"`
}

type ListCustomTokenBalance struct {
	PaymentAddress         string               `json:"PaymentAddress"`
	ListCustomTokenBalance []CustomTokenBalance `json:"ListCustomTokenBalance"`
}

type CustomTokenBalance struct {
	Name    string `json:"Name"`
	Symbol  string `json:"Symbol"`
	Amount  uint64 `json:"Amount"`
	TokenID string `json:"TokenID"`
}

type ProofDetail struct {
	InputCoins  []*CoinDetail
	OutputCoins []*CoinDetail
}

type PrivacyCustomTokenProofDetail struct {
	InputCoins  []*CoinDetail
	OutputCoins []*CoinDetail
}

type PrivacyCustomTokenData struct {
	PropertyID string `json:"PropertyID"`
}

type CoinDetail struct {
	CoinDetails          *Coin
	CoinDetailsEncrypted string
}

type Coin struct {
	PublicKey      string
	CoinCommitment string
	SerialNumber   string
	Value          uint64
	Info           string
}

type TransactionDetail struct {
	BlockHash string `json:"BlockHash"`
	Index     uint64 `json:"Index"`
	ChainId   byte   `json:"ChainId"`
	Hash      string `json:"Hash"`
	Version   int8   `json:"Version"`
	Type      string `json:"Type"` // Transaction type
	LockTime  string `json:"LockTime"`
	Fee       uint64 `json:"Fee"` // Fee applies: always consant

	Descs    []interface{} `json:"Descs"`
	JSPubKey []byte        `json:"JSPubKey,omitempty"` // 64 bytes
	JSSig    []byte        `json:"JSSig,omitempty"`    // 64 bytes

	AddressLastByte byte   `json:"AddressLastByte"`
	Metadata        string `json:"Metadata"`
	IsInMempool     bool   `json:"IsInMempool"`
	IsInBlock       bool   `json:"IsInBlock"`
	Info            string `json:"Info"`
	ShardID         int8   `json:"ShardID"`

	ProofDetail                   *ProofDetail                   `json:"ProofDetail"`
	PrivacyCustomTokenData        string                         `json:"PrivacyCustomTokenData"`
	PrivacyCustomTokenProofDetail *PrivacyCustomTokenProofDetail `json:"PrivacyCustomTokenProofDetail"`
}

func (self TransactionDetail) String() string {
	str, _ := json.MarshalIndent(self, "", "\t")
	return string(str)
}

type GetBlockChainInfoResult struct {
	ChainName  string                      `json:"ChainName"`
	BestBlocks map[string]GetBestBlockItem `json:"BestBlocks"`
}
type GetBlockInfo struct {
	Txs []TxInfo `json:"Txs"`
}

type TxInfo struct {
	Hash string `json:"Hash"`
}

type GetBestBlockItem struct {
	Height           int32  `json:"Height"`
	Hash             string `json:"Hash"`
	TotalTxs         uint64 `json:"TotalTxs"`
	SalaryFund       uint64 `json:"SalaryFund"`
	BasicSalary      uint64 `json:"BasicSalary"`
	SalaryPerTx      uint64 `json:"SalaryPerTx"`
	BlockProducer    string `json:"BlockProducer"`
	BlockProducerSig string `json:"BlockProducerSig"`
	Time             int64  `json:"Time"`
}

type PCustomToken struct {
	Amount    uint64
	ID        string
	Image     string
	IsPrivacy bool
	Name      string
	Symbol    string
}

type ReceivedTransactions struct {
	ReceivedTransactions []ReceivedTransaction `json:"ReceivedTransactions"`
}

type ReceivedTransaction struct {
	Hash            string                `json:"Hash"`
	Info            string                `json:"Info"`
	ReceivedAmounts map[string]CoinDetail `json:"ReceivedAmounts"`
}

//
type AccountAddressResp struct {
	Result struct {
		PrivateKey     string `json:"PrivateKey"`
		PaymentAddress string `json:"PaymentAddress"`
		ReadonlyKey    string `json:"ReadonlyKey"`
		Pubkey         string `json:"Pubkey"`
		ValidatorKey   string `json:"ValidatorKey"`
	} `json:"Result"`
	Error *string `json:"Error"`
	ID    int     `json:"Id"`
}

//
type CommitteeKeyString struct {
	IncPubKey    string
	MiningPubKey map[string]string
}

type CommitteeKeySetAutoStake struct {
	IncPubKey    string
	MiningPubKey map[string]string
	IsAutoStake  bool
}

type BeaconBestStateResp struct {
	Result struct {
		RewardReceiver map[string]string             `json:"RewardReceiver"`
		ShardCommittee map[byte][]CommitteeKeyString `json:"ShardCommittee"`
		AutoStaking    []CommitteeKeySetAutoStake    `json:"AutoStaking"`
	} `json:"Result"`
	Error *string `json:"Error"`
	ID    int     `json:"Id"`
}

type DecrypTransactionPRV struct {
	AmountPRV uint64 `json:"0000000000000000000000000000000000000000000000000000000000000004"`
	IsMempool bool   `json:"IsMempool"`
}

type BurningForDepositToSCRes struct {
	TxID string `json:"TxID"`
}

type Unstake struct {
	IncPubKey   string
	IsAutoStake bool
}

type RewardData struct {
	PublicKey   string
	RewardItems []RewardDataItem
}

type RewardDataItem struct {
	TokenId string
	Reward  uint64
}

type RewardItems struct {
	TokenId string
	Reward  float64
}

type Utxo struct {
	Value        uint64
	SnDerivator  string
	SerialNumber string
}

type TotalStaker struct {
	TotalStaker  uint64
}