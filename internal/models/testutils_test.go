package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// 为我们的测试数据库建立 sql.DB 连接池。由于我们的设置和拆卸脚本包含多条 SQL 语句，因此我们需要在 DSN 中使用 "multiStatements=true "参数。
	// 这将指示我们的 MySQL 数据库驱动程序支持在一次 db.Exec() 调用中执行多条 SQL 语句。
	db, err := sql.Open("mysql", "test_web:25804769@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
	// 使用 t.Cleanup() 注册一个函数，当调用 newTestDB() 的当前测试（或子测试）结束时，Go 会自动调用该函数*。
	// 在该函数中，我们将读取并执行拆除脚本，并关闭数据库连接池。
	t.Cleanup(func() {
		script, err = os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
	return db
}
