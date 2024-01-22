package domain

import (
	"errors"
)

var (
	ErrorWhenOpenFile     = errors.New("the application received an error when opening the file")
	ErrorWhenCloseFile    = errors.New("the application received an error when closing the file")
	ErrorWhenReadFile     = errors.New("the application received an error while reading the file")
	ErrorMeanOfDeathParse = errors.New("error during parse mean of death")
)
