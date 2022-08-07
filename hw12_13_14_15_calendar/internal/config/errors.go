package config

import (
	"errors"
)

var (
	ErrFilePathEmpty = errors.New("file path is empty")
	ErrOpenFailed    = errors.New("opening config file is failed")
	ErrReadFile      = errors.New("can't read file")
)
