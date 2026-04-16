package sequence

import "errors"

// ErrEmptyField is returned when New is called with an empty field name.
var ErrEmptyField = errors.New("sequence: field name must not be empty")
