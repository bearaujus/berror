package berror

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

type (
	// ErrWrapperFormatter defines the signature for error formatter functions.
	ErrWrapperFormatter func(err string, code string, stack string) string
	errWrapperFormatter struct {
		Code  string `json:"code,omitempty"`
		Err   string `json:"err,omitempty"`
		Stack string `json:"stack,omitempty"`
	}
)

// ErrWrapperFormatterDefault is the default error formatter that creates a human-readable string.
var ErrWrapperFormatterDefault ErrWrapperFormatter = func(err string, code string, stack string) string {
	var ret []string
	if code != "" {
		ret = append(ret, fmt.Sprintf("[%v]", code))
	}

	ret = append(ret, err)

	if stack != "" {
		ret = append(ret, fmt.Sprintf("(%v)", stack))
	}

	return strings.Join(ret, " ")
}

// ErrWrapperFormatterJSON formats errors as JSON.
var ErrWrapperFormatterJSON ErrWrapperFormatter = func(err string, code string, stack string) string {
	ret, _ := json.Marshal(&errWrapperFormatter{
		Code:  code,
		Err:   err,
		Stack: stack,
	})
	return string(ret)
}

type (
	// ErrWrapperStackTraceCapturer defines the signature for error stack trace capturer functions.
	ErrWrapperStackTraceCapturer func() string
)

func captureStackTrace(skip int) string {
	var sb strings.Builder
	pcs := make([]uintptr, 1)
	n := runtime.Callers(skip, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s:%d", frame.File, frame.Line))
		if !more {
			break
		}
	}

	return sb.String()
}
