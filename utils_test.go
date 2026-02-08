package berror_test

import (
	"github.com/bearaujus/berror"
	"testing"
)

func Test_err_wrapper_util(t *testing.T) {
	t.Run("test ErrWrapperFormatterDefault", func(t *testing.T) {
		gotErrWrapperFormatterDefault := berror.ErrWrapperFormatterDefault("test1", "test2", "test3")
		expectedErrWrapperFormatterDefault := "[test2] test1 (test3)"
		if gotErrWrapperFormatterDefault != expectedErrWrapperFormatterDefault {
			t.Errorf("gotErrWrapperFormatterDefault = %v, want %v", gotErrWrapperFormatterDefault, expectedErrWrapperFormatterDefault)
		}
	})

	t.Run("test ErrWrapperFormatterDefault 2", func(t *testing.T) {
		gotErrWrapperFormatterDefault := berror.ErrWrapperFormatterDefault("test1", "", "")
		expectedErrWrapperFormatterDefault := "test1"
		if gotErrWrapperFormatterDefault != expectedErrWrapperFormatterDefault {
			t.Errorf("gotErrWrapperFormatterDefault = %v, want %v", gotErrWrapperFormatterDefault, expectedErrWrapperFormatterDefault)
		}
	})

	t.Run("test ErrWrapperFormatterDefault 3", func(t *testing.T) {
		gotErrWrapperFormatterDefault := berror.ErrWrapperFormatterDefault("test1", "test2", "")
		expectedErrWrapperFormatterDefault := "[test2] test1"
		if gotErrWrapperFormatterDefault != expectedErrWrapperFormatterDefault {
			t.Errorf("gotErrWrapperFormatterDefault = %v, want %v", gotErrWrapperFormatterDefault, expectedErrWrapperFormatterDefault)
		}
	})

	t.Run("test ErrWrapperFormatterDefault 4", func(t *testing.T) {
		gotErrWrapperFormatterDefault := berror.ErrWrapperFormatterDefault("test1", "", "test3")
		expectedErrWrapperFormatterDefault := "test1 (test3)"
		if gotErrWrapperFormatterDefault != expectedErrWrapperFormatterDefault {
			t.Errorf("gotErrWrapperFormatterDefault = %v, want %v", gotErrWrapperFormatterDefault, expectedErrWrapperFormatterDefault)
		}
	})

	t.Run("test ErrWrapperFormatterJSON", func(t *testing.T) {
		gotErrWrapperFormatterDefault := berror.ErrWrapperFormatterJSON("test1", "test2", "test3")
		expectedErrWrapperFormatterDefault := `{"code":"test2","err":"test1","stack":"test3"}`
		if gotErrWrapperFormatterDefault != expectedErrWrapperFormatterDefault {
			t.Errorf("gotErrWrapperFormatterDefault = %v, want %v", gotErrWrapperFormatterDefault, expectedErrWrapperFormatterDefault)
		}
	})
}
