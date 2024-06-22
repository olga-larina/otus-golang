package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrNotStruct = errors.New("not a struct")

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, len(v))
	for i, e := range v {
		errs[i] = e.Error()
	}
	return strings.Join(errs, "\n")
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field: %s, error: %s", v.Field, v.Err.Error())
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("not valid interface type %v: error %w", val.Kind().String(), ErrNotStruct)
	}

	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		fieldVal := val.Field(i)
		fieldName := field.Name
		err := validateField(fieldName, fieldVal, tag)

		if err = processError(&validationErrors, err); err != nil {
			return err
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(fieldName string, value reflect.Value, tag string) error {
	var validationErrors ValidationErrors
	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		switch {
		case value.Kind() == reflect.String:
			err := validateString(fieldName, value.String(), rule)
			if err = processError(&validationErrors, err); err != nil {
				return err
			}
		case value.Kind() == reflect.Int:
			err := validateInt(fieldName, value.Int(), rule)
			if err = processError(&validationErrors, err); err != nil {
				return err
			}
		case value.Kind() == reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				elem := value.Index(i)
				switch {
				case elem.Kind() == reflect.String:
					err := validateString(fmt.Sprintf("%s index=%d", fieldName, i), elem.String(), rule)
					if err = processError(&validationErrors, err); err != nil {
						return err
					}
				case elem.Kind() == reflect.Int:
					err := validateInt(fmt.Sprintf("%s index=%d", fieldName, i), elem.Int(), rule)
					if err = processError(&validationErrors, err); err != nil {
						return err
					}
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

// Проверка типа ошибки и её обработка
// Если ошибка err имеет тип ValidationError или ValidationErrors,
//
//	то она добавляется в validationErrors, а возвращаемое значение - null
//
// Иначе, возвращаемое значение == ошибке err (т.е. это программная ошибка).
func processError(validationErrors *ValidationErrors, err error) error {
	if err == nil {
		return nil
	}

	var fieldValidationError ValidationError
	if errors.As(err, &fieldValidationError) { // ошибка валидации
		*validationErrors = append(*validationErrors, fieldValidationError)
		return nil
	}

	var fieldValidationErrors ValidationErrors
	if errors.As(err, &fieldValidationErrors) { // ошибки валидации
		*validationErrors = append(*validationErrors, fieldValidationErrors...)
		return nil
	}

	// программная ошибка
	return err
}
