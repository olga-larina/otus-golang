package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type CurrentState struct {
	resultStr     strings.Builder // билдер для сбора строки
	prevRuneValue rune            // предыдущая руна
	isHanging     bool            // есть ли "висящая" руна (последняя), которая не была добавлена в строку
	isBackslashed bool            // экранирован ли символ
}

func Unpack(str string) (string, error) {
	state := &CurrentState{}
	for i, runeValue := range str {
		switch {
		case unicode.IsDigit(runeValue):
			if err := state.processDigit(i, runeValue); err != nil {
				return "", makeError(err, i, runeValue)
			}
		case isBackslash(runeValue):
			state.processBackslash()
		default:
			if err := state.processRune(); err != nil {
				return "", makeError(err, i, runeValue)
			}
		}
		state.prevRuneValue = runeValue
	}
	if state.isBackslashed {
		return "", fmt.Errorf("backslash can be before digit or backslash: %w", ErrInvalidString)
	}
	if state.isHanging {
		state.resultStr.WriteRune(state.prevRuneValue)
	}
	return state.resultStr.String(), nil
}

func (state *CurrentState) processDigit(curIndex int, curRuneValue rune) error {
	if curIndex == 0 || (!state.isHanging && unicode.IsDigit(state.prevRuneValue)) {
		return fmt.Errorf("digit should be after backslash, rune or backslashed digit: %w", ErrInvalidString)
	}
	if state.isBackslashed {
		state.isHanging = true
		state.isBackslashed = false
	} else {
		cntRunes := toInt(curRuneValue)
		state.repeat(state.prevRuneValue, cntRunes)
		state.isHanging = false
	}
	return nil
}

func (state *CurrentState) processBackslash() {
	if state.isHanging {
		state.resultStr.WriteRune(state.prevRuneValue)
	}
	if state.isBackslashed {
		state.isHanging = true
		state.isBackslashed = false
	} else {
		state.isHanging = false
		state.isBackslashed = true
	}
}

func (state *CurrentState) processRune() error {
	if state.isBackslashed {
		return fmt.Errorf("backslash should be before backslash or digit: %w", ErrInvalidString)
	}
	if state.isHanging {
		state.resultStr.WriteRune(state.prevRuneValue)
	}
	state.isHanging = true
	return nil
}

func (state *CurrentState) repeat(r rune, times int) {
	for i := 0; i < times; i++ {
		state.resultStr.WriteRune(r)
	}
}

func toInt(r rune) int {
	return int(r - '0')
}

func isBackslash(r rune) bool {
	return r == '\\'
}

func makeError(err error, curIndex int, curRuneValue rune) error {
	return fmt.Errorf("processing rune %c at index %d: %w", curRuneValue, curIndex, err)
}
