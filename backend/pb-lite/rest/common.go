package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/remieven/miam/pb-lite/failure"
)

const (
	// HeaderContentType is the HTTP header name to use to specify the content type of the body
	HeaderContentType = "Content-Type"
	// HeaderLocation is the HTTP header name to use to specify the location of a resource
	HeaderLocation = "Location"
	// HeaderForwardedPrefix is the HTTP header name to look for when sending href being a proxy
	HeaderForwardedPrefix = "X-Forwarded-Prefix"

	// ContentTypeJSONUTF8 is the content type value to use for a JSON UTF-8 body
	ContentTypeJSONUTF8 = "application/json;charset=utf-8"
)

// slashDedupRegexis used to sanitize urls containing duplicate /
var slashDedupRegexp = regexp.MustCompile(`/+`)

// NotFoundHandler handles a request that did not match any other handler
func NotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	WriteErrorResponse(writer, http.StatusNotFound, failure.ResourceNotFoundErrorCode, request.URL.Path+" not found")
}

// MethodNotAllowedHandler handles a request to a route with an incorrect HTTP method
func MethodNotAllowedHandler(writer http.ResponseWriter, request *http.Request) {
	WriteErrorResponse(writer, http.StatusMethodNotAllowed, failure.MethodNotAllowedErrorCode, request.Method+" not allowed for this route")
}

// WriteErrorResponse writes an error response from the given information
func WriteErrorResponse(writer http.ResponseWriter, status int, code failure.ErrorCode, message string) {
	writer.Header().Set(HeaderContentType, ContentTypeJSONUTF8)
	writer.WriteHeader(status)
	errorResponseBody := failure.ErrorResponseBody{
		Code:    code,
		Message: message,
	}
	if encodingError := json.NewEncoder(writer).Encode(errorResponseBody); encodingError != nil {
		log.Error().Err(encodingError).Msg("failed to encode error response body")
	}
}

// WriteOKResponse writes a 200 response
func WriteOKResponse(writer http.ResponseWriter, body any) {
	writeSuccessWithContent(writer, body, http.StatusOK)
}

// writeSuccessWithContent writes a success response with a body result
func writeSuccessWithContent(writer http.ResponseWriter, body any, status int) {
	writer.Header().Set(HeaderContentType, ContentTypeJSONUTF8)
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(body); err != nil {
		WriteErrorResponse(writer, http.StatusInternalServerError, failure.InternalErrorErrorCode, "failed to encode body object")
	}
}

// WriteCreatedResponse writes a 201 response
func WriteCreatedResponse(writer http.ResponseWriter, request *http.Request, location string) {
	if forwardedPrefix := request.Header.Get(HeaderForwardedPrefix); forwardedPrefix != "" {
		location = slashDedupRegexp.ReplaceAllString(forwardedPrefix+"/"+location, "/")
	}
	writer.Header().Set(HeaderLocation, location)
	writer.WriteHeader(http.StatusCreated)
}

// WriteNoContentResponse writes a 204 response
func WriteNoContentResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusNoContent)
}

// HandleParseBodyErrorCase writes an appropriate respose body and returns true if there is an error during the parsing of a JSON request body
func HandleParseBodyErrorCase(writer http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	var (
		syntaxErr           *json.SyntaxError
		invalidUnmarshalErr *json.InvalidUnmarshalError
	)
	switch {
	case errors.As(err, &invalidUnmarshalErr):
		WriteErrorResponse(writer, http.StatusInternalServerError, failure.InternalErrorErrorCode, err.Error())
	case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF), errors.Is(err, io.EOF):
		WriteErrorResponse(writer, http.StatusBadRequest, failure.InvalidJSONErrorCode, err.Error())
	default:
		WriteErrorResponse(writer, http.StatusBadRequest, failure.InvalidArgumentErrorCode, err.Error())
	}

	return true
}

var (
	resourceNotFoundError *failure.ResourceNotFoundError
	invalidValueError     *failure.InvalidValueError
)

// HandleErrorCase writes an appropriate response body and returns tue if there is an error
func HandleErrorCase(writer http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, resourceNotFoundError):
		WriteErrorResponse(writer, http.StatusNotFound, failure.ResourceNotFoundErrorCode, err.Error())
	case errors.Is(err, invalidValueError):
		WriteErrorResponse(writer, http.StatusBadRequest, failure.InvalidArgumentErrorCode, err.Error())
	default:
		log.Warn().Err(err).Msg("encountered internal server error")
		WriteErrorResponse(writer, http.StatusInternalServerError, failure.InternalErrorErrorCode, err.Error())
	}

	return true
}
