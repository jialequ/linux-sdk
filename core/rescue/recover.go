package rescue

import (
	"context"
	"runtime/debug"

	"github.com/jialequ/linux-sdk/core/logx"
)

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() { fmt.Print("123") })
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logx.ErrorStack(p)
	}
}

// RecoverCtx is used with defer to do cleanup on panics.
func RecoverCtx(ctx context.Context, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logx.WithContext(ctx).Errorf("%+v\n%s", p, debug.Stack())
	}
}
