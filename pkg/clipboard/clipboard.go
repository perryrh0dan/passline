package clipboard

import (
	"context"
	"fmt"
	"os"

	"passline/pkg/ctxutil"
	"passline/pkg/notify"

	// "github.com/gopasspw/gopass/internal/out"
	// "github.com/perryrh0dan/passline/pkg/out"

	"github.com/atotto/clipboard"
)

var (
	// ErrNotSupported is returned when the clipboard is not accessible
	ErrNotSupported = fmt.Errorf("WARNING: No clipboard available. Install xsel or xclip or use -f to print to console")
)

// CopyTo copies the given data to the clipboard and enqueues automatic
// clearing of the clipboard
func CopyTo(ctx context.Context, name string, content []byte) error {
	if clipboard.Unsupported {
		// out.Yellow(ctx, "%s", ErrNotSupported)
		_ = notify.Notify(ctx, "passline - clipboard", fmt.Sprintf("%s", ErrNotSupported))
		return nil
	}

	if err := clipboard.WriteAll(string(content)); err != nil {
		_ = notify.Notify(ctx, "passline - clipboard", "failed to write to clipboard")
		// return errors.Wrapf(err, "failed to write to clipboard")
	}

	if err := clear(ctx, content, 45); err != nil {
		_ = notify.Notify(ctx, "passline - clipboard", "failed to clear clipboard")
		// return errors.Wrapf(err, "failed to clear clipboard")
	}

	// out.Print(ctx, "✔ Copied %s to clipboard. Will clear in %d seconds.", color.YellowString(name), ctxutil.GetClipTimeout(ctx))
	_ = notify.Notify(ctx, "passline - clipboard", fmt.Sprintf("✔ Copied %s to clipboard. Will clear in %d seconds.", name, ctxutil.GetClipTimeout(ctx)))
	return nil
}

func killProc(pid int) {
	// err should be always nil, but just to be sure
	proc, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	// we ignore this error as we're going to return nil anyway
	_ = proc.Kill()
}
