package syserr

import "errors"

func GetStackFormattedFromGenericError(err error) []string {
	var sysErr *Error
	if errors.As(err, &sysErr) {
		return sysErr.StackFormatted()
	}

	return formatStack(extractStackFromGenericError(err))
}

func GetCodeFromGenericError(err error) Code {
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

func GetFieldsFromGenericError(err error) []*Field {
	var result []*Field

	for {
		if err == nil {
			return result
		}

		if sErr, ok := err.(*Error); ok {
			result = append(result, sErr.Fields()...)
		}

		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		default:
			return result
		}
	}
}
