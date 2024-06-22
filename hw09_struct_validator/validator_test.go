package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code   int    `validate:"in:200,404,500"`
		Values []int  `validate:"min:11|max:22"`
		Body   string `json:"omitempty"`
	}

	WrongStringTag struct {
		Name string `validate:"equals:1"`
	}

	WrongIntMin struct {
		Tag int `validate:"min:abc"`
	}

	WrongIntArrTag struct {
		Tag []int `validate:"equals:1"`
	}

	WrongRegexp struct {
		Name string `validate:"regexp:[a-z"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			"not_a_struct",
			ErrNotStruct,
		},
		{
			User{
				ID:    "12345", // len != 36
				Name:  "Test",
				Age:   60,             // > max=50
				Email: "wrongmail.ru", // not correspond regexp
				Role:  "tester",       // not in admin,stuff
				Phones: []string{
					"123456789",
					"1234567890",
				}, // len != 11
				meta: nil,
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrStringIncorrectLen,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrIntMax,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrStringNotMatchRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrStringNotInValues,
				},
				ValidationError{
					Field: "Phones index=0",
					Err:   ErrStringIncorrectLen,
				},
				ValidationError{
					Field: "Phones index=1",
					Err:   ErrStringIncorrectLen,
				},
			},
		},
		{
			User{
				ID:    "123456789012345678901234567890123456",
				Name:  "Test",
				Age:   40,
				Email: "test@mail.ru",
				Role:  "admin",
				Phones: []string{
					"12345678901",
					"12345678901",
				},
				meta: nil,
			},
			nil,
		},
		{
			Token{
				Header:    []byte("Header"),
				Payload:   nil,
				Signature: []byte{},
			},
			nil,
		},
		{
			App{
				Version: "version12345", // len != 5
			},
			ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrStringIncorrectLen,
				},
			},
		},
		{
			Response{
				Code:   888,               // not in (200,404,500)
				Values: []int{10, 11, 22}, // >= 11 && <= 22
				Body:   "abc",
			},
			ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrIntNotInValues,
				},
				ValidationError{
					Field: "Values index=0",
					Err:   ErrIntMin,
				},
			},
		},
		{
			WrongStringTag{
				Name: "name",
			},
			ErrStringInvalidRule,
		},
		{
			WrongIntMin{
				Tag: 1,
			},
			ErrIntInvalidRule,
		},
		{
			WrongIntArrTag{
				Tag: []int{1, 2, 3},
			},
			ErrIntInvalidRule,
		},
		{
			WrongRegexp{
				Name: "name",
			},
			ErrStringInvalidRule,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				var actualValidationErrors ValidationErrors
				if errors.As(err, &actualValidationErrors) {
					var expectedValidationErrors ValidationErrors
					require.ErrorAs(t, tt.expectedErr, &expectedValidationErrors)
					require.Equal(t, len(expectedValidationErrors), len(actualValidationErrors))
					for j, actualErr := range actualValidationErrors {
						require.ErrorIs(t, actualErr, expectedValidationErrors[j])
					}
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}
			}
			_ = tt
		})
	}
}
