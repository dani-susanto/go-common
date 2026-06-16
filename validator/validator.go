package validator

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"strings"
)

type Validator interface {
	Prepare() Validator
	Set(field string, value any) Validator
	Validate() Errors
	Required() Validator
	Email() Validator
	Password() Validator
	PasswordConfirm(password string) Validator
	CountryCode() Validator
	Phone() Validator
	OTP() Validator
	Nullable() Validator
	IN([]string) Validator
}

type validator struct {
	attr     ValidateAttr
	errs     Errors
	nullable bool
}

func New() Validator {
	return &validator{}
}

type ValidateAttr struct {
	Field string
	Value any
}

type Errors map[string][]string

func (v *validator) Prepare() Validator {
	v.errs = make(Errors)
	return v
}

func (v *validator) Set(field string, value any) Validator {
	v.attr = ValidateAttr{
		Field: field,
		Value: value,
	}
	v.nullable = false
	return v
}

func (v *validator) Append(field string, message string) {
	v.errs[field] = append(v.errs[field], message)
}

func (v *validator) Validate() Errors {
	return v.errs
}

func (v *validator) Required() Validator {
	if v.attr.Value == nil || v.attr.Value == "" {
		v.Append(v.attr.Field, fmt.Sprintf("The %v field is required.", v.attr.Field))
	}
	return v
}

func (v *validator) Email() Validator {
	_, err := mail.ParseAddress(fmt.Sprintf("%v", v.attr.Value))
	if err != nil {
		v.Append(v.attr.Field, fmt.Sprintf("The %v field must be a valid email address.", v.attr.Field))
	}
	return v
}

func (v *validator) Password() Validator {
	pass, ok := v.attr.Value.(string)
	if !ok {
		v.Append(v.attr.Field, fmt.Sprintf("The %v field must be a string.", v.attr.Field))
		return v
	}

	if len(pass) < 8 {
		v.Append(v.attr.Field, "Password must be at least 8 characters.")
	}

	if match, _ := regexp.MatchString(`[A-Z]`, pass); !match {
		v.Append(v.attr.Field, "Password must contain at least one uppercase letter.")
	}

	if match, _ := regexp.MatchString(`[a-z]`, pass); !match {
		v.Append(v.attr.Field, "Password must contain at least one lowercase letter.")
	}

	if match, _ := regexp.MatchString(`[0-9]`, pass); !match {
		v.Append(v.attr.Field, "Password must contain at least one number.")
	}

	if match, _ := regexp.MatchString(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>\/?\\|]`, pass); !match {
		v.Append(v.attr.Field, "Password must contain at least one special character.")
	}

	return v
}

func (v *validator) PasswordConfirm(password string) Validator {
	if v.attr.Value != password {
		v.Append(v.attr.Field, fmt.Sprintf("The %v field confirmation does not match.", v.attr.Field))
	}
	return v
}

func (v *validator) Phone() Validator {
	if v.nullable && IsEmpty(v.attr.Value) {
		return v
	}

	phone, ok := v.attr.Value.(string)
	if !ok {
		v.Append(v.attr.Field, fmt.Sprintf("The %v field must be a string.", v.attr.Field))
		return v
	}

	match, _ := regexp.MatchString(`^\+?[0-9]{10,15}$`, phone)
	if !match {
		v.Append(v.attr.Field, fmt.Sprintf("The %v field must be a valid phone number.", v.attr.Field))
	}

	return v
}

func (v *validator) CountryCode() Validator {
	if v.nullable && IsEmpty(v.attr.Value) {
		return v
	}

	code, ok := v.attr.Value.(string)
	if !ok {
		v.Append(v.attr.Field, fmt.Sprintf("The %v field must be a string.", v.attr.Field))
		return v
	}

	match, _ := regexp.MatchString(`^[0-9]{1,3}$`, code)
	if !match {
		v.Append(v.attr.Field, fmt.Sprintf("The %v field must be a valid country code (digits only, 1-3 length).", v.attr.Field))
	}

	return v
}

func (v *validator) OTP() Validator {
	if v.nullable && IsEmpty(v.attr.Value) {
		return v
	}

	otp := v.attr.Value.(string)

	if len(otp) != 6 {
		v.Append(
			v.attr.Field,
			fmt.Sprintf("The %v field must be a 6-digit OTP.", v.attr.Field),
		)
		return v
	}

	for _, c := range otp {
		if c < '0' || c > '9' {
			v.Append(
				v.attr.Field,
				fmt.Sprintf("The %v field must contain only digits.", v.attr.Field),
			)
			return v
		}
	}

	return v
}

func (v *validator) Nullable() Validator {
	v.nullable = true
	return v
}

func (v *validator) IN(matches []string) Validator {
	if v.nullable && IsEmpty(v.attr.Value) {
		return v
	}

	input := strings.TrimSpace(fmt.Sprintf("%v", v.attr.Value))

	match := false
	for _, val := range matches {
		cleanVal := strings.TrimSpace(val)

		if cleanVal == input {
			match = true
			break
		}
	}

	if !match {
		v.Append(
			v.attr.Field,
			fmt.Sprintf(
				"The selected %v is invalid, field must contain one of: [%v].",
				v.attr.Field,
				strings.Join(matches, ", "),
			),
		)
	}

	return v
}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case int:
		return val == 0
	case int8:
		return val == 0
	case int16:
		return val == 0
	case int32:
		return val == 0
	case int64:
		return val == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(val).Uint() == 0
	case float32:
		return val == 0
	case float64:
		return val == 0
	case bool:
		return false
	case *string:
		return val == nil || *val == ""
	case *int:
		return val == nil
	case *int64:
		return val == nil
	case *bool:
		return val == nil
	case *float64:
		return val == nil
	case []byte:
		return len(val) == 0
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Map, reflect.Array:
			return rv.Len() == 0
		case reflect.Ptr, reflect.Interface:
			return rv.IsNil()
		}
		return false
	}
}
