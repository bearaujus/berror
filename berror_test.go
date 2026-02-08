package berror_test

import (
	"errors"
	"fmt"
	"github.com/bearaujus/berror"
	"testing"
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
