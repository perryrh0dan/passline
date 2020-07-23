// +build windows

package notify

import (
	"context"
	"os"
	"os/exec"

	"passline/pkg/ctxutil"

	"golang.org/x/sys/windows/registry"
	"gopkg.in/toast.v1"
)

var isWindows10 bool

func init() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer k.Close()

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return
	}

	isWindows10 = maj == 10
}

// Notify displays a desktop notification through msg
func Notify(ctx context.Context, subj, msg string) error {
	if os.Getenv("PASSLINE_NO_NOTIFY") != "" || !ctxutil.IsNotifications(ctx) {
		return nil
	}

	if isWindows10 {
		return sendBaloon(subj, msg)
	} else {
		return sendToast(subj, msg)
	}
}

func sendToast(subj, msg string) error {
	notification := toast.Notification{
		AppID:   "Passline",
		Title:   subj,
		Message: msg,
	}
	return notification.Push()
}

func sendBaloon(subj, msg string) error {
	winmsg, err := exec.LookPath("msg")
	if err != nil {
		return err
	}

	return exec.Command(winmsg,
		"*",
		"/TIME:3",
		subj+"\n\n"+msg,
	).Start()
}
