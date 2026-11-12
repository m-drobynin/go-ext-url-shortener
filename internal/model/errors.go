package model

import "errors"

var ErrUrlsConflict = errors.New("urls conflict")
var ErrCodeConflict = errors.New("codes conflict")
var ErrDBNotFound = errors.New("not found")
