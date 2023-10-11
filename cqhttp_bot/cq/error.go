package cq

import "errors"

var (
	ErrorNotCQCode   = errors.New("error not cq code")
	ErrorUnknownCode = errors.New("error unknown code")
)
