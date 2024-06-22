package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reLen  = regexp.MustCompile(`^len:(\d+)$`)
	reExpr = regexp.MustCompile(`^regexp:(.+)$`)
	reIn   = regexp.MustCompile(`^in:(.+)$`)
)

var (
	ErrStringIncorrectLen   = errors.New("incorrect len of the string")
	ErrStringNotMatchRegexp = errors.New("the string does not match regexp")
	ErrStringNotInValues    = errors.New("the string is not in the values")
	ErrStringInvalidRule    = errors.New("invalid string rule")
)

func validateString(fieldName string, value string, rule string) error {
	if len(rule) < 1 {
		return nil
	}

	lenRuleMatch := reLen.FindStringSubmatch(rule)
	if len(lenRuleMatch) > 1 {
		number, err := strconv.Atoi(lenRuleMatch[1])
		if err != nil {
			return fmt.Errorf("%w: failed converting len rule %s for %s: %w", ErrStringInvalidRule, rule, fieldName, err)
		}
		if len(value) != number {
			return ValidationError{Field: fieldName, Err: ErrStringIncorrectLen}
		}
		return nil
	}

	exprRuleMatch := reExpr.FindStringSubmatch(rule)
	if len(exprRuleMatch) > 1 {
		expr, err := regexp.Compile(exprRuleMatch[1])
		if err != nil {
			return fmt.Errorf("%w: failed converting regexp rule %s for %s: %w", ErrStringInvalidRule, rule, fieldName, err)
		}
		if !expr.MatchString(value) {
			return ValidationError{Field: fieldName, Err: ErrStringNotMatchRegexp}
		}
		return nil
	}

	inRuleMatch := reIn.FindStringSubmatch(rule)
	if len(inRuleMatch) > 1 {
		possibleValues := strings.Split(inRuleMatch[1], ",")
		matched := false
		for _, possibleValue := range possibleValues {
			if value == possibleValue {
				matched = true
				break
			}
		}
		if !matched {
			return ValidationError{Field: fieldName, Err: ErrStringNotInValues}
		}
		return nil
	}

	return ErrStringInvalidRule
}
