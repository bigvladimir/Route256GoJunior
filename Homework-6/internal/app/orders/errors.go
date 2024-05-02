package orders

import (
	"errors"
)

// ErrNotFound is an error that means object not found
var ErrNotFound = errors.New("Не найдено в базе")
