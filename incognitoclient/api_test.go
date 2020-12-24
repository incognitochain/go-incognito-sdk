package incognitoclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IncognitoTestSuite struct {
	suite.Suite
	client *PublicIncognito
	wallet *Wallet
	stake  *Stake
	block  *BlockInfo
	pdex   *PDex
}

func (t *IncognitoTestSuite) SetupTest() {
	client := &http.Client{}
	publicIncognito := NewPublicIncognito(client, "https://testnet.incognito.org/fullnode")
	t.client = publicIncognito

	blockInfo := NewBlockInfo(publicIncognito)
	t.wallet = NewWallet(publicIncognito, blockInfo)
	t.stake = NewStake(publicIncognito)
	t.block = NewBlockInfo(publicIncognito)
	t.pdex = NewPDex(publicIncognito, blockInfo)

}

func TestIncognitoTestSuite(t *testing.T) {
	suite.Run(t, new(IncognitoTestSuite))
}
