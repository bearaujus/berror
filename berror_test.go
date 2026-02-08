package berror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bearaujus/berror"
)

const (
	testFormat                      = "error timeout"
	testFormatWithArgs              = "error timeout: %v %v"
	testStackTraceCaputurerMockPath = "./test"
	testErrorCode                   = "t001"
)

var (
	testErrDef = berror.NewErrDefinition("xxx %v",
		berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock))
	testSampleErr                                              = errors.New("test sample error")
	testSampleErr2                                             = errors.New("test sample error 2")
	testSampleErr3                                             = errors.New("test sample error 3")
	stackTraceCapturerMock berror.ErrWrapperStackTraceCapturer = func() string {
		return testStackTraceCaputurerMockPath
	}
)

func TestNew(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "from nil error",
			args: args{
				err: nil,
			},
			want: nil,
		},
		{
			name: "from valid error",
			args: args{
				err: testSampleErr,
			},
			want: testSampleErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := berror.New(tt.args.err)
			if got == nil {
				if tt.want != nil {
					t.Fatal("got is nil but want is not nil")
				}
				return
			}
			if !got.Is(tt.want) {
				t.Fatal("got is not equals to want")
			}
		})
	}
}

func TestNewf(t *testing.T) {
	t.Run("simple format", func(t *testing.T) {
		err := berror.Newf("test error")
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.RawError() != "test error" {
			t.Fatalf("expected 'test error', got: %v", err.RawError())
		}
	})

	t.Run("format with args", func(t *testing.T) {
		err := berror.Newf("error: %v %v", "arg1", "arg2")
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.RawError() != "error: arg1 arg2" {
			t.Fatalf("expected 'error: arg1 arg2', got: %v", err.RawError())
		}
	})

	t.Run("format with error arg", func(t *testing.T) {
		err := berror.Newf("wrapped: %w", testSampleErr)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if !err.Is(testSampleErr) {
			t.Fatal("expected error to wrap testSampleErr")
		}
	})
}

func TestErrDefinition_Format(t *testing.T) {
	format := "test format %v"
	ed := berror.NewErrDefinition(format)
	if ed.Format() != format {
		t.Fatalf("expected format '%v', got: %v", format, ed.Format())
	}
}

func TestWrappedErr_Is_SelfComparison(t *testing.T) {
	err := berror.Newf("test error")
	if !errors.Is(err, err) {
		t.Fatal("expected error to be equal to itself")
	}
}

func TestCastToWrappedErrFromErr(t *testing.T) {
	t.Run("with nil error", func(t *testing.T) {
		we, ok := berror.CastToWrappedErrFromErr(nil)
		if ok || we != nil {
			t.Fatal("expected nil and false for nil error")
		}
	})

	t.Run("with non-wrapped error", func(t *testing.T) {
		we, ok := berror.CastToWrappedErrFromErr(testSampleErr)
		if ok || we != nil {
			t.Fatal("expected nil and false for non-wrapped error")
		}
	})

	t.Run("with wrapped error", func(t *testing.T) {
		original := berror.Newf("test")
		we, ok := berror.CastToWrappedErrFromErr(original)
		if !ok || we == nil {
			t.Fatal("expected wrapped error to be cast successfully")
		}
	})
}

func TestEmptyFormatString(t *testing.T) {
	ed := berror.NewErrDefinition("")
	err := ed.New()
	if err == nil {
		t.Fatal("expected non-nil error even with empty format")
	}
	if err.RawError() != "" {
		t.Fatalf("expected empty raw error, got: %v", err.RawError())
	}
}

func TestNilErrDefinition(t *testing.T) {
	var ed *struct {
		berror.ErrDefinition
	}
	// Test that nil ErrDefinition doesn't panic when accessed via interface
	// This is handled internally, but we test the public API behavior
	t.Run("NewErrDefinition returns non-nil", func(t *testing.T) {
		def := berror.NewErrDefinition("test")
		if def == nil {
			t.Fatal("NewErrDefinition should never return nil")
		}
		_ = ed // silence unused warning
	})
}

func TestErrDefinition_Is(t *testing.T) {
	t.Run("with matching error", func(t *testing.T) {
		ed := berror.NewErrDefinition("test error",
			berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock))
		err := ed.New()
		if !ed.Is(err) {
			t.Fatal("expected ErrDefinition.Is to return true for matching error")
		}
	})

	t.Run("with non-matching error", func(t *testing.T) {
		ed := berror.NewErrDefinition("test error",
			berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock))
		if ed.Is(testSampleErr) {
			t.Fatal("expected ErrDefinition.Is to return false for non-matching error")
		}
	})

	t.Run("with nil error", func(t *testing.T) {
		ed := berror.NewErrDefinition("test error",
			berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock))
		if ed.Is(nil) {
			t.Fatal("expected ErrDefinition.Is to return false for nil error")
		}
	})
}

func Test_err_wrapper(t *testing.T) {
	type args struct {
		format string
		opts   []berror.ErrDefinitionOption
		args   []any
	}
	tests := []struct {
		name           string
		args           args
		wantRawErrStr  string
		wantErrStr     string
		wantErrCode    string
		wantErrUnwarp  map[string]bool
		wantErrIs      map[error]bool
		wantStackTrace string
	}{
		{
			name: "simple usage",
			args: args{
				format: testFormat,
				opts:   []berror.ErrDefinitionOption{berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock)},
				args:   nil,
			},
			wantRawErrStr:  "error timeout",
			wantErrStr:     "error timeout (./test)",
			wantErrCode:    "",
			wantErrUnwarp:  map[string]bool{},
			wantErrIs:      map[error]bool{},
			wantStackTrace: "./test",
		},
		{
			name: "usage with args",
			args: args{
				format: testFormatWithArgs,
				opts:   []berror.ErrDefinitionOption{berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock)},
				args:   []interface{}{"arg1", "arg2"},
			},
			wantRawErrStr:  "error timeout: arg1 arg2",
			wantErrStr:     "error timeout: arg1 arg2 (./test)",
			wantErrCode:    "",
			wantErrUnwarp:  map[string]bool{},
			wantErrIs:      map[error]bool{},
			wantStackTrace: "./test",
		},
		{
			name: "usage with args and options",
			args: args{
				format: testFormatWithArgs,
				opts: []berror.ErrDefinitionOption{
					berror.OptionErrDefinitionWithErrCode(testErrorCode),
					berror.OptionErrDefinitionWithCustomFormater(func(err string, code string, stack string) string {
						return fmt.Sprintf("%v--%v--%v", err, code, stack)
					}),
					berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
				},
				args: []interface{}{"arg1", "arg2"},
			},
			wantRawErrStr:  "error timeout: arg1 arg2",
			wantErrStr:     "error timeout: arg1 arg2--t001--./test",
			wantErrCode:    "t001",
			wantErrUnwarp:  map[string]bool{},
			wantErrIs:      map[error]bool{},
			wantStackTrace: "./test",
		},
		{
			name: "unwarp and errors is",
			args: args{
				format: testFormatWithArgs,
				opts: []berror.ErrDefinitionOption{
					berror.OptionErrDefinitionWithErrCode(testErrorCode),
					berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
				},
				args: []interface{}{testSampleErr, testErrDef.New(testSampleErr2)},
			},
			wantRawErrStr: "error timeout: test sample error xxx test sample error 2 (./test)",
			wantErrStr:    "[t001] error timeout: test sample error xxx test sample error 2 (./test) (./test)",
			wantErrCode:   "t001",
			wantErrUnwarp: map[string]bool{
				testSampleErr.Error():                  false,
				testErrDef.New(testSampleErr2).Error(): false,
			},
			wantErrIs: map[error]bool{
				testSampleErr:    true,
				testSampleErr2:   true,
				testErrDef.New(): true,
				testSampleErr3:   false,
				nil:              false,
			},
			wantStackTrace: "./test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := berror.NewErrDefinition(tt.args.format, tt.args.opts...)
			err := d.New(tt.args.args...)

			gotWantRawErrStr := err.RawError()
			if tt.wantRawErrStr != gotWantRawErrStr {
				t.Fatal(fmt.Sprintf("expected wantRawErrStr: %v, got: %v", tt.wantRawErrStr, gotWantRawErrStr))
			}

			gotErrStr := err.Error()
			if tt.wantErrStr != gotErrStr {
				t.Fatal(fmt.Sprintf("expected wantErrStr: %v, got: %v", tt.wantErrStr, gotErrStr))
			}

			gotErrCode := err.Code()
			if tt.wantErrCode != gotErrCode {
				t.Fatal(fmt.Sprintf("expected wantErrCode: %v, got: %v", tt.wantErrCode, gotErrCode))
			}

			for _, e := range err.Unwrap() {
				_, ok := tt.wantErrUnwarp[e.Error()]
				if ok {
					tt.wantErrUnwarp[e.Error()] = true
				} else {
					t.Fatal(fmt.Sprintf("err: %v is not expected by wantErrUnwarp", e.Error()))
				}
			}

			for k, v := range tt.wantErrUnwarp {
				if !v {
					t.Fatal(fmt.Sprintf("err: %v is expected by wantErrUnwarp", k))
				}
			}

			for k, v := range tt.wantErrIs {
				gotErrIs := errors.Is(err, k)
				if gotErrIs != v {
					t.Fatal(fmt.Sprintf("err: %v is expected to be %v by wantErrIs, got: %v", k, v, gotErrIs))
				}

				gotErrIs = err.Is(k)
				if gotErrIs != v {
					t.Fatal(fmt.Sprintf("parent err: %v is expected to be %v by wantErrIs, got: %v", k, v, gotErrIs))
				}
			}

			if !d.Is(err) || err.ErrorDefinition() != d {
				t.Fatal("error definition check fail")
			}

			gotStackTrace := err.StackTrace()
			if tt.wantStackTrace != gotStackTrace {
				t.Fatal(fmt.Sprintf("expected stack trace %v, got: %v", tt.wantStackTrace, gotStackTrace))
			}

			gotString := fmt.Sprintf("%s", err)
			if err.String() != gotString {
				t.Fatal(fmt.Sprintf("expected err String(): %v, got: %v", err.Error(), gotString))
			}
		})
	}
}

func TestWrap(t *testing.T) {
	t.Run("wrap nil error", func(t *testing.T) {
		err := berror.Wrap(nil, "context")
		if err != nil {
			t.Fatal("expected nil when wrapping nil error")
		}
	})

	t.Run("wrap valid error", func(t *testing.T) {
		err := berror.Wrap(testSampleErr, "failed to process")
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if !errors.Is(err, testSampleErr) {
			t.Fatal("wrapped error should contain original error")
		}
		if !containsString(err.RawError(), "failed to process") {
			t.Fatalf("expected message to contain 'failed to process', got: %v", err.RawError())
		}
	})
}

func TestWrapf(t *testing.T) {
	t.Run("wrapf nil error", func(t *testing.T) {
		err := berror.Wrapf(nil, "context %v", "arg1")
		if err != nil {
			t.Fatal("expected nil when wrapping nil error")
		}
	})

	t.Run("wrapf valid error", func(t *testing.T) {
		err := berror.Wrapf(testSampleErr, "failed at step %d", 5)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if !errors.Is(err, testSampleErr) {
			t.Fatal("wrapped error should contain original error")
		}
		if !containsString(err.RawError(), "failed at step 5") {
			t.Fatalf("expected message to contain 'failed at step 5', got: %v", err.RawError())
		}
	})
}

func TestIsCode(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		if berror.IsCode(nil, "E001") {
			t.Fatal("expected false for nil error")
		}
	})

	t.Run("empty code", func(t *testing.T) {
		err := berror.Newf("test")
		if berror.IsCode(err, "") {
			t.Fatal("expected false for empty code")
		}
	})

	t.Run("matching code", func(t *testing.T) {
		err := berror.NewErrDefinition("test error",
			berror.OptionErrDefinitionWithErrCode("E001"),
			berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
		).New()
		if !berror.IsCode(err, "E001") {
			t.Fatal("expected true for matching code")
		}
	})

	t.Run("non-matching code", func(t *testing.T) {
		err := berror.NewErrDefinition("test error",
			berror.OptionErrDefinitionWithErrCode("E001"),
			berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
		).New()
		if berror.IsCode(err, "E002") {
			t.Fatal("expected false for non-matching code")
		}
	})

	t.Run("nested wrapped error with code", func(t *testing.T) {
		innerErr := berror.NewErrDefinition("inner error",
			berror.OptionErrDefinitionWithErrCode("INNER"),
			berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
		).New()
		outerErr := berror.Wrap(innerErr, "outer context")
		if !berror.IsCode(outerErr, "INNER") {
			t.Fatal("expected to find inner error code")
		}
	})

	t.Run("standard error without code", func(t *testing.T) {
		if berror.IsCode(testSampleErr, "E001") {
			t.Fatal("expected false for standard error")
		}
	})
}

func TestGetCode(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		code := berror.GetCode(nil)
		if code != "" {
			t.Fatal("expected empty string for nil error")
		}
	})

	t.Run("non-wrapped error", func(t *testing.T) {
		code := berror.GetCode(testSampleErr)
		if code != "" {
			t.Fatal("expected empty string for non-wrapped error")
		}
	})

	t.Run("wrapped error with code", func(t *testing.T) {
		err := berror.NewErrDefinition("test",
			berror.OptionErrDefinitionWithErrCode("E001"),
			berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
		).New()
		code := berror.GetCode(err)
		if code != "E001" {
			t.Fatalf("expected 'E001', got: %v", code)
		}
	})
}

func TestGetStackTrace(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		stack := berror.GetStackTrace(nil)
		if stack != "" {
			t.Fatal("expected empty string for nil error")
		}
	})

	t.Run("non-wrapped error", func(t *testing.T) {
		stack := berror.GetStackTrace(testSampleErr)
		if stack != "" {
			t.Fatal("expected empty string for non-wrapped error")
		}
	})

	t.Run("wrapped error with stack trace", func(t *testing.T) {
		err := berror.NewErrDefinition("test",
			berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
		).New()
		stack := berror.GetStackTrace(err)
		if stack != testStackTraceCaputurerMockPath {
			t.Fatalf("expected '%v', got: %v", testStackTraceCaputurerMockPath, stack)
		}
	})
}

func TestIsWrappedErr(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		if berror.IsWrappedErr(nil) {
			t.Fatal("expected false for nil error")
		}
	})

	t.Run("standard error", func(t *testing.T) {
		if berror.IsWrappedErr(testSampleErr) {
			t.Fatal("expected false for standard error")
		}
	})

	t.Run("wrapped error", func(t *testing.T) {
		err := berror.Newf("test")
		if !berror.IsWrappedErr(err) {
			t.Fatal("expected true for wrapped error")
		}
	})
}

func TestJoin(t *testing.T) {
	t.Run("all nil errors", func(t *testing.T) {
		err := berror.Join(nil, nil, nil)
		if err != nil {
			t.Fatal("expected nil when all errors are nil")
		}
	})

	t.Run("single error", func(t *testing.T) {
		err := berror.Join(testSampleErr)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if !errors.Is(err, testSampleErr) {
			t.Fatal("expected error to wrap testSampleErr")
		}
	})

	t.Run("single wrapped error", func(t *testing.T) {
		original := berror.Newf("original")
		err := berror.Join(original)
		if err != original {
			t.Fatal("expected same WrappedErr instance for single wrapped error")
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		err := berror.Join(testSampleErr, testSampleErr2, testSampleErr3)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if !errors.Is(err, testSampleErr) {
			t.Fatal("expected error to contain testSampleErr")
		}
		if !errors.Is(err, testSampleErr2) {
			t.Fatal("expected error to contain testSampleErr2")
		}
		if !errors.Is(err, testSampleErr3) {
			t.Fatal("expected error to contain testSampleErr3")
		}
	})

	t.Run("mixed nil and non-nil errors", func(t *testing.T) {
		err := berror.Join(nil, testSampleErr, nil, testSampleErr2, nil)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if !errors.Is(err, testSampleErr) {
			t.Fatal("expected error to contain testSampleErr")
		}
		if !errors.Is(err, testSampleErr2) {
			t.Fatal("expected error to contain testSampleErr2")
		}
	})
}

func TestWithCode(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		err := berror.WithCode(nil, "E001")
		if err != nil {
			t.Fatal("expected nil when error is nil")
		}
	})

	t.Run("valid error with code", func(t *testing.T) {
		err := berror.WithCode(testSampleErr, "E001")
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.Code() != "E001" {
			t.Fatalf("expected code 'E001', got: %v", err.Code())
		}
		if !errors.Is(err, testSampleErr) {
			t.Fatal("expected error to wrap testSampleErr")
		}
	})
}

func TestMust(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		result := berror.Must("success", nil)
		if result != "success" {
			t.Fatalf("expected 'success', got: %v", result)
		}
	})

	t.Run("with error panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		berror.Must("value", testSampleErr)
	})
}

func TestIgnore(t *testing.T) {
	t.Run("returns value ignoring error", func(t *testing.T) {
		result := berror.Ignore("success", testSampleErr)
		if result != "success" {
			t.Fatalf("expected 'success', got: %v", result)
		}
	})

	t.Run("returns value with nil error", func(t *testing.T) {
		result := berror.Ignore(42, nil)
		if result != 42 {
			t.Fatalf("expected 42, got: %v", result)
		}
	})
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
