package config

import "errors"

// ErrInvalidInterval is returned when scan_interval is less than 1 second.
var ErrInvalidInterval = errors.New("config: scan_interval must be at least 1s")

// ErrNotFound is returned when the config file does not exist.
var ErrNotFound = errors.New("config: file not found")
