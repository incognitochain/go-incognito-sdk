package privacy_util

type PrivacyUtilLogger struct {
	//Log common.Logger
}

func (logger *PrivacyUtilLogger) Init() {
	//logger.Log = inst
}

// Global instant to use
var Logger = PrivacyUtilLogger{}
