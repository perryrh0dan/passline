package ctxutil

import (
	"context"

	"github.com/urfave/cli/v2"
)

type contextKey int

const (
	ctxKeyColor contextKey = iota
	ctxKeyTerminal
	ctxKeyInteractive
	ctxKeyStdin
	ctxKeyClipTimeout
	ctxKeyConfirm
	ctxKeyShowSafeContent
	ctxKeyAlwaysYes
	ctxKeyVerbose
	ctxKeyAutoClip
	ctxKeyNotifications
	ctxKeyForce
	ctxKeyDefaultUsername
	ctxKeyEmail
	ctxKeyAdvanced
	ctxKeyQuickSelect
	ctxKeyCategory
	ctxKeyPhoneNumber
  ctxNoClip
)

// WithGlobalFlags parses any global flags from the cli context and returns
// a regular context
func WithGlobalFlags(c *cli.Context) context.Context {
	ctx := c.Context
	if c.Bool("yes") {
		ctx = WithAlwaysYes(ctx, true)
	}
	if c.IsSet("category") {
		ctx = WithCategory(ctx, c.String("category"))
	}
  if c.Bool("noclip") {
    ctx = WithNoClip(ctx, true)
  }

	return ctx
}

// ProgressCallback is a callback for updateing progress
type ProgressCallback func()

// WithColor returns a context with an explicit value for color
func WithColor(ctx context.Context, color bool) context.Context {
	return context.WithValue(ctx, ctxKeyColor, color)
}

// HasColor returns true if a value for Color has been set in this context
func HasColor(ctx context.Context) bool {
	return hasBool(ctx, ctxKeyColor)
}

// IsColor returns the value of color or the default (true)
func IsColor(ctx context.Context) bool {
	return is(ctx, ctxKeyColor, true)
}

// WithTerminal returns a context with an explicit value for terminal
func WithTerminal(ctx context.Context, isTerm bool) context.Context {
	return context.WithValue(ctx, ctxKeyTerminal, isTerm)
}

// HasTerminal returns true if a value for Terminal has been set in this context
func HasTerminal(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKeyTerminal).(bool)
	return ok
}

// IsTerminal returns the value of terminal or the default (true)
func IsTerminal(ctx context.Context) bool {
	bv, ok := ctx.Value(ctxKeyTerminal).(bool)
	if !ok {
		return true
	}
	return bv
}

// WithInteractive returns a context with an explicit value for interactive
func WithInteractive(ctx context.Context, isInteractive bool) context.Context {
	return context.WithValue(ctx, ctxKeyInteractive, isInteractive)
}

// HasInteractive returns true if a value for Interactive has been set in this context
func HasInteractive(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKeyInteractive).(bool)
	return ok
}

// IsInteractive returns the value of interactive or the default (true)
func IsInteractive(ctx context.Context) bool {
	bv, ok := ctx.Value(ctxKeyInteractive).(bool)
	if !ok {
		return true
	}
	return bv
}

// WithStdin returns a context with the value for Stdin set. If true some input
// is available on Stdin (e.g. something is being piped into it)
func WithStdin(ctx context.Context, isStdin bool) context.Context {
	return context.WithValue(ctx, ctxKeyStdin, isStdin)
}

// HasStdin returns true if a value for Stdin has been set in this context
func HasStdin(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKeyStdin).(bool)
	return ok
}

// IsStdin returns the value of stdin, i.e. if it's true some data is being
// piped to stdin. If not set it returns the default value (false)
func IsStdin(ctx context.Context) bool {
	bv, ok := ctx.Value(ctxKeyStdin).(bool)
	if !ok {
		return false
	}
	return bv
}

// WithClipTimeout returns a context with the value for clip timeout set
func WithClipTimeout(ctx context.Context, to int) context.Context {
	return context.WithValue(ctx, ctxKeyClipTimeout, to)
}

// HasClipTimeout returns true if a value for ClipTimeout has been set in this context
func HasClipTimeout(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKeyClipTimeout).(int)
	return ok
}

// GetClipTimeout returns the value of clip timeout or the default (45)
func GetClipTimeout(ctx context.Context) int {
	iv, ok := ctx.Value(ctxKeyClipTimeout).(int)
	if !ok || iv < 1 {
		return 45
	}
	return iv
}

// WithConfirm returns a context with the value for ask for more set
func WithConfirm(ctx context.Context, bv bool) context.Context {
	return context.WithValue(ctx, ctxKeyConfirm, bv)
}

// HasConfirm returns true if a value for Confirm has been set in this context
func HasConfirm(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKeyConfirm).(bool)
	return ok
}

// IsConfirm returns the value of ask for more or the default (false)
func IsConfirm(ctx context.Context) bool {
	bv, ok := ctx.Value(ctxKeyConfirm).(bool)
	if !ok {
		return false
	}
	return bv
}

// WithShowSafeContent returns a context with the value for ShowSafeContent set
func WithShowSafeContent(ctx context.Context, bv bool) context.Context {
	return context.WithValue(ctx, ctxKeyShowSafeContent, bv)
}

// HasShowSafeContent returns true if a value for ShowSafeContent has been set in this context
func HasShowSafeContent(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKeyShowSafeContent).(bool)
	return ok
}

// IsShowSafeContent returns the value of ShowSafeContent or the default (false)
func IsShowSafeContent(ctx context.Context) bool {
	bv, ok := ctx.Value(ctxKeyShowSafeContent).(bool)
	if !ok {
		return false
	}
	return bv
}

// WithAlwaysYes returns a context with the value of always yes set
func WithAlwaysYes(ctx context.Context, bv bool) context.Context {
	return context.WithValue(ctx, ctxKeyAlwaysYes, bv)
}

// HasAlwaysYes returns true if a value for AlwaysYes has been set in this context
func HasAlwaysYes(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKeyAlwaysYes).(bool)
	return ok
}

// IsAlwaysYes returns the value of always yes or the default (false)
func IsAlwaysYes(ctx context.Context) bool {
	bv, ok := ctx.Value(ctxKeyAlwaysYes).(bool)
	if !ok {
		return false
	}
	return bv
}

// WithVerbose returns a context with the value for verbose set
func WithVerbose(ctx context.Context, verbose bool) context.Context {
	return context.WithValue(ctx, ctxKeyVerbose, verbose)
}

// HasVerbose returns true if a value for Verbose has been set in this context
func HasVerbose(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKeyVerbose).(bool)
	return ok
}

// IsVerbose returns the value of verbose or the default (false)
func IsVerbose(ctx context.Context) bool {
	return is(ctx, ctxKeyVerbose, false)
}

// WithNotifications returns a context with the value for Notifications set
func WithNotifications(ctx context.Context, notification bool) context.Context {
	return context.WithValue(ctx, ctxKeyNotifications, notification)
}

// HasNotifications returns true if a value for Notifications has been set in this context
func HasNotifications(ctx context.Context) bool {
	return hasBool(ctx, ctxKeyNotifications)
}

// IsNotifications returns the value of Notifications or the default (true)
func IsNotifications(ctx context.Context) bool {
	return is(ctx, ctxKeyNotifications, true)
}

// WithAdvanced returns a context with the value for Advanced set
func WithAdvanced(ctx context.Context, advanced bool) context.Context {
	return context.WithValue(ctx, ctxKeyAdvanced, advanced)
}

// HasAdvanced returns true if a value for Advanced has been set in this context
func HasAdvanced(ctx context.Context) bool {
	return hasBool(ctx, ctxKeyAdvanced)
}

// IsAdvanced returns the value of Advanced or the default (false)
func IsAdvanced(ctx context.Context) bool {
	return is(ctx, ctxKeyAdvanced, false)
}

// WithAutoClip returns a context with the value for AutoClip set
func WithAutoClip(ctx context.Context, bv bool) context.Context {
	return context.WithValue(ctx, ctxKeyAutoClip, bv)
}

// HasAutoClip returns true if a value for AutoClip has been set in this context
func HasAutoClip(ctx context.Context) bool {
	return hasBool(ctx, ctxKeyAutoClip)
}

// IsAutoClip returns the value of AutoClip or the default (true)
func IsAutoClip(ctx context.Context) bool {
	return is(ctx, ctxKeyAutoClip, true)
}

// WithForce returns a context with the force flag set
func WithForce(ctx context.Context, bv bool) context.Context {
	return context.WithValue(ctx, ctxKeyForce, bv)
}

// HasForce returns true if the context has the force flag set
func HasForce(ctx context.Context) bool {
	return hasBool(ctx, ctxKeyForce)
}

// IsForce returns the force flag value of the default (false)
func IsForce(ctx context.Context) bool {
	return is(ctx, ctxKeyForce, false)
}

// WithDefaultUsername returns a context with the username set in the context
func WithDefaultUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, ctxKeyDefaultUsername, username)
}

// HasDefaultUsername returns true if the context has defaultUsername set
func HasDefaultUsername(ctx context.Context) bool {
	return hasString(ctx, ctxKeyDefaultUsername)
}

// GetDefaultUsername returns the username from the context
func GetDefaultUsername(ctx context.Context) string {
	sv, ok := ctx.Value(ctxKeyDefaultUsername).(string)
	if !ok {
		return ""
	}
	return sv
}

// WithQuickSelect returns a context with the quickselect flag set
func WithQuickSelect(ctx context.Context, qs bool) context.Context {
	return context.WithValue(ctx, ctxKeyQuickSelect, qs)
}

// HasQuickSelect returns true if the context has the quickselect flag set
func HasQuickSelect(ctx context.Context) bool {
	return hasBool(ctx, ctxKeyQuickSelect)
}

// IsQuickSelect returns the quickselect flag value of the default (false)
func IsQuickSelect(ctx context.Context) bool {
	return is(ctx, ctxKeyQuickSelect, false)
}

// WithCategory returns a context with the category set in the context
func WithCategory(ctx context.Context, category string) context.Context {
	return context.WithValue(ctx, ctxKeyCategory, category)
}

// HasCategory returns true if the context has the category set
func HasCategory(ctx context.Context) bool {
	return hasString(ctx, ctxKeyCategory)
}

// GetCategory returns the category from the context
func GetCategory(ctx context.Context) string {
	sv, ok := ctx.Value(ctxKeyCategory).(string)
	if !ok {
		return "*"
	}
	return sv
}

// WithCategory returns a context with the category set in the context
func WithPhoneNumber(ctx context.Context, phoneNumber string) context.Context {
	return context.WithValue(ctx, ctxKeyPhoneNumber, phoneNumber)
}

// HasPhoneNumber returns true if the context has the phone number set
func HasPhoneNumber(ctx context.Context) bool {
	return hasString(ctx, ctxKeyPhoneNumber)
}

// GetPhoneNumber returns the phone number from the context
func GetPhoneNumber(ctx context.Context) string {
	sv, ok := ctx.Value(ctxKeyPhoneNumber).(string)
	if !ok {
		return ""
	}
	return sv
}

// WithAlwaysYes returns a context with the value of always yes set
func WithNoClip(ctx context.Context, bv bool) context.Context {
	return context.WithValue(ctx, ctxNoClip, bv)
}

// IsAlwaysYes returns the value of always yes or the default (false)
func IsNoClip(ctx context.Context) bool {
	bv, ok := ctx.Value(ctxNoClip).(bool)
	if !ok {
		return false
	}
	return bv
}
