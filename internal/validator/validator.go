package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX 使用 regexp.MustCompile() 函数解析正则表达式模式，以检查电子邮件地址的格式是否正确。该函数会返回一个指向 "已编译 "regexp.Regexp 类型的指针，如果出现错误则会宕机。在启动时解析一次该模式并将编译后的 regexp.Regexp 保存在一个变量中，比每次需要时都重新解析该模式更高效
// https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator 定义一个新的验证器类型，其中包含表单字段的验证错误映射。
type Validator struct {
	// 在结构体中添加一个新的 NonFieldErrors []string 字段，用于保存与特定表单字段无关的验证错误。
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid 如果 validator 的 FieldErrors 不包含任何项时 返回 true
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func PermittedInt(value int, permittedValue ...int) bool {
	for i := range permittedValue {
		if value == permittedValue[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
