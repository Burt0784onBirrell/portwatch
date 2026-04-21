package groupby

import "errors"

// ErrNilKeyFunc is returned when New is called with a nil key function.
var ErrNilKeyFunc = errors.New("groupby: key function must not be nil")
