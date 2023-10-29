package failure

const (
 // InvalidJSONErrorCode is used when a received body should be json but is syntactically invalid
 InvalidJSONErrorCode = ErrorCode("INVALID_JSON")

 // InvalidArgumentErrorCode is used when a received body, a path param or a query param is semantically invalid
 InvalidArgumentErrorCode = ErrorCode("INVALID_ARGUMENT")

 // InternalErrorErrorCode is used when something went wrong while processing a (valid) request
 InternalErrorErrorCode = ErrorCode("INTERNAL_ERROR")

 // ResourceNotFoundErrorCode is used when a request points to a route that does not exist
 ResourceNotFoundErrorCode = ErrorCode("RESOURCE_NOT_FOUND")

 // MethodNotAllowedErrorCode is used when a route is called with an incorrect http method
 MethodNotAllowedErrorCode = ErrorCode("METHOD_NOT_ALLOWED")

 // TimeoutErrorCode is used when a request request takes too long to read/process/reply to (this can be because of the client or because of the server)
 TimeoutErrorCode = ErrorCode("TIMEOUT")
)
