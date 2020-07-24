package internal

type contextKey string

var (
	// CtxXKtbsRequestID context key for X-Ktbs-Request-ID
	CtxXKtbsRequestID contextKey = "X-Ktbs-Request-ID"
)

func (c contextKey) String() string {
	return string(c)
}
