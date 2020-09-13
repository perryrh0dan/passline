package config

import (
	"context"
	"passline/pkg/ctxutil"
)

// WithContext returns a context with all config options set for this store
// config, if they have not been already set in the context
func (c *Config) WithContext(ctx context.Context) context.Context {
	if !ctxutil.HasAutoClip(ctx) {
		ctx = ctxutil.WithAutoClip(ctx, c.AutoClip)
	}
	if !ctxutil.HasNotifications(ctx) {
		ctx = ctxutil.WithNotifications(ctx, c.Notifications)
	}
	if !ctxutil.HasDefaultUsername(ctx) {
		ctx = ctxutil.WithDefaultUsername(ctx, c.DefaultUsername)
	}

	return ctx
}
