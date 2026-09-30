package model

import "errors"

var ErrDBConflict = errors.New("conflict")
var ErrDBNotFound = errors.New("not found")

var ErrBadRequest = errors.New("bad request")
var ErrInternal = errors.New("internal error")
