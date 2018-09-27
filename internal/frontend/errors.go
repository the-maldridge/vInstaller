package frontend

import (
	"errors"
)

var (
	// ErrUnknownFrontend is to be returned when a frontend has
	// been requested that is not known to the system.
	ErrUnknownFrontend = errors.New("frontend name is not known")

	// ErrConfigUnobtainable is to be returned if the
	// configuration cannot be obtained due to some terminal
	// failure.
	ErrConfigUnobtainable = errors.New("config cannot be obtained")

	// ErrInstallationAborted is to be returned by the
	// ConfirmInstallation() function if the user chooses not to
	// install the system.
	ErrInstallationAborted = errors.New("the installation was cancelled")
)
