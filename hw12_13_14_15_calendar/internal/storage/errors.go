package storage

import "errors"

var (
	ErrEventNotFound = errors.New("event not found")
	ErrBusyTime      = errors.New("time is busy by another event")
)
