package common

import "errors"

// TODO: use same pattern as in pb-core here

// ErrNotFound is used when an element could not be found
var ErrNotFound = errors.New("Element not found")

// ErrInvalidID is used when an invalid id was received
var ErrInvalidID = errors.New("Invalid ID")

// ErrInvalidOperation is used when a constraint violation occured
var ErrInvalidOperation = errors.New("Invalid operation")
