package main

import (
	"errors"
	"fmt"
)

func formatError(label string, err error) error {
	errorString := fmt.Sprintf("%s: %s", label, err)
	return errors.New(errorString)
}
