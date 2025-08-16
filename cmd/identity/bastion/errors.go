package bastion

import "errors"

// ErrAborted is returned when the user cancels a TUI or selection.
var ErrAborted = errors.New("aborted by user")
