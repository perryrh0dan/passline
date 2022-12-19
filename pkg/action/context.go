package action

import "context"

type contextKey int

const (
	ctxKeyClip contextKey = iota
)

// WithClip returns a context with the value for clip (for copy to clipboard)
// set
func WithClip(ctx context.Context, clip bool) context.Context {
	return context.WithValue(ctx, ctxKeyClip, clip)
}

// IsClip returns the value of clip or the default (false)
func IsClip(ctx context.Context) bool {
	bv, ok := ctx.Value(ctxKeyClip).(bool)
	if !ok {
		return false
	}
	return bv
}
