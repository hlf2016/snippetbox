package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) 
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// 返回的 ID 类型为 int64，因此我们在返回前将其转换为 int 类型。
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// 初始化指向已清零的新 Snippet 结构的指针。
	s := &Snippet{}
	// 使用 row.Scan() 将 sql.Row 中每个字段的值复制到 Snippet 结构中的相应字段。
	// 请注意，row.Scan 的参数是指向要将数据复制到的位置的指针，参数数必须与语句返回的列数完全相同。
	err := m.DB.QueryRow(`SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// 如果查询没有返回记录，那么 row.Scan() 将返回一个 sql.ErrNoRows 错误。
		// 我们使用 errors.Is() 函数专门检查该错误，并返回我们自己的 ErrNoRecord 错误（我们稍后将创建该错误）
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	rows, err := m.DB.Query(`SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}
	// 使用 defer rows.Close() 关闭结果集至关重要。只要结果集处于打开状态，底层数据库连接就会一直处于打开状态......因此，
	// 如果该方法出现问题，结果集没有关闭，就会迅速导致池中的所有连接被用完
	defer rows.Close()

	var snippets []*Snippet

	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
