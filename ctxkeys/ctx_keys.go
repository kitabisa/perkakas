package ctxkeys

type contextKey string
type logger string

var (
	// CtxXKtbsRequestID context key for X-Ktbs-Request-ID
	CtxXKtbsRequestID contextKey = "X-Ktbs-Request-ID"

	// CtxKtbsDonationIdentifier context key for X-Ktbs-Donation-Identifier
	CtxKtbsDonationIdentifier contextKey = "X-Ktbs-Donation-Identifier"

	// CtxLogger context key for logger
	CtxLogger logger = "Ktbs-Logger"

	CtxWatermillProcessID contextKey = "Ktbs-watermill-process-id"
)

func (c contextKey) String() string {
	return string(c)
}

func (c logger) String() string {
	return string(c)
}
