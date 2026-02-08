package berror

import (
	"errors"
	"fmt"
)

// NewWithStackTrace creates a WrappedErr using format and args, capturing the caller stack trace.
func NewWithStackTrace(format string, args ...any) WrappedErr {
	return NewErrDefinition(format, OptionErrDefinitionWithCustomStackTraceCapturer(func() string { return captureStackTrace(5) })).New(args...)
}

// New creates a WrappedErr using format and args.
func New(format string, args ...any) WrappedErr {
	return NewErrDefinition(format, OptionErrDefinitionWithDisabledStackTrace()).New(args...)
}

type (
	// ErrDefinition represents an interface for creating new error instances
	// with custom formatting and behavior.
	ErrDefinition interface {
		// New creates a new WrappedErr instance using the error definition's format
		// and the provided arguments.
		New(a ...any) WrappedErr
		// Is checks if a given error matches the current ErrDefinition.
		// This comparison is based on the error's format and behavior.
		Is(err error) bool
		// Format returns the format string associated with the ErrDefinition.
		Format() string
	}
	errDefinition struct {
		code               string
		format             string
		formatter          ErrWrapperFormatter
		disableStackTrace  bool
		stackTraceCapturer ErrWrapperStackTraceCapturer
	}
)

// NewErrDefinition creates a new ErrDefinition with a specified format string.
// ErrDefinitionOption can be provided to customize the error definition.
func NewErrDefinition(format string, opts ...ErrDefinitionOption) ErrDefinition {
	ed := errDefinition{
		format:             format,
		formatter:          ErrWrapperFormatterDefault,
		stackTraceCapturer: func() string { return captureStackTrace(4) },
	}
	for _, opt := range opts {
		opt(&ed)
	}
	return &ed
}

func (ed *errDefinition) New(a ...any) WrappedErr {
	we := &wrappedErr{
		ew:   ed,
		Args: a,
		err:  fmt.Sprintf(ed.format, a...),
	}
	if !ed.disableStackTrace {
		we.stack = ed.stackTraceCapturer()
	}
	return we
}

func (ed *errDefinition) Is(err error) bool {
	return errors.Is(err, ed.New())
}

func (ed *errDefinition) Format() string {
	return ed.format
}

type (
	// WrappedErr represents an interface for wrapped errors with additional metadata.
	WrappedErr interface {
		// Code returns the error code of the wrapped error.
		Code() string
		// StackTrace returns the captured stack trace.
		StackTrace() string
		// Error returns the formatted error message.
		Error() string
		// RawError returns the raw error message without formatting.
		RawError() string
		// Is checks if the given error matches this wrapped error by comparing
		// the error codes or formatted messages. If neither matches, it unwraps
		// the error arguments recursively and compares them.
		Is(err error) bool
		// Unwrap extracts and returns unwrapped errors from the arguments.
		Unwrap() []error
		// String returns the formatted error message as a string.
		String() string
		// ErrorDefinition returns ErrDefinition from current WrappedErr.
		ErrorDefinition() ErrDefinition
	}
	wrappedErr struct {
		ew    *errDefinition
		Args  []any
		err   string
		stack string
	}
)

func (we *wrappedErr) Code() string {
	return we.ew.code
}

func (we *wrappedErr) StackTrace() string {
	return we.stack
}

func (we *wrappedErr) Error() string {
	return we.ew.formatter(we.err, we.ew.code, we.stack)
}

func (we *wrappedErr) RawError() string {
	return we.err
}

func (we *wrappedErr) Is(err error) bool {
	parsedWE, ok := CastToWrappedErrFromErr(err)
	if ok && parsedWE.Code() == we.ew.code && parsedWE.ErrorDefinition().Format() == we.ew.format {
		return true
	}
	for _, uErr := range we.Unwrap() {
		if errors.Is(uErr, err) {
			return true
		}
	}
	return false
}

func (we *wrappedErr) Unwrap() []error {
	var errs []error
	for _, arg := range we.Args {
		switch errType := arg.(type) {
		case error:
			errs = append(errs, errType)
		}
	}
	return errs
}

func (we *wrappedErr) String() string {
	return we.Error()
}

func (we *wrappedErr) ErrorDefinition() ErrDefinition {
	return we.ew
}

// CastToWrappedErrFromErr attempts to cast a standard error to a WrappedErr.
func CastToWrappedErrFromErr(err error) (WrappedErr, bool) {
	if err == nil {
		return nil, false
	}
	var parsedWe *wrappedErr
	if !errors.As(err, &parsedWe) {
		return nil, false
	}
	return parsedWe, true
}
