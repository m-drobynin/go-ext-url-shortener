package model

import "errors"

var DbConflictError = errors.New("conflict")
var DbNotFoundError = errors.New("not found")

var BadRequestError = errors.New("bad request")
var InternalError = errors.New("internal error")
