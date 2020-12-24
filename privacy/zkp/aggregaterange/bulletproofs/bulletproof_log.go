package bulletproofs

type BulletproofsLogger struct {
	//Log common.Logger
}

func (logger *BulletproofsLogger) Init() {
	//logger.Log = inst
}

// Global instant to use
var Logger = BulletproofsLogger{}
