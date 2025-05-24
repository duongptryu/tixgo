package syserr

import (
	"errors"
	"fmt"
	"strings"

	pkgError "github.com/pkg/errors"
)

type Error struct {
	Message      string
	code         Code
	Stack        []*ErrorStackItem
	fields       []*Field
	WrappedError error
}

type ErrorStackItem struct {
	File     string
	Line     string
	Function string
}

type Field struct {
	Key   string
	Value any
}

func F(key string, value any) *Field {
	return &Field{
		Key:   key,
		Value: value,
	}
}

func New(code Code, message string, fields ...*Field) *Error {
	stack := extractStackFromGenericError(pkgError.New(""))
	return &Error{
		Message: message,
		code:    code,
		fields:  fields,
		Stack:   stack,
	}
}

func Wrap(err error, code Code, message string, fields ...*Field) *Error {
	newError := New(code, message, fields...)
	newError.WrappedError = err
	return newError
}
func WrapAsIs(err error, message string, fields ...*Field) *Error {
	newError := New(extractCodeFromGenericError(err), message, fields...)
	newError.WrappedError = err
	return newError
}

func (e *Error) Unwrap() error {
	return e.WrappedError
}

func (e *Error) Code() Code {
	return e.code
}

func (e *Error) Fields() []*Field {
	return e.fields
}

func (e *Error) StackTrace() []*ErrorStackItem {
	return e.Stack
}

func (e *Error) Error() string {
	if e.WrappedError != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.WrappedError.Error())
	}
	return e.Message
}

func (e *Error) StackFormatted() []string {
	return formatStack(e.Stack)
}

func formatStack(stack []*ErrorStackItem) []string {
	result := make([]string, len(stack))

	for index, stackItem := range stack {
		result[index] = fmt.Sprintf("%s:%s %s", stackItem.File, stackItem.Line, stackItem.Function)
	}

	return result
}

func extractStackFromGenericError(err error) []*ErrorStackItem {
	stackTrace := extractStackTraceFromGenericError(err)

	result := make([]*ErrorStackItem, len(stackTrace))

	for index, frame := range stackTrace {
		result[index] = &ErrorStackItem{
			File:     getFrameFilePath(frame),
			Line:     fmt.Sprintf("%d", frame),
			Function: fmt.Sprintf("%s", frame),
		}
	}

	return result
}

type stackTracer interface {
	StackTrace() pkgError.StackTrace
}

func extractStackTraceFromGenericError(err error) pkgError.StackTrace {
	var result pkgError.StackTrace

	var traceableError stackTracer
	ok := errors.As(err, &traceableError)
	if ok {
		result = traceableError.StackTrace()
	}

	return result
}

func getFrameFilePath(frame pkgError.Frame) string {
	frameString := strings.Split(fmt.Sprintf("%+s", frame), "\n\t")
	return frameString[1]
}

func extractCodeFromGenericError(err error) Code {
	if err == nil {
		return InternalCode
	}

	for {
		if sErr, ok := err.(*Error); ok {
			return sErr.Code()
		}

		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err := x.Unwrap()
			if err == nil {
				return InternalCode
			}
		default:
			return InternalCode
		}
	}
}
