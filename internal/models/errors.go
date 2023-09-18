package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials 添加新的 ErrInvalidCredentials 错误。如果用户尝试使用错误的电子邮件地址或密码登录
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail 添加新的 ErrDuplicateEmail 错误。如果用户尝试使用已被使用的电子邮件地址注册
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
