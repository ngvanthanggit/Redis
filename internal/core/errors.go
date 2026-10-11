package core

import "errors"

var ErrInvalidData error = errors.New("invalid data")
var ErrNoData error = errors.New("no data")
var ErrNumberOfArgs error = errors.New("ERR wrong number of arguments for 'ping' command")
