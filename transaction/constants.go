package transaction

const (
	// txVersion is the current latest supported transaction version.
	txVersion                        = 1
	ValidateTimeForOneoutOfManyProof = 1574985600 // GMT: Friday, November 29, 2019 12:00:00 AM
)

const (
	CustomTokenInit = iota
	CustomTokenTransfer
	CustomTokenCrossShard
)

const (
	NormalCoinType = iota
	CustomTokenPrivacyType
)

const MaxSizeInfo = 512
