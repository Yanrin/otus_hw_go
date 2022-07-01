package hw09structvalidator

import (
	"errors"
	"strings"
)

var (
	ErrExpectedStruct         = errors.New("expected struct")
	ErrRuleIncorrect          = errors.New("incorrect rule")
	ErrRuleValueIncorrect     = errors.New("incorrect rule value")
	ErrRuleUnsupported        = errors.New("unsupported rule")
	ErrValidationOccurrence   = errors.New("not equal with any occurrence")
	ErrValidationStringLength = errors.New("length differs from value in tag")
	ErrValidationStringRegexp = errors.New("does not match the mask")
	ErrValidationIntMax       = errors.New("greater than max required")
	ErrValidationIntMin       = errors.New("less than min required")
)

type ValidationError struct {
	Field string
	Err   error
}

// Error returns the describing string of the error.
func (v ValidationError) Error() string {
	var b strings.Builder
	b.WriteString(v.Field)
	b.WriteString(": ")
	b.WriteString(v.Err.Error())

	return b.String()
}

type ValidationErrors []ValidationError

// Error returns the accumulated string of all errors in the slice.
func (v ValidationErrors) Error() string {
	var b strings.Builder
	for _, e := range v {
		b.WriteString(e.Error())
		b.WriteString("\n")
	}

	return b.String()
}

// Add adds new error into the slice.
func (v *ValidationErrors) Add(field string, err error) {
	*v = append(*v, ValidationError{Field: field, Err: err})
}

// AddList adds the list of the existing errors into the slice.
func (v *ValidationErrors) AddList(list ValidationErrors) {
	*v = append(*v, list...)
}
