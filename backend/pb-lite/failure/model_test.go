package failure

import (
	"errors"
	"testing"
)

func TestResourceNotFound(t *testing.T) {
 testCases := map[string]struct {
  cause error
  expectedMessage string
 }{
  "Error with cause": {
   cause: errors.New("the cause"),
   expectedMessage: "resource not found: the cause",
  },
  "Error without cause": {
   cause: nil,
   expectedMessage: "resource not found",
  },
 }

 for name, test := range testCases {
  t.Run(name, func(t *testing.T) {
   resourceNotFoundError := ResourceNotFoundError{
    Message: "resource not found",
    Cause: test.cause,
   }

   unwrappedCause := resourceNotFoundError.Unwrap()
   if unwrappedCause != test.cause {
    t.Errorf("Unwrap returned cause [%v], expected [%v]", unwrappedCause, test.cause)
   }

   actualMessage := resourceNotFoundError.Error()
   if actualMessage != test.expectedMessage {
    t.Errorf("Got message [%q], expected [%q]", actualMessage, test.expectedMessage)
   }

   if !errors.Is(&resourceNotFoundError, &ResourceNotFoundError{}) {
    t.Error("error must be a ResourceNotFoundError")
   }
  })
 }
}

func TestInvalidValueError(t *testing.T) {
 testCases := map[string]struct {
  cause error
  expectedMessage string
 }{
  "Error with cause": {
   cause: errors.New("the cause"),
   expectedMessage: "invalid value: the cause",
  },
  "Error without cause": {
   cause: nil,
   expectedMessage: "invalid value",
  },
 }

 for name, test := range testCases {
  t.Run(name, func(t *testing.T) {
   invalidValueError := InvalidValueError{
    Message: "invalid value",
    Cause: test.cause,
   }

   unwrappedCause := invalidValueError.Unwrap()
   if unwrappedCause != test.cause {
    t.Errorf("Unwrap returned cause [%v], expected [%v]", unwrappedCause, test.cause)
   }

   actualMessage := invalidValueError.Error()
   if actualMessage != test.expectedMessage {
    t.Errorf("Got message [%q], expected [%q]", actualMessage, test.expectedMessage)
   }

   if !errors.Is(&invalidValueError, &InvalidValueError{}) {
    t.Error("error must be a InvalidValueError")
   }
  })
 }
}

func TestIsOneOf(t *testing.T) {
 tests := map[string]struct {
  erb *ErrorResponseBody
  codes []ErrorCode
  expectedResult bool
 }{
  "nil receiver": {
   erb: nil,
   codes: []ErrorCode{InternalErrorErrorCode},
   expectedResult: false,
  },
  "code is one of wanted ones": {
   erb: &ErrorResponseBody{
    Code: InvalidJSONErrorCode,
   },
   codes: []ErrorCode{InvalidJSONErrorCode, InvalidArgumentErrorCode},
   expectedResult: true,
  },
  "code is not one of wanted ones": {
   erb: &ErrorResponseBody{
    Code: InternalErrorErrorCode,
   },
   codes: []ErrorCode{ResourceNotFoundErrorCode},
   expectedResult: false,
  },
 }

 for name, test := range tests {
  t.Run(name, func(t *testing.T) {
   actual := test.erb.IsOneOf(test.codes...)
   if actual != test.expectedResult {
    t.Errorf("unexpected result: got [%t], wanted [%t]", actual, test.expectedResult)
   }
  })
 }
}
