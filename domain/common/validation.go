package common

import (
	"authenticator-backend/extension/logger"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

const (
	errorMessageCannotBeBlank         = "cannot be blank"
	errorMessageInvalidUUID           = "invalid UUID"
	errorMessageLast6DigitsNotNumeric = "the last 6 digits must always be numeric"
)

// StringUUIDValid
// Summary: This is function which checks whether the string is a valid UUID
// input: value (interface{})
// output: (error) error object
func StringUUIDValid(value interface{}) error {
	s, _ := value.(string)
	if s == "" {
		logger.Set(nil).Warnf(errorMessageCannotBeBlank)

		return fmt.Errorf(errorMessageCannotBeBlank)
	} else {
		if len(s) != UUIDLen {
			logger.Set(nil).Warnf(errorMessageInvalidUUID)

			return fmt.Errorf(errorMessageInvalidUUID)
		}

		_, err := uuid.Parse(s)
		if err != nil {
			logger.Set(nil).Warnf(errorMessageInvalidUUID)

			return fmt.Errorf(errorMessageInvalidUUID)
		}
	}

	return nil
}

// StringPtrNilOrUUIDValid
// Summary: This is function which checks whether the string pointer is nil or a valid UUID
// input: value (interface{})
// output: (error) error object
func StringPtrNilOrUUIDValid(value interface{}) error {
	sp, _ := value.(*string)
	if sp == nil {
		return nil
	}
	s := *sp

	if len(s) != UUIDLen {
		logger.Set(nil).Warnf(errorMessageInvalidUUID)

		return fmt.Errorf(errorMessageInvalidUUID)
	}

	_, err := uuid.Parse(s)
	if err != nil {
		logger.Set(nil).Warnf(errorMessageInvalidUUID)

		return fmt.Errorf(errorMessageInvalidUUID)
	}

	return nil
}

// StringPtrLast6CharsNumeric
// Summary: This is function which checks whether the last 6 characters of the string pointer are numeric
// input: value (interface{})
// output: (error) error object
func StringPtrLast6CharsNumeric(value interface{}) error {
	sp, _ := value.(*string)
	if sp == nil {
		return nil
	}
	s := *sp

	if len(s) < 6 {
		logger.Set(nil).Warnf(errorMessageLast6DigitsNotNumeric)

		return fmt.Errorf(errorMessageLast6DigitsNotNumeric)
	}

	last6 := s[len(s)-6:]
	_, err := strconv.Atoi(last6)
	if err != nil {
		logger.Set(nil).Warnf(errorMessageLast6DigitsNotNumeric)

		return fmt.Errorf(errorMessageLast6DigitsNotNumeric)
	}

	return nil
}

// JoinErrors
// Summary: This is function which joins multiple error messages
// input: errors ([]error) error object
// output: (error) error object
func JoinErrors(errors []error) error {
	errDetails := ""
	errorMap := make(map[string]bool)

	filteredErrors := []error{}
	for _, v := range errors {
		if _, ok := errorMap[v.Error()]; ok {
			continue
		}
		errorMap[v.Error()] = true
		filteredErrors = append(filteredErrors, v)
	}

	for i, v := range filteredErrors {
		errDetails += v.Error()

		if i+1 < len(filteredErrors) {
			errDetails += "; "
		} else {
			errDetails += "."
		}
	}

	return fmt.Errorf(errDetails)
}
