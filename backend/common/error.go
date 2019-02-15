package common

import "errors"

// ErrNotFound is used when an element could not be found
var ErrNotFound = errors.New("Element not found")

// ErrInvalidID is used when an invalid id was received
var ErrInvalidID = errors.New("Invalid ID")
