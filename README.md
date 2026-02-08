# BError - Error Utilities in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bearaujus/berror/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bearaujus/berror)](https://goreportcard.com/report/github.com/bearaujus/berror)
[![Go Reference](https://pkg.go.dev/badge/github.com/bearaujus/berror.svg)](https://pkg.go.dev/github.com/bearaujus/berror)

This package provides advanced error handling utilities, allowing developers to create structured, formatted, and traceable errors with minimal effort. It supports features like custom error definitions, error codes, stack traces, and flexible formatting.

## Installation

```shell
go get github.com/bearaujus/berror
```

## Import

```go
import "github.com/bearaujus/berror"
```

## Quick Start

```go
// Simple error wrapping
err := berror.New(originalErr)

// Create formatted error
err := berror.Newf("failed to process: %v", details)

// Wrap with context
err := berror.Wrap(originalErr, "failed to connect")

// Wrap with formatted context
err := berror.Wrapf(originalErr, "failed at step %d", stepNum)
```

## Error Definitions

### 1. Basic Error Definitions

```go
var (
    ErrStartProcess = berror.NewErrDefinition("failed to start process: %v")
    ErrReadFile     = berror.NewErrDefinition("failed to read file: %v")
    ErrOpenFile     = berror.NewErrDefinition("failed to open file: %v")
)
```

### 2. Error Definitions with Options

```go
var (
    ErrWithCode = berror.NewErrDefinition("error with code: %v",
        berror.OptionErrDefinitionWithErrCode("E1001"))

    ErrNoStackTrace = berror.NewErrDefinition("no stack trace: %v",
        berror.OptionErrDefinitionWithErrCode("E1002"),
        berror.OptionErrDefinitionWithDisabledStackTrace())

    ErrCustomFormat = berror.NewErrDefinition("custom format: %v",
        berror.OptionErrDefinitionWithErrCode("E1003"),
        berror.OptionErrDefinitionWithCustomFormater(func(err, code, stack string) string {
            return fmt.Sprintf("Custom: %v | Code: %v | Trace: %v", err, code, stack)
        }))

    ErrJSONFormat = berror.NewErrDefinition("json format: %v",
        berror.OptionErrDefinitionWithErrCode("E1004"),
        berror.OptionErrDefinitionWithCustomFormater(berror.ErrWrapperFormatterJSON))
)
```

### 3. Using Error Definitions

```go
func StartProcess(path string) error {
    data, err := ReadFile(path)
    if err != nil {
        return ErrStartProcess.New(err)
    }
    _ = data
    return nil
}

func ReadFile(path string) ([]byte, error) {
    file, err := OpenFile(path)
    if err != nil {
        return nil, ErrReadFile.New(err)
    }
    defer file.Close()
    return io.ReadAll(file)
}

func OpenFile(path string) (*os.File, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, ErrOpenFile.New(err)
    }
    return file, nil
}
```

## Helper Functions

### Wrap and Wrapf

```go
// Wrap error with message
err := berror.Wrap(originalErr, "database connection failed")

// Wrap error with formatted message
err := berror.Wrapf(originalErr, "failed to process user %s", userID)
```

### WithCode

```go
// Add error code to existing error
err := berror.WithCode(originalErr, "E1001")
```

### Join

```go
// Combine multiple errors
func ValidateUser(user *User) error {
    var errs []error
    if user.Name == "" {
        errs = append(errs, berror.Newf("name is required"))
    }
    if user.Email == "" {
        errs = append(errs, berror.Newf("email is required"))
    }
    return berror.Join(errs...) // Returns nil if no errors
}
```

### IsCode and GetCode

```go
// Check if error has specific code
if berror.IsCode(err, "E1001") {
    // Handle specific error
}

// Get error code
code := berror.GetCode(err)
```

### GetStackTrace

```go
// Get stack trace from error
stack := berror.GetStackTrace(err)
```

### IsWrappedErr

```go
// Check if error is a WrappedErr
if berror.IsWrappedErr(err) {
    // Handle wrapped error
}
```

### Must and Ignore

```go
// Panic if error is not nil (useful for initialization)
config := berror.Must(LoadConfig("config.yaml"))

// Ignore error and return only value (use with caution)
value := berror.Ignore(ParseOptional(input))
```

## Error Comparison

```go
err := StartProcess("invalid-path")

// Using errors.Is
fmt.Printf("Matches os.ErrNotExist: %v\n", errors.Is(err, os.ErrNotExist))

// Using ErrDefinition.Is
fmt.Printf("Matches ErrStartProcess: %v\n", ErrStartProcess.Is(err))
fmt.Printf("Matches ErrOpenFile: %v\n", ErrOpenFile.Is(err))
```

## Output Examples

```shell
# Basic errors
failed to start process: Invalid input (./main.go:40)
failed to start process: file already closed (./main.go:41)

# Error with code
[E1001] error with code: Resource not found (./main.go:46)

# No stack trace
[E1002] no stack trace: Operation timed out

# Custom format
Custom: custom format: example | Code: E1003 | Trace: ./main.go:48

# JSON format
{"code":"E1004","err":"json format: error","stack":"./main.go:49"}

# Nested errors
failed to start process: failed to read file: failed to open file: open invalid-path: file not found (./main.go:91) (./main.go:81) (./main.go:71)
```

## Stack Trace Tips

By default, stack traces include full file paths. Use `-trimpath` flag for cleaner output:

```shell
# Without -trimpath
Error: something failed (D:/Projects/myapp/internal/service.go:45)

# With -trimpath
go build -trimpath -o myapp .
Error: something failed (internal/service.go:45)
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/bearaujus/berror/blob/master/LICENSE) file for details.
