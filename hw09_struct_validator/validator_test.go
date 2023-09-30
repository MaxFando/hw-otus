package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		// Тесты для структуры User
		{
			in: User{
				ID:     "1234567890123456789012345678901234",
				Name:   "John Doe",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{"ID", errors.New("длина поля ID должна быть 36")},
				{"Phones", errors.New("длина поля Phones должна быть 11")},
			},
		},
		{
			in: User{
				ID:     "1234567890",
				Name:   "Alice",
				Age:    17,
				Email:  "alice.com",
				Role:   "guest",
				Phones: []string{"12345"},
			},
			expectedErr: ValidationErrors{
				{"ID", errors.New("длина поля ID должна быть 36")},
				{"Age", errors.New("Age меньше минимального значения 18")},
				{"Email", errors.New("Email не удовлетворяет регулярному выражению")},
				{"Role", errors.New("Role не входит в список разрешенных значений")},
				{"Phones", errors.New("длина поля Phones должна быть 11")},
			},
		},

		// Тесты для структуры App
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0.00",
			},
			expectedErr: ValidationErrors{
				{"Version", errors.New("длина поля Version должна быть 5")},
			},
		},

		// Тесты для структуры Token
		{
			in: Token{},
			expectedErr: ValidationErrors{
				{"Header", errors.New("неверное поле для правила len: Header")},
				{"Payload", errors.New("неверное поле для правила len: Payload")},
				{"Signature", errors.New("неверное поле для правила len: Signature")},
			},
		},

		// Тесты для структуры Response
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "Not Found",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 500,
				Body: "",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 400,
				Body: "Bad Request",
			},
			expectedErr: ValidationErrors{
				{"Code", errors.New("Code не входит в список разрешенных значений")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("Ожидается успешная валидация, получена ошибка: %v", err)
				}
			} else {
				if err != nil && err.Error() != tt.expectedErr.Error() {
					t.Errorf("Ожидается ошибка: %v, получена ошибка: %v", tt.expectedErr, err)
				}
			}
		})
	}
}
