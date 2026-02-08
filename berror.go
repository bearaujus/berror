package berror

import (
	"errors"
	"fmt"
)

// New creates a WrappedErr using error, capturing the caller stack trace.
func New(err error) WrappedErr {
	if err == nil {
		return nil
	}
	return NewErrDefinition("%w", OptionErrDefinitionWithCustomStackTraceCapturer(func() string { return captureStackTrace(5) })).New(err)
}

// Newf creates a WrappedErr using format and args, capturing the caller stack trace.
func Newf(format string, args ...any) WrappedErr {
	return NewErrDefinition(format, OptionErrDefinitionWithCustomStackTraceCapturer(func() string { return captureStackTrace(5) })).New(args...)
}

// Wrap wraps an existing error with an additional context message.
// Returns nil if err is nil.
func Wrap(err error, message string) WrappedErr {
	if err == nil {
		return nil
	}
	return NewErrDefinition(message+": %w", OptionErrDefinitionWithCustomStackTraceCapturer(func() string { return captureStackTrace(5) })).New(err)
}

// Wrapf wraps an existing error with a formatted context message.
// Returns nil if err is nil.
func Wrapf(err error, format string, args ...any) WrappedErr {
	if err == nil {
		return nil
	}
	args = append(args, err)
	return NewErrDefinition(format+": %w", OptionErrDefinitionWithCustomStackTraceCapturer(func() string { return captureStackTrace(5) })).New(args...)
}

// IsCode checks if the given error (or any error in its chain) has the specified error code.
func IsCode(err error, code string) bool {
	if err == nil || code == "" {
		return false
	}
	we, ok := CastToWrappedErrFromErr(err)
	if ok && we.Code() == code {
		return true
	}
	// Check wrapped errors
	var unwrapper interface{ Unwrap() []error }
	if errors.As(err, &unwrapper) {
		for _, e := range unwrapper.Unwrap() {
			if IsCode(e, code) {
				return true
			}
		}
	}
	// Check single unwrap
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		return IsCode(unwrapped, code)
	}
	return false
}

// GetCode returns the error code from the error if it's a WrappedErr.
// Returns empty string if err is nil or not a WrappedErr.
func GetCode(err error) string {
	if err == nil {
		return ""
	}
	we, ok := CastToWrappedErrFromErr(err)
	if !ok {
		return ""
	}
	return we.Code()
}

// GetStackTrace returns the stack trace from the error if it's a WrappedErr.
// Returns empty string if err is nil or not a WrappedErr.
func GetStackTrace(err error) string {
	if err == nil {
		return ""
	}
	we, ok := CastToWrappedErrFromErr(err)
	if !ok {
		return ""
	}
	return we.StackTrace()
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
		format:    format,
		formatter: ErrWrapperFormatterDefault,
	}
	for _, opt := range opts {
		opt(&ed)
	}
	if ed.stackTraceCapturer == nil {
		ed.stackTraceCapturer = func() string { return captureStackTrace(4) }
	}
	return &ed
}

func (ed *errDefinition) New(a ...any) WrappedErr {
	if ed == nil {
		return nil
	}
	we := &wrappedErr{
		ew:   ed,
		Args: a,
		err:  fmt.Errorf(ed.format, a...).Error(),
	}
	if !ed.disableStackTrace {
		we.stack = ed.stackTraceCapturer()
	}
	return we
}

func (ed *errDefinition) Is(err error) bool {
	if ed == nil {
		return false
	}
	return errors.Is(err, ed.New())
}

func (ed *errDefinition) Format() string {
	if ed == nil {
		return ""
	}
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
	if we == nil {
		return ""
	}
	return we.ew.code
}

func (we *wrappedErr) StackTrace() string {
	if we == nil {
		return ""
	}
	return we.stack
}

func (we *wrappedErr) Error() string {
	if we == nil {
		return fmt.Sprint(nil)
	}
	return we.ew.formatter(we.err, we.ew.code, we.stack)
}

func (we *wrappedErr) RawError() string {
	if we == nil {
		return fmt.Sprint(nil)
	}
	return we.err
}

func (we *wrappedErr) Is(err error) bool {
	if we == nil && err == nil {
		return true
	}
	if we == nil {
		return false
	}
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
	if we == nil {
		return errs
	}
	for _, arg := range we.Args {
		switch errType := arg.(type) {
		case error:
			errs = append(errs, errType)
		}
	}
	return errs
}

func (we *wrappedErr) String() string {
	if we == nil {
		return fmt.Sprint(nil)
	}
	return we.Error()
}

func (we *wrappedErr) ErrorDefinition() ErrDefinition {
	if we == nil {
		return nil
	}
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

// IsWrappedErr checks if the given error is a WrappedErr.
func IsWrappedErr(err error) bool {
	_, ok := CastToWrappedErrFromErr(err)
	return ok
}

// Join combines multiple errors into a single WrappedErr.
// Returns nil if all errors are nil.
func Join(errs ...error) WrappedErr {
	var nonNilErrs []any
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}
	if len(nonNilErrs) == 0 {
		return nil
	}
	if len(nonNilErrs) == 1 {
		if we, ok := nonNilErrs[0].(WrappedErr); ok {
			return we
		}
		return New(nonNilErrs[0].(error))
	}
	format := make([]string, len(nonNilErrs))
	for i := range nonNilErrs {
		format[i] = "%w"
	}
	return NewErrDefinition(
		fmt.Sprintf("multiple errors: [%s]", joinStrings(format, ", ")),
		OptionErrDefinitionWithCustomStackTraceCapturer(func() string { return captureStackTrace(5) }),
	).New(nonNilErrs...)
}

// joinStrings joins a slice of strings with a separator.
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// WithCode creates a new WrappedErr from an existing error with an error code.
// Returns nil if err is nil.
func WithCode(err error, code string) WrappedErr {
	if err == nil {
		return nil
	}
	return NewErrDefinition("%w",
		OptionErrDefinitionWithErrCode(code),
		OptionErrDefinitionWithCustomStackTraceCapturer(func() string { return captureStackTrace(5) }),
	).New(err)
}

// Must panics if err is not nil. Useful for initialization code.
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// Ignore discards the error and returns only the value.
// Use with caution - only when you're certain the error can be safely ignored.
func Ignore[T any](val T, _ error) T {
	return val
}
