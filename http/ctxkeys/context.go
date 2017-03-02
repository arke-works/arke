package ctxkeys

type ctxKey int

const (
	CtxLoggerKey ctxKey = iota
	CtxPivotIDKey
	CtxSizeKey
	CtxFountainKey
)
