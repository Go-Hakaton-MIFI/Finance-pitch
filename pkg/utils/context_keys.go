package utils

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	ContextKeyRequestId = contextKey("RequestID")
	ContextKeyUser      = contextKey("User")
)

func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	requestID, ok := ctx.Value(ContextKeyRequestId).(string)
	return requestID, ok
}
