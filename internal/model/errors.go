package model

import "errors"

var ErrDBConflict = errors.New("conflict")
var ErrDBNotFound = errors.New("not found")
