package jsonform

import "errors"

var (
	errUnexpectedInput    = errors.New("unexpected input")
	errUnexpectedOutput   = errors.New("unexpected output")
	errSchemaAlreadyAdded = errors.New("schema is already added")
	errMissingFormSchema  = errors.New("missing form schema")
)
