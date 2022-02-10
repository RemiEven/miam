package testutils

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestEmptyResponseBodyTest(t *testing.T) {
	tests := map[string]struct {
		bodyToCheck     string
		expectedResult  bool
		expectedMessage string
	}{
		"Empty string": {
			bodyToCheck:    "",
			expectedResult: true,
		},
		"Non empty string": {
			bodyToCheck:     "WingardiumLeviosa",
			expectedResult:  false,
			expectedMessage: "expected empty response body but got [WingardiumLeviosa]",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualMessage, actualResult := EmptyResponseBodyTest(test.bodyToCheck)
			if test.expectedMessage != "" && actualMessage != test.expectedMessage {
				t.Errorf("wanted message [%q], got [%q]", test.expectedMessage, actualMessage)
			}
			if actualResult != test.expectedResult {
				t.Errorf("wanted result [%v], got [%v]", test.expectedResult, actualResult)
			}
		})
	}
}

func TestExactResponseBodyTest(t *testing.T) {
	testCases := map[string]struct {
		bodyToCheck     string
		wantedBody      string
		expectedResult  bool
		expectedMessage string
	}{
		"bodies are equal": {
			bodyToCheck:    "All bodies are equal, but some bodies are more equal than others.",
			wantedBody:     "All bodies are equal, but some bodies are more equal than others.",
			expectedResult: true,
		},
		"bodies are not equal": {
			bodyToCheck:     "Political language is designed to make lies sound truthful and murder respectable, and to give an appearance of solidity to pure wind.",
			wantedBody:      "Doublethink means the power of holding two contradictory beliefs in one's mind simultaneously, and accepting both of them.",
			expectedResult:  false,
			expectedMessage: "expected body to be [Doublethink means the power of holding two contradictory beliefs in one's mind simultaneously, and accepting both of them.], but got [Political language is designed to make lies sound truthful and murder respectable, and to give an appearance of solidity to pure wind.]",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			actualMessage, actualResult := ExactResponseBodyTest(test.wantedBody)(test.bodyToCheck)
			if test.expectedMessage != "" && actualMessage != test.expectedMessage {
				t.Errorf("wanted message [%q], got [%q]", test.expectedMessage, actualMessage)
			}
			if actualResult != test.expectedResult {
				t.Errorf("wanted result [%v], got [%v]", test.expectedResult, actualResult)
			}
		})
	}
}

func TestJsonResponseBodyTest(t *testing.T) {
	tests := map[string]struct {
		bodyToCheck         string
		wantedBody          string
		shouldReturnMessage string
		shouldReturnMatch   bool
	}{
		"invalid json in actual": {
			bodyToCheck:         `invalid`,
			wantedBody:          `null`,
			shouldReturnMessage: "failed to parse actual json: invalid character 'i' looking for beginning of value",
		},
		"invalid json in expected": {
			bodyToCheck:         `null`,
			wantedBody:          `invalid`,
			shouldReturnMessage: "failed to parse expected json: invalid character 'i' looking for beginning of value",
		},
		"equal primitive value": {
			bodyToCheck:       `null`,
			wantedBody:        `null`,
			shouldReturnMatch: true,
		},
		"different primitive values": {
			bodyToCheck:         `2`,
			wantedBody:          `null`,
			shouldReturnMessage: `JSON contents do not match: difference(s) found between actual and expected: 2 != <nil pointer>`,
		},
		"equal arrays": {
			bodyToCheck:       `[3, 4]`,
			wantedBody:        `[3, 4]`,
			shouldReturnMatch: true,
		},
		"arrays with different orders": {
			bodyToCheck:         `[3, 4]`,
			wantedBody:          `[4, 3]`,
			shouldReturnMessage: "JSON contents do not match: difference(s) found between actual and expected:\n- slice[0]: 3 != 4\n- slice[1]: 4 != 3",
		},
		"equal objects": {
			bodyToCheck:       `{"a": 3, "b": 4}`,
			wantedBody:        `{"a": 3, "b": 4}`,
			shouldReturnMatch: true,
		},
		"different objects": {
			bodyToCheck:         `{"a": 3, "b": 4}`,
			wantedBody:          `{"a": 4, "b": 4}`,
			shouldReturnMessage: `JSON contents do not match: difference(s) found between actual and expected: map[a]: 3 != 4`,
		},
		"equal complex objects": {
			bodyToCheck:       `{"a": [3, 4], "b": {"c": 5, "d": "some message"}}`,
			wantedBody:        `{"a": [3, 4], "b": {"c": 5, "d": "some message"}}`,
			shouldReturnMatch: true,
		},
		"different complex objects": {
			bodyToCheck:         `{"a": [3, 4], "b": {"c": 5, "d": "some message"}}`,
			wantedBody:          `{"a": [3, 4], "b": {"c": 5, "d": "another message"}}`,
			shouldReturnMessage: `JSON contents do not match: difference(s) found between actual and expected: map[b].map[d]: some message != another message`,
		},
		"different object attributes order": {
			bodyToCheck:       `{"a": 3, "b": 4}`,
			wantedBody:        `{"b": 4, "a": 3}`,
			shouldReturnMatch: true,
		},
		"different whitespacing": {
			bodyToCheck: `{
				"a": 3,
				"b": 4
			}`,
			wantedBody:        `{"b": 4, "a": 3}`,
			shouldReturnMatch: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			message, match := JsonResponseBodyTest(test.wantedBody)(test.bodyToCheck)
			if message != test.shouldReturnMessage {
				t.Errorf("got message [%q], wanted [%q]", message, test.shouldReturnMessage)
			}
			if match != test.shouldReturnMatch {
				t.Errorf("got match [%v], wanted [%v]", match, test.shouldReturnMatch)
			}
		})
	}
}

func TestErrorEqual(t *testing.T) {
	testCases := map[string]struct {
		actualErr    error
		expectedErr  error
		shouldReturn bool
	}{
		"Both nil": {
			actualErr:    nil,
			expectedErr:  nil,
			shouldReturn: true,
		},
		"Both same messages": {
			actualErr:    errors.New("message"),
			expectedErr:  errors.New("message"),
			shouldReturn: true,
		},
		"Actual is nil": {
			actualErr:    nil,
			expectedErr:  errors.New("message"),
			shouldReturn: false,
		},
		"Expected is nil": {
			actualErr:    errors.New("message"),
			expectedErr:  nil,
			shouldReturn: false,
		},
		"Messages differ": {
			actualErr:    errors.New("a message"),
			expectedErr:  errors.New("another message"),
			shouldReturn: false,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			actualReturn := ErrorEqual(test.actualErr, test.expectedErr)
			if actualReturn != test.shouldReturn {
				t.Errorf("Got result [%v], expected [%v]", actualReturn, test.shouldReturn)
			}
		})
	}
}

func TestErrorsEqual(t *testing.T) {
	testCases := map[string]struct {
		actualErrs   []error
		expectedErrs []error
		shouldReturn bool
	}{
		"Both empty": {
			actualErrs:   []error{},
			expectedErrs: []error{},
			shouldReturn: true,
		},
		"more errors than expected": {
			actualErrs:   []error{errors.New("some error")},
			expectedErrs: []error{},
			shouldReturn: false,
		},
		"less errors than expected": {
			actualErrs:   []error{},
			expectedErrs: []error{errors.New("some error")},
			shouldReturn: false,
		},
		"Some errors have different messages": {
			actualErrs:   []error{errors.New("some error"), errors.New("some other error")},
			expectedErrs: []error{errors.New("some error"), errors.New("yet another error")},
			shouldReturn: false,
		},
		"Ordering is not the same": {
			actualErrs:   []error{errors.New("some error"), errors.New("some other error")},
			expectedErrs: []error{errors.New("some other error"), errors.New("some error")},
			shouldReturn: false,
		},
		"Many errors with same messages and ordering": {
			actualErrs:   []error{errors.New("some error"), errors.New("some other error")},
			expectedErrs: []error{errors.New("some error"), errors.New("some other error")},
			shouldReturn: true,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			actualReturn := ErrorsEqual(test.actualErrs, test.expectedErrs)
			if actualReturn != test.shouldReturn {
				t.Errorf("Got result [%v], expected [%v]", actualReturn, test.shouldReturn)
			}
		})
	}
}

func TestJsonEqual(t *testing.T) {
	tests := map[string]struct {
		actualJSON          string
		expectedJSON        string
		shouldReturnMessage string
		shouldReturnMatch   bool
	}{
		"invalid json in actual": {
			actualJSON:          `invalid`,
			expectedJSON:        `null`,
			shouldReturnMessage: "failed to parse actual json: invalid character 'i' looking for beginning of value",
		},
		"invalid json in expected": {
			actualJSON:          `null`,
			expectedJSON:        `invalid`,
			shouldReturnMessage: "failed to parse expected json: invalid character 'i' looking for beginning of value",
		},
		"equal primitive value": {
			actualJSON:        `null`,
			expectedJSON:      `null`,
			shouldReturnMatch: true,
		},
		"different primitive values": {
			actualJSON:          `2`,
			expectedJSON:        `null`,
			shouldReturnMessage: `JSON contents do not match: difference(s) found between actual and expected: 2 != <nil pointer>`,
		},
		"equal arrays": {
			actualJSON:        `[3, 4]`,
			expectedJSON:      `[3, 4]`,
			shouldReturnMatch: true,
		},
		"arrays with different orders": {
			actualJSON:          `[3, 4]`,
			expectedJSON:        `[4, 3]`,
			shouldReturnMessage: "JSON contents do not match: difference(s) found between actual and expected:\n- slice[0]: 3 != 4\n- slice[1]: 4 != 3",
		},
		"equal objects": {
			actualJSON:        `{"a": 3, "b": 4}`,
			expectedJSON:      `{"a": 3, "b": 4}`,
			shouldReturnMatch: true,
		},
		"different objects": {
			actualJSON:          `{"a": 3, "b": 4}`,
			expectedJSON:        `{"a": 4, "b": 4}`,
			shouldReturnMessage: `JSON contents do not match: difference(s) found between actual and expected: map[a]: 3 != 4`,
		},
		"equal complex objects": {
			actualJSON:        `{"a": [3, 4], "b": {"c": 5, "d": "some message"}}`,
			expectedJSON:      `{"a": [3, 4], "b": {"c": 5, "d": "some message"}}`,
			shouldReturnMatch: true,
		},
		"different complex objects": {
			actualJSON:          `{"a": [3, 4], "b": {"c": 5, "d": "some message"}}`,
			expectedJSON:        `{"a": [3, 4], "b": {"c": 5, "d": "another message"}}`,
			shouldReturnMessage: `JSON contents do not match: difference(s) found between actual and expected: map[b].map[d]: some message != another message`,
		},
		"different object attributes order": {
			actualJSON:        `{"a": 3, "b": 4}`,
			expectedJSON:      `{"b": 4, "a": 3}`,
			shouldReturnMatch: true,
		},
		"different whitespacing": {
			actualJSON: `{
				"a": 3,
				"b": 4
			}`,
			expectedJSON:      `{"b": 4, "a": 3}`,
			shouldReturnMatch: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			message, match := JSONEqual(test.actualJSON, test.expectedJSON)
			if message != test.shouldReturnMessage {
				t.Errorf("got message [%q], wanted [%q]", message, test.shouldReturnMessage)
			}
			if match != test.shouldReturnMatch {
				t.Errorf("got match [%v], wanted [%v]", match, test.shouldReturnMatch)
			}
		})
	}
}

type exampleStruct struct {
	FirstName, LastName string
	unexportedName      string
}

func TestDeepEqual(t *testing.T) {
	tests := map[string]struct {
		a, b           interface{}
		expectedResult string
	}{
		"Equal objects": {
			a: exampleStruct{
				FirstName:      "FirstName",
				LastName:       "LastName",
				unexportedName: "Name",
			},
			b: exampleStruct{
				FirstName:      "FirstName",
				LastName:       "LastName",
				unexportedName: "Name",
			},
			expectedResult: ``,
		},
		"Different object types": {
			a: exampleStruct{
				FirstName:      "FirstName",
				LastName:       "LastName",
				unexportedName: "Name",
			},
			b:              time.Time{},
			expectedResult: `difference(s) found between actual and expected: testutils.exampleStruct != time.Time`,
		},
		"Objects with one difference": {
			a: exampleStruct{
				FirstName:      "FirstName",
				LastName:       "LastName",
				unexportedName: "Name",
			},
			b: exampleStruct{
				FirstName:      "Other FirstName",
				LastName:       "LastName",
				unexportedName: "Name",
			},
			expectedResult: `difference(s) found between actual and expected: FirstName: FirstName != Other FirstName`,
		},
		"Objects with several differences": {
			a: exampleStruct{
				FirstName:      "FirstName",
				LastName:       "LastName",
				unexportedName: "Name",
			},
			b: exampleStruct{
				FirstName:      "Other FirstName",
				LastName:       "Other LastName",
				unexportedName: "Name",
			},
			expectedResult: `difference(s) found between actual and expected:
- FirstName: FirstName != Other FirstName
- LastName: LastName != Other LastName`,
		},
		"Objects with difference but it is in an unexported field": {
			a: exampleStruct{
				FirstName:      "FirstName",
				LastName:       "LastName",
				unexportedName: "Name",
			},
			b: exampleStruct{
				FirstName:      "FirstName",
				LastName:       "LastName",
				unexportedName: "Other Secret Name",
			},
			expectedResult: ``,
		},
		"Nil with empty slice": {
			a:              nil,
			b:              []string{},
			expectedResult: `difference(s) found between actual and expected: <nil pointer> != []`,
		},
		"Int with solid float": {
			a:              int(1),
			b:              float64(1.0),
			expectedResult: `difference(s) found between actual and expected: int != float64`,
		},
		"Slices containing different types": {
			a:              []string{},
			b:              []interface{}{},
			expectedResult: `difference(s) found between actual and expected: []string != []interface {}`,
		},
		"Errors with same messages but different underlying types": { // in those cases, avoid testutils.DeepEqual and use testutils.ErrorEqual
			a:              fmt.Errorf("wrapping error: %w", errors.New("wrapped error")),
			b:              errors.New("wrapping error: wrapped error"),
			expectedResult: `difference(s) found between actual and expected: *fmt.wrapError != *errors.errorString`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualResult := DeepEqual(test.a, test.b)

			if actualResult != test.expectedResult {
				t.Errorf("got result [%v], wanted [%v]", actualResult, test.expectedResult)
			}
		})
	}
}
