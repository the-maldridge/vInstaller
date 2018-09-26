package frontend

import (
	"errors"
)

var (
	// ErrUnknownFrontend is to be returned when a frontend has
	// been requested that is not known to the system.
	ErrUnknownFrontend = errors.New("frontend name is not known")
)
