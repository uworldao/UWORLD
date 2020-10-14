package core

import "errors"

var (
	ErrDuplicateBlock = errors.New("duplicate block")
	ErrNoParent       = errors.New("not find block parent header")
)
