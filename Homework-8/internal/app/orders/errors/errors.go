package errors

import (
	"errors"
)

// ErrNotFound is an error that means object not found
var ErrNotFound = errors.New("Order not found")
