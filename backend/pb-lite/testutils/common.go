package testutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/go-test/deep"

	"github.com/remieven/miam/pb-lite/failure"
)

// ErrorResponseBodyTest creates a function checking that a string match an error response body with the wanted error code
func ErrorResponseBodyTest(code failure.ErrorCode) func(string) (string, bool) {
	return func(body string) (string, bool) {
		var errorBody failure.ErrorResponseBody
		if err := json.Unmarshal([]byte(body), &errorBody); err != nil {
			return fmt.Sprintf("failed to parse body [%q]: %+v", body, err), false
		}
		if errorBody.Code == "" {
			return fmt.Sprintf("expected an error response body with code [%q], but got body [%q]", code, body), false
		}
		if errorBody.Code != code {
			return fmt.Sprintf("got error response body with code [%q] (message: [%q]), expected [%q]", errorBody.Code, errorBody.Message, code), false
		}
		if errorBody.Message == "" {
			return "got error response body with an empty or missing message", false
		}
		return "", true
	}
}

// EmptyResponseBodyTest checks that a string is empty
func EmptyResponseBodyTest(body string) (string, bool) {
	return fmt.Sprintf("expected empty response body but got [%v]", body), len(body) == 0
}

// ExactResponseBodyTest creates a function checking that a string is exactly equal to the one given as parameter
func ExactResponseBodyTest(expectedBody string) func(string) (string, bool) {
	return func(body string) (string, bool) {
		return fmt.Sprintf("expected body to be [%v], but got [%v]", expectedBody, body), body == expectedBody
	}
}

// JsonResponseBodyTest creates a function checking that the given body is valid JSON and that its content is semantically equivalent to expectedJsonBody
func JsonResponseBodyTest(expectedJsonBody string) func(string) (string, bool) {
	return func(body string) (string, bool) {
		return JSONEqual(body, expectedJsonBody)
	}
}

// ErrorsEqual checks whether two error slices contains errors with the same messages and in the same order
func ErrorsEqual(actual, expected []error) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range expected {
		if !ErrorEqual(actual[i], expected[i]) {
			return false
		}
	}
	return true
}

// ErrorEqual checks whether two errors have the same message (or are both nil)
func ErrorEqual(actual, expected error) bool {
	if actual == nil && expected == nil {
		return true
	}
	if (actual != nil && expected == nil) || (actual == nil && expected != nil) {
		return false
	}
	if actual.Error() != expected.Error() {
		return false
	}
	return true
}

// JSONEqual checks whether two strings are valid JSON and that their contents are semantically equivalent
func JSONEqual(actual, expected string) (string, bool) {
	var parsedActual, parsedExpected interface{}

	if err := json.Unmarshal([]byte(actual), &parsedActual); err != nil {
		return fmt.Sprintf("failed to parse actual json: %v", err), false
	}
	if err := json.Unmarshal([]byte(expected), &parsedExpected); err != nil {
		return fmt.Sprintf("failed to parse expected json: %v", err), false
	}

	if diff := DeepEqual(parsedActual, parsedExpected); diff != "" {
		return "JSON contents do not match: " + diff, false
	}

	return "", true
}

// CheckStructValidation checks if the given model is valid and matches the given expected errors
func CheckStructValidation(t *testing.T, modelToValidate interface{}, expectedErrorFields []string) {
	expectError := len(expectedErrorFields) != 0
	err := validator.New().Struct(modelToValidate)
	if err != nil && !expectError {
		t.Errorf("got error [%v] while expecting none", err)
	}
	if err == nil && expectError {
		t.Errorf("got no error while expecting one")
	}
	if err == nil {
		return
	}
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Errorf("expected validationErrors but got [%v]", err)
	}
	if len(validationErrors) != len(expectedErrorFields) {
		t.Errorf("expect %d but got %d validationErrors", len(expectedErrorFields), len(validationErrors))
	}
	for _, err := range validationErrors {
		if !containFields(err.Field(), expectedErrorFields) {
			t.Errorf("got unexpected [%q] error field ", err.Field())
		}
	}
}

func containFields(field string, expectedErrorFields []string) bool {
	for _, expectedErrorField := range expectedErrorFields {
		if field == expectedErrorField {
			return true
		}
	}
	return false
}

// DeepEqual compares a and b and returns a formatted string explaining the differences if there are any
func DeepEqual(actual, expected interface{}) string {
	diff := deep.Equal(actual, expected)
	diffHeader := "difference(s) found between actual and expected:"
	switch len(diff) {
	case 0:
		return ""
	case 1:
		return diffHeader + " " + diff[0]
	default:
		return diffHeader + "\n- " + strings.Join(diff, "\n- ")
	}
}

// WriteMessageIDs write a sample response sent by the Pub/Sub API when publishing a message
func WriteMessageIDs(responseWriter http.ResponseWriter) error {
	_, err := fmt.Fprintln(responseWriter, `{
  "messageIds": [
   "755993117165864"
  ]
 }`)
	return err
}

// ErrSampleTechnical can be used in mocks when simulating errors
var ErrSampleTechnical = errors.New("there's been a great disturbance in the Force")
