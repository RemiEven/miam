package failure

import "fmt"

// ErrorCode is a type used for the codes in ErrorResponseBody
type ErrorCode string

// ErrorResponseBody is the structure to use as HTTP body when responding an error to a request
type ErrorResponseBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message,omitempty"`
}

func (erb *ErrorResponseBody) IsOneOf(codes ...ErrorCode) bool {
	if erb == nil {
		return false
	}
	for _, code := range codes {
		if erb.Code == code {
			return true
		}
	}
	return false
}

var _ error = (*InvalidValueError)(nil) // ensure that *InvalidValueError implements error

// InvalidValueError is used when a value is invalid (error on type/pattern, missing value...)
type InvalidValueError struct {
	Message string
	Cause   error
}

// Error is used to implement the error interface
func (err *InvalidValueError) Error() string {
	return MergeErrorWithCause(err.Message, err.Cause)
}

// Unwrap is used to implement the errors.Unwrap interface
func (err *InvalidValueError) Unwrap() error {
	return err.Cause
}

// Is compares an error to the InvalidValueError type, and returns true if they are the same
func (err *InvalidValueError) Is(otherErr error) bool {
	_, ok := otherErr.(*InvalidValueError)
	return ok
}

var _ error = (*ResourceNotFoundError)(nil) // ensure that *ResourceNotFoundError implements error

// ResourceNotFoundError is used when a resource is not found
type ResourceNotFoundError struct {
	Message string
	Cause   error
}

// Error is used to implement the error interface
func (err *ResourceNotFoundError) Error() string {
	return MergeErrorWithCause(err.Message, err.Cause)
}

// Unwrap is used to implement the errors.Unwrap interface
func (err *ResourceNotFoundError) Unwrap() error {
	return err.Cause
}

// Is compares an error to the ResourceNotFoundError type, and returns true if they are the same
func (err *ResourceNotFoundError) Is(otherErr error) bool {
	_, ok := otherErr.(*ResourceNotFoundError)
	return ok
}

// MergeErrorWithCause is used to merge error message with a cause if there is one
func MergeErrorWithCause(message string, cause error) string {
	if cause != nil {
		return fmt.Sprintf("%s: %s", message, cause.Error())
	}
	return message
}
