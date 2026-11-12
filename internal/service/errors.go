package service

import (
	"errors"
)

var ErrSaveURLMaxAttempts = errors.New("max url save retries reached")
var ErrUniqueValue = errors.New("non unique value")
