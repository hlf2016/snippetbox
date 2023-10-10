package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES (? ,?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// 如果返回错误，我们将使用 errors.As() 函数检查错误是否属于 mysql.MySQLError 类型。
		// 如果是，该错误将被赋值给 mySQLError 变量。然后，我们可以通过检查错误代码是否等于 1062 以及错误消息字符串的内容，检查错误是否与 users_uc_email 密钥有关。如果是，我们将返回 ErrDuplicateEmail 错误信息
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			if mysqlError.Number == 1062 && strings.Contains(mysqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	// 读取与给定电子邮件相关的 ID 和哈希密码。如果不存在匹配的电子邮件，我们将返回 ErrInvalidCredentials 错误信息
	var id int
	var hashedPassword []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email= ?`
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredential
		} else {
			return 0, err
		}
	}
	// 检查提供的散列密码和纯文本密码是否匹配。如果不匹配，我们将返回 ErrInvalidCredentials 错误。
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredential
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id=?)`
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) Get(id int) (*User, error) {
	var user User
	stmt := "SELECT id, name, email, created from users WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return &user, nil
}
