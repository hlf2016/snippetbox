package models

import (
	"github.com/hlf2016/snippetbox/internal/assert"
	"testing"
)

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用 newTestDB() 辅助函数获取测试数据库的连接池。在 t.Run() 中调用此函数意味着将为每个子测试设置和删除新的数据库表和数据。
			db := newTestDB(t)
			m := UserModel{db}
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
