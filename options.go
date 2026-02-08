package berror

// ErrDefinitionOption defines a function type for customizing error definitions.
type ErrDefinitionOption func(*errDefinition)

// OptionErrDefinitionWithErrCode adds a custom error code to the error definition.
func OptionErrDefinitionWithErrCode(code string) ErrDefinitionOption {
	return func(ed *errDefinition) {
		if code == "" {
			return
		}
		ed.code = code
	}
}

// OptionErrDefinitionWithCustomFormater adds a custom formatter to the error definition.
func OptionErrDefinitionWithCustomFormater(formatter ErrWrapperFormatter) ErrDefinitionOption {
	return func(ed *errDefinition) {
		if formatter == nil {
			return
		}
		ed.formatter = formatter
	}
}

// OptionErrDefinitionWithDisabledStackTrace disables stack trace capture for the error definition.
func OptionErrDefinitionWithDisabledStackTrace() ErrDefinitionOption {
	return func(ed *errDefinition) {
		ed.disableStackTrace = true
	}
}

// OptionErrDefinitionWithCustomStackTraceCapturer adds a custom stack trace capturer to the error definition.
func OptionErrDefinitionWithCustomStackTraceCapturer(capturer ErrWrapperStackTraceCapturer) ErrDefinitionOption {
	return func(ed *errDefinition) {
		if capturer == nil {
			return
		}
		ed.stackTraceCapturer = capturer
	}
}
