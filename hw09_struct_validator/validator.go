//nolint:lll
package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	messages := make([]string, 0, len(v))
	for _, ve := range v {
		messages = append(messages, ve.Err.Error())
	}

	return strings.Join(messages, "; ")
}

// Validate валидирует структуру v на основе структурных тегов "validate".
//
//nolint:gocognit,gocyclo,funlen
func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return errors.New("ожидается структура")
	}

	var validationErrors []ValidationError

	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typ.Field(i)

		// Получаем значение тэга validate
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue // Пропускаем поля без тега validate
		}

		fieldName := fieldType.Name

		// Разделяем тэги по символу "|"
		validationRules := strings.Split(validateTag, "|")

		for _, rule := range validationRules {
			ruleParts := strings.Split(rule, ":")
			if len(ruleParts) < 2 {
				return fmt.Errorf("неверный формат правила валидации для поля %s", fieldName)
			}

			ruleType := ruleParts[0]
			ruleValue := strings.Join(ruleParts[1:], ":")

			switch ruleType {
			case "len":
				length, err := strconv.Atoi(ruleValue)
				if err != nil {
					return fmt.Errorf("неверное значение правила len для поля %s: %w", fieldName, err)
				}

				//nolint:exhaustive
				switch field.Kind() {
				case reflect.String:
					if field.Len() != length {
						validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("длина поля %s должна быть %d", fieldName, length)})
					}
				case reflect.Slice:
					for j := 0; j < field.Len(); j++ {
						sliceElem := field.Index(j)
						if sliceElem.Len() != length {
							validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("длина элемента %s в списке %s должна быть %d", sliceElem.String(), fieldName, length)})
						}
					}
				default:
					return fmt.Errorf("неверное поле для правила len: %s", fieldName)
				}

			case "regexp":
				//nolint:exhaustive
				switch field.Kind() {
				case reflect.Slice:
					for j := 0; j < field.Len(); j++ {
						sliceElem := field.Index(j)
						if sliceElem.Kind() != reflect.String {
							return fmt.Errorf("неверное поле для правила regexp: %s", fieldName)
						}
						match, err := regexp.MatchString(ruleValue, field.String())
						if err != nil {
							return err
						}
						if !match {
							validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("%s не удовлетворяет регулярному выражению", fieldName)})
						}
					}
				default:
					if field.Kind() != reflect.String {
						return fmt.Errorf("неверное поле для правила regexp: %s", fieldName)
					}
					match, err := regexp.MatchString(ruleValue, field.String())
					if err != nil {
						return fmt.Errorf("ошибка при проверке регулярного выражения для поля %s: %w", fieldName, err)
					}
					if !match {
						validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("%s не удовлетворяет регулярному выражению", fieldName)})
					}
				}
			case "in":
				//nolint:exhaustive
				switch field.Kind() {
				case reflect.Slice:
					for j := 0; j < field.Len(); j++ {
						sliceElem := field.Index(j)
						inValues := strings.Split(ruleValue, ",")

						valid, err := validateIn(sliceElem, fieldName, inValues)
						if err != nil {
							return err
						}

						if !valid {
							validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("%s не входит в список разрешенных значений", fieldName)})
						}
					}
				default:
					inValues := strings.Split(ruleValue, ",")
					valid, err := validateIn(field, fieldName, inValues)
					if err != nil {
						return err
					}

					if !valid {
						validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: fmt.Errorf("%s не входит в список разрешенных значений", fieldName)})
					}
				}

			case "min":
				minValue, err := strconv.Atoi(ruleValue)
				if err != nil {
					return fmt.Errorf("неверное значение для правила min: %s", fieldName)
				}

				//nolint:exhaustive
				switch field.Kind() {
				case reflect.Int:
					errMin := validateMin(fieldName, int(field.Int()), minValue)
					if errMin != nil {
						validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: errMin})
					}
				case reflect.Slice:
					for j := 0; j < field.Len(); j++ {
						sliceElem := field.Index(j)
						errMin := validateMin(fieldName, int(sliceElem.Int()), minValue)
						if errMin != nil {
							validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: errMin})
						}
					}
				default:
					return fmt.Errorf("неверное поле для правила min: %s", fieldName)
				}
			case "max":
				maxValue, err := strconv.Atoi(ruleValue)
				if err != nil {
					return fmt.Errorf("неверное значение для правила max: %s", fieldName)
				}

				//nolint:exhaustive
				switch field.Kind() {
				case reflect.Slice:
					for j := 0; j < field.Len(); j++ {
						sliceElem := field.Index(j)
						errMax := validateMax(fieldName, int(sliceElem.Int()), maxValue)
						if errMax != nil {
							validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: errMax})
						}
					}
				case reflect.Int:
					errMax := validateMax(fieldName, int(field.Int()), maxValue)
					if errMax != nil {
						validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: errMax})
					}
				default:
					return fmt.Errorf("неверное поле для правила max: %s", fieldName)
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return ValidationErrors(validationErrors)
	}

	return nil
}

func validateMax(fieldName string, v, maxValue int) error {
	if v > maxValue {
		return fmt.Errorf("%s больше максимального значения %d", fieldName, maxValue)
	}
	return nil
}

func validateMin(fieldName string, v, minValue int) error {
	if v < minValue {
		return fmt.Errorf("%s меньше минимального значения %d", fieldName, minValue)
	}
	return nil
}

func validateIn(field reflect.Value, fieldName string, inValues []string) (bool, error) {
	valid := false

	for _, inValue := range inValues {
		inValue = strings.TrimSpace(inValue)
		if field.Kind() == reflect.String && field.String() == inValue {
			valid = true
			break
		}
		if field.Kind() == reflect.Int {
			intValue, err := strconv.Atoi(inValue)
			if err != nil {
				return false, fmt.Errorf("неверное значение для правила in: %s", fieldName)
			}
			if int(field.Int()) == intValue {
				valid = true
				break
			}
		}
	}

	return valid, nil
}
