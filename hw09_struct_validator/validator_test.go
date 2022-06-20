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
		Version string `validate:"regexp:^[\\d\\.]*$|len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	InvalidTags struct {
		WorldWidth       string `validate:"len:entireworld"`
		EmptyTag         string `validate:"in:"`
		EmptyAlias       string `validate:""`
		NonValidationTag string `name:"phantasmagoria"`
		UnsupportedTag   string `validate:"firebird:654"`
	}
)

var errVT ValidationErrors

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App{Version: "12.15"},
			expectedErr: nil,
		},
		{
			in: App{Version: "2.3a"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrValidationStringRegexp},
				ValidationError{Field: "Version", Err: ErrValidationStringLength},
			},
		},
		{
			in: InvalidTags{
				WorldWidth:       "345",
				EmptyTag:         "no tag",
				EmptyAlias:       "nothing",
				NonValidationTag: "nuts",
				UnsupportedTag:   "deer",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "WorldWidth", Err: ErrRuleValueIncorrect},
				ValidationError{Field: "EmptyTag", Err: ErrRuleIncorrect},
				ValidationError{Field: "EmptyAlias", Err: ErrRuleIncorrect},
				ValidationError{Field: "UnsupportedTag", Err: ErrRuleUnsupported},
			},
		},
		{
			in: User{
				ID:     "58ca9f62-0000-0000-0000-0000a2c4cb2c",
				Name:   "Jaime Lannister",
				Age:    33,
				Email:  "JL@web.com",
				Role:   "admin",
				Phones: []string{"89991234567"},
				meta:   json.RawMessage(`{"message": "anything"}`),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "58ca9f62-0000-0000-0000-00002c4cb2c",
				Name:   "Cersei Lannister",
				Age:    16,
				Email:  "kings@landing",
				Role:   "queen",
				Phones: []string{"8999000"},
				meta:   json.RawMessage(`{"message": "mhm"}`),
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrValidationStringLength},
				ValidationError{Field: "Age", Err: ErrValidationIntMin},
				ValidationError{Field: "Email", Err: ErrValidationStringRegexp},
				ValidationError{Field: "Role", Err: ErrValidationOccurrence},
				ValidationError{Field: "Phones", Err: ErrValidationStringLength},
			},
		},
		{
			in:          Response{Code: 404},
			expectedErr: nil,
		},
		{
			in: Response{Code: 401, Body: "Who are you?"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrValidationOccurrence},
			},
		},
		{
			in:          "Every family has a sheep",
			expectedErr: ErrExpectedStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				if errors.As(err, &errVT) {
					require.Equal(t, err.Error(), tt.expectedErr.Error())
				} else {
					require.ErrorIs(t, err, ErrExpectedStruct)
				}
			}
		})
	}
}
