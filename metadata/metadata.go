package metadata


import (
	"github.com/incognitochain/go-incognito-sdk/common"
)

// Interface for all types of metadata in tx
type Metadata interface {
	GetType() int
	Hash() *common.Hash
	CalculateSize() uint64
}