package application

import "errors"

var (
	ErrDuplicateLogPath = errors.New("log path already exists")
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidPath      = errors.New("invalid path")
	ErrArchiveNotFound  = errors.New("archive not found")
	ErrTopologyNotReady = errors.New("topology is only available for successfully parsed logs")
	ErrParseFailed      = errors.New("parse failed")
	ErrPersistFailed    = errors.New("persist failed")
)
