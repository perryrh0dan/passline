package config

import (
	"context"
	"passline/pkg/ctxutil"
)

// WithContext returns a context with all config options set for this store
// config, iff they have not been already set in the context
func (c *Config) WithContext(ctx context.Context) context.Context {
	if !ctxutil.HasAutoClip(ctx) {
		ctx = ctxutil.WithAutoClip(ctx, c.AutoClip)
	}

	return ctx
}
