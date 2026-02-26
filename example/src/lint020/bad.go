package types

import "errors"

var ErrNotFound = errors.New("not found")        // want `LINT-020: error variable "ErrNotFound" must be defined in errors\.go`
var ErrUnauthorized = errors.New("unauthorized") // want `LINT-020: error variable "ErrUnauthorized" must be defined in errors\.go`

var regularVar = "ok" // ok — does not start with Err
