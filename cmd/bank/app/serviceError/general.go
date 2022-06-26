package serviceError

import "errors"

var ErrNoIdentifier = errors.New("no identifier")
var ErrNoAuthentication = errors.New("no authentication")
var ErrContextEmpty = errors.New("Context is empty!")
