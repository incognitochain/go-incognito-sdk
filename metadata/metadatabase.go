package metadata

import (
	"strconv"

	"github.com/incognitochain/go-incognito-sdk/common"
)

type MetadataBase struct {
	Type int
}

func NewMetadataBase(thisType int) *MetadataBase {
	return &MetadataBase{Type: thisType}
}

func (mb MetadataBase) GetType() int {
	return mb.Type
}

func (mb MetadataBase) Hash() *common.Hash {
	record := strconv.Itoa(mb.Type)
	data := []byte(record)
	hash := common.HashH(data)
	return &hash
}

func (mb *MetadataBase) CalculateSize() uint64 {
	return 0
}