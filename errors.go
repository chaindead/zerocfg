package zfg

import "errors"

var (
	ErrNoSuchKey      = errors.New("no such key")
	ErrCollidingAlias = errors.New("colliding alias with key")
	ErrDuplicateKey   = errors.New("duplicate key")
)
