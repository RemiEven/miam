package model

// CodeError is an error with a code and a message
type CodeError struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}
