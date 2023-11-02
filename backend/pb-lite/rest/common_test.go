package rest_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/remieven/miam/pb-lite/failure"
	"github.com/remieven/miam/pb-lite/rest"
	"github.com/remieven/miam/pb-lite/testutils"
)

type FakeObject struct {
	ID string
}

func TestNotFoundHandler(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/route-which-does-not-exist?some-params=some-value", nil)
	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()

	rest.NotFoundHandler(rr, request)

	expectedStatus := http.StatusNotFound
	if rr.Code != expectedStatus {
		t.Errorf("got unexpected http status: wanted [%d], got [%d]", expectedStatus, rr.Code)
	}

	resp := rr.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if message, ok := testutils.ErrorResponseBodyTest(failure.ResourceNotFoundErrorCode)(string(body)); !ok {
		t.Errorf("got unexpected http response body: %s", message)
	}

	var errorBody failure.ErrorResponseBody
	if err := json.Unmarshal(body, &errorBody); err != nil {
		t.Error(err)
	}
	expectedMessage := "/route-which-does-not-exist not found"
	if errorBody.Message != expectedMessage {
		t.Errorf("got unexpected error message: wanted [%q], got [%q]", expectedMessage, errorBody.Message)
	}
}

func TestMethodNotAllowedHandler(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/some-route", nil)
	if err != nil {
		t.Error(err)
		return
	}

	rr := httptest.NewRecorder()

	rest.MethodNotAllowedHandler(rr, request)

	expectedStatus := http.StatusMethodNotAllowed
	if rr.Code != expectedStatus {
		t.Errorf("got unexpected http status: wanted [%d], got [%d]", expectedStatus, rr.Code)
	}

	resp := rr.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if message, ok := testutils.ErrorResponseBodyTest(failure.MethodNotAllowedErrorCode)(string(body)); !ok {
		t.Errorf("got unexpected http response body: %s", message)
	}

	var errorBody failure.ErrorResponseBody
	if err := json.Unmarshal(body, &errorBody); err != nil {
		t.Error(err)
	}
	expectedMessage := "GET not allowed for this route"
	if errorBody.Message != expectedMessage {
		t.Errorf("got unexpected error message: wanted [%q], got [%q]", expectedMessage, errorBody.Message)
	}
}

func TestWriteNoContentResponse(t *testing.T) {
	_, err := http.NewRequest(http.MethodGet, "/some-route", nil)
	if err != nil {
		t.Error(err)
		return
	}
	rr := httptest.NewRecorder()

	rest.WriteNoContentResponse(rr)
	if rr.Code != http.StatusNoContent {
		t.Errorf("got unexpected http status: wanted [%d], got [%d]", http.StatusNoContent, rr.Code)
	}
}

func TestWriteCreatedResponse(t *testing.T) {
	_, err := http.NewRequest(http.MethodGet, "/some-route", nil)
	if err != nil {
		t.Error(err)
		return
	}
	rr := httptest.NewRecorder()

	rest.WriteCreatedResponse(rr, &http.Request{}, "/some-location")
	if rr.Code != http.StatusCreated {
		t.Errorf("got unexpected http status: wanted [%d], got [%d]", http.StatusCreated, rr.Code)
	}
	expectedLocationHeader := "/some-location"
	actualLocationHeader := rr.Header().Get(rest.HeaderLocation)
	if actualLocationHeader != expectedLocationHeader {
		t.Errorf("unexpected location header: got [%v], wanted [%v]", actualLocationHeader, expectedLocationHeader)
	}
}

func TestWriteCreatedResponseWithPrefix(t *testing.T) {
	_, err := http.NewRequest(http.MethodGet, "/some-route", nil)
	if err != nil {
		t.Error(err)
		return
	}
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.Header.Add("X-Forwarded-Prefix", "/some-prefix/")

	rest.WriteCreatedResponse(rr, req, "/some-location")
	if rr.Code != http.StatusCreated {
		t.Errorf("got unexpected http status: wanted [%d], got [%d]", http.StatusCreated, rr.Code)
	}
	expectedLocationHeader := "/some-prefix/some-location"
	actualLocationHeader := rr.Header().Get(rest.HeaderLocation)
	if actualLocationHeader != expectedLocationHeader {
		t.Errorf("unexpected location header: got [%v], wanted [%v]", actualLocationHeader, expectedLocationHeader)
	}
}

func TestWriteOKResponse(t *testing.T) {
	_, err := http.NewRequest(http.MethodGet, "/some-route", nil)
	if err != nil {
		t.Error(err)
		return
	}
	rr := httptest.NewRecorder()

	rest.WriteOKResponse(rr, FakeObject{ID: "id"})
	if rr.Code != http.StatusOK {
		t.Errorf("got unexpected http status: wanted [%d], got [%d]", http.StatusOK, rr.Code)
	}
	expectedContentTypeHeader := rest.ContentTypeJSONUTF8
	actualContentTypeHeader := rr.Header().Get(rest.HeaderContentType)
	if actualContentTypeHeader != expectedContentTypeHeader {
		t.Errorf("unexpected content-type header: got [%v], wanted [%v]", actualContentTypeHeader, expectedContentTypeHeader)
	}
}

func TestWriteParseBodyErrorResponse(t *testing.T) {
	tests := map[string]struct {
		err              error
		responseBodyTest func(string) (string, bool)
		expectedStatus   int
	}{
		"Syntax Error": {
			err:              json.Unmarshal([]byte("{"), &FakeObject{}),
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.InvalidJSONErrorCode),
			expectedStatus:   http.StatusBadRequest,
		},
		"Unexpected EOF": {
			err: func() error {
				reader := strings.NewReader("{")
				return json.NewDecoder(reader).Decode(&FakeObject{})
			}(),
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.InvalidJSONErrorCode),
			expectedStatus:   http.StatusBadRequest,
		},
		"EOF due to empty input": {
			err: func() error {
				reader := strings.NewReader("")
				return json.NewDecoder(reader).Decode(&FakeObject{})
			}(),
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.InvalidJSONErrorCode),
			expectedStatus:   http.StatusBadRequest,
		},
		"Invalid Unmarshal Error": {
			err:              &json.InvalidUnmarshalError{},
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.InternalErrorErrorCode),
			expectedStatus:   http.StatusInternalServerError,
		},
		"Invalid Argument Error": {
			err:              testutils.ErrSampleTechnical,
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.InvalidArgumentErrorCode),
			expectedStatus:   http.StatusBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			if !rest.HandleParseBodyErrorCase(rr, test.err) {
				rest.WriteNoContentResponse(rr)
			}

			body, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
				return
			}
			if rr.Code != test.expectedStatus {
				t.Errorf("got unexpected http status: wanted [%d], got [%d]", test.expectedStatus, rr.Code)
			}
			if message, ok := test.responseBodyTest(string(body)); !ok {
				t.Errorf("got unexpected http response body: %s", message)
			}
		})
	}
}

func TestHandleErrorCase(t *testing.T) {
	tests := map[string]struct {
		err              error
		expectedStatus   int
		responseBodyTest func(string) (string, bool)
	}{
		"no error": {
			err:              nil,
			expectedStatus:   http.StatusNoContent,
			responseBodyTest: testutils.EmptyResponseBodyTest,
		},
		"resource not found": {
			err: &failure.ResourceNotFoundError{
				Message: "some resource not found",
			},
			expectedStatus:   http.StatusNotFound,
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.ResourceNotFoundErrorCode),
		},
		"invalid value error": {
			err: &failure.InvalidValueError{
				Message: "missing field error",
			},
			expectedStatus:   http.StatusBadRequest,
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.InvalidArgumentErrorCode),
		},
		"other error": {
			err:              testutils.ErrSampleTechnical,
			expectedStatus:   http.StatusInternalServerError,
			responseBodyTest: testutils.ErrorResponseBodyTest(failure.InternalErrorErrorCode),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			rr := httptest.NewRecorder()
			if !rest.HandleErrorCase(rr, test.err) {
				rest.WriteNoContentResponse(rr)
			}

			body, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
				return
			}
			if rr.Code != test.expectedStatus {
				t.Errorf("got unexpected http status: wanted [%d], got [%d]", test.expectedStatus, rr.Code)
			}
			if message, ok := test.responseBodyTest(string(body)); !ok {
				t.Errorf("got unexpected http response body: %s", message)
			}
		})
	}
}

func BenchmarkHandleErrorCase(b *testing.B) {
	errorList := []error{
		nil,
		&failure.ResourceNotFoundError{},
		&failure.InvalidValueError{},
	}
	errorsLen := len(errorList)
	rr := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		rest.HandleErrorCase(rr, errorList[n%errorsLen])
	}
}
