package constant

import "errors"

type PDexTradeStatus int

const (
	PDexTradePending  PDexTradeStatus = iota
	PDexTradeSuccess                  // 1
	PDexTradeRefunded                 // 2
)

var (
	ErrTxHashNotExists          = errors.New("tx hash does not exist")
	ErrTxHashInvalidFromAddress = errors.New("tx hash invalid: from address does not matched!")
	ErrTxHashInvalidToAddress   = errors.New("tx hash invalid: to address does not matched!")
)

const (
	PDEX_TRADE_STEPS = 2
	EstimateFee      = 5 //30000000
)
