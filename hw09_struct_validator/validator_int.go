package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reMin    = regexp.MustCompile(`^min:(\d+)$`)
	reMax    = regexp.MustCompile(`^max:(.+)$`)
	reInNums = regexp.MustCompile(`^in:(.+)$`)
)

var (
	ErrIntMin         = errors.New("the number is less than min")
	ErrIntMax         = errors.New("the number is more than max")
	ErrIntNotInValues = errors.New("the number is not in the values")
	ErrIntInvalidRule = errors.New("invalid int rule")
)

func validateInt(fieldName string, value int64, rule string) error {
	if len(rule) < 1 {
		return nil
	}

	minRuleMatch := reMin.FindStringSubmatch(rule)
	if len(minRuleMatch) > 1 {
		number, err := strconv.Atoi(minRuleMatch[1])
		if err != nil {
			return fmt.Errorf("%w: failed converting min rule %s for %s: %w", ErrIntInvalidRule, rule, fieldName, err)
		}
		if value < int64(number) {
			return ValidationError{Field: fieldName, Err: ErrIntMin}
		}
		return nil
	}

	maxRuleMatch := reMax.FindStringSubmatch(rule)
	if len(maxRuleMatch) > 1 {
		number, err := strconv.Atoi(maxRuleMatch[1])
		if err != nil {
			return fmt.Errorf("%w: failed converting max rule %s for %s: %w", ErrIntInvalidRule, rule, fieldName, err)
		}
		if value > int64(number) {
			return ValidationError{Field: fieldName, Err: ErrIntMax}
		}
		return nil
	}

	inNumsRuleMatch := reInNums.FindStringSubmatch(rule)
	if len(inNumsRuleMatch) > 1 {
		possibleValues := strings.Split(inNumsRuleMatch[1], ",")
		matched := false
		for _, possibleValue := range possibleValues {
			possibleNumber, err := strconv.Atoi(possibleValue)
			if err != nil {
				return fmt.Errorf("%w: failed converting in rule %s for %s: %w", ErrIntInvalidRule, rule, fieldName, err)
			}
			if value == int64(possibleNumber) {
				matched = true
				break
			}
		}
		if !matched {
			return ValidationError{Field: fieldName, Err: ErrIntNotInValues}
		}
		return nil
	}

	return ErrIntInvalidRule
}
