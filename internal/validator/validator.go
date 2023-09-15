package validator

import (
	"strings"
	"unicode/utf8"
)

// Validator 定义一个新的验证器类型，其中包含表单字段的验证错误映射。
type Validator struct {
	FieldErrors map[string]string
}

// Valid 如果 validator 的 FieldErrors 不包含任何项时 返回 true
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValue ...int) bool {
	for i := range permittedValue {
		if value == permittedValue[i] {
			return true
		}
	}
	return false
}
