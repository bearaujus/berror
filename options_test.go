package berror_test

import (
	"fmt"
	"github.com/bearaujus/berror"
	"testing"
)

func Test_err_wrapper_options(t *testing.T) {
	type args struct {
		opts []berror.ErrDefinitionOption
	}
	tests := []struct {
		name               string
		args               args
		wantRawErrStr      string
		wantErrStr         string
		wantErrCode        string
		ignoreErrStrAssert bool
	}{
		{
			name: "test OptionErrDefinitionWithErrCode",
			args: args{[]berror.ErrDefinitionOption{
				berror.OptionErrDefinitionWithErrCode(testErrorCode),
				berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
			}},
			wantRawErrStr:      "test",
			wantErrStr:         "[t001] test (./test)",
			wantErrCode:        "t001",
			ignoreErrStrAssert: false,
		},
		{
			name: "test empty args OptionErrDefinitionWithErrCode",
			args: args{[]berror.ErrDefinitionOption{
				berror.OptionErrDefinitionWithErrCode(""),
				berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
			}},
			wantRawErrStr:      "test",
			wantErrStr:         "test (./test)",
			wantErrCode:        "",
			ignoreErrStrAssert: false,
		},
		{
			name: "test OptionErrDefinitionWithCustomFormater",
			args: args{[]berror.ErrDefinitionOption{
				berror.OptionErrDefinitionWithErrCode(testErrorCode),
				berror.OptionErrDefinitionWithCustomFormater(berror.ErrWrapperFormatterJSON),
				berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
			}},
			wantRawErrStr:      "test",
			wantErrStr:         `{"code":"t001","err":"test","stack":"./test"}`,
			wantErrCode:        "t001",
			ignoreErrStrAssert: false,
		},
		{
			name: "test empty args OptionErrDefinitionWithCustomFormater",
			args: args{[]berror.ErrDefinitionOption{
				berror.OptionErrDefinitionWithCustomFormater(nil),
				berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
			}},
			wantRawErrStr:      "test",
			wantErrStr:         "test (./test)",
			wantErrCode:        "",
			ignoreErrStrAssert: false,
		},
		{
			name: "test OptionErrDefinitionWithDisabledStackTrace",
			args: args{[]berror.ErrDefinitionOption{
				berror.OptionErrDefinitionWithDisabledStackTrace(),
				berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
			}},
			wantRawErrStr:      "test",
			wantErrStr:         "test",
			wantErrCode:        "",
			ignoreErrStrAssert: false,
		},
		{
			name: "test OptionErrDefinitionWithCustomStackTraceCapturer",
			args: args{[]berror.ErrDefinitionOption{
				berror.OptionErrDefinitionWithCustomStackTraceCapturer(stackTraceCapturerMock),
			}},
			wantRawErrStr:      "test",
			wantErrStr:         "test (./test)",
			wantErrCode:        "",
			ignoreErrStrAssert: false,
		},
		{
			name: "test empty args OptionErrDefinitionWithCustomStackTraceCapturer",
			args: args{[]berror.ErrDefinitionOption{
				berror.OptionErrDefinitionWithCustomStackTraceCapturer(nil),
			}},
			wantRawErrStr:      "test",
			wantErrStr:         "test (./test)",
			wantErrCode:        "",
			ignoreErrStrAssert: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := berror.NewErrDefinition("test", tt.args.opts...).New()
			gotWantRawErrStr := err.RawError()
			if tt.wantRawErrStr != gotWantRawErrStr {
				t.Fatal(fmt.Sprintf("expected wantRawErrStr: %v, got: %v", tt.wantRawErrStr, gotWantRawErrStr))
			}

			gotErrStr := err.Error()
			if tt.wantErrStr != gotErrStr && !tt.ignoreErrStrAssert {
				t.Fatal(fmt.Sprintf("expected wantErrStr: %v, got: %v", tt.wantErrStr, gotErrStr))
			}

			gotErrCode := err.Code()
			if tt.wantErrCode != gotErrCode {
				t.Fatal(fmt.Sprintf("expected wantErrCode: %v, got: %v", tt.wantErrCode, gotErrCode))
			}
		})
	}
}
