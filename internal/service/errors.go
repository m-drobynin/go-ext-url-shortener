package service

import (
	"errors"
)

var ErrSaveUrlMaxAttempts = errors.New("max url save retries reached")
