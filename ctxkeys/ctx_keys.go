package ctxkeys

type contextKey string
type logger string

var (
	// CtxXKtbsRequestID context key for X-Ktbs-Request-ID
	CtxXKtbsRequestID contextKey = "X-Ktbs-Request-ID"

	// CtxLogger context key for logger
	CtxLogger logger = "Ktbs-Logger"
)

func (c contextKey) String() string {
	return string(c)
}

func (c logger) String() string {
	return string(c)
}
