package lint020nontypes

import "errors"

var ErrInvalidEnvironment = errors.New("invalid environment") // ok - non-types packages are not checked by LINT-020
