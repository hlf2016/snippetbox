package main

import (
	"github.com/hlf2016/snippetbox/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		// 在第三个测试用例中，我们使用 CET（中欧时间）作为时区，它比 UTC 早一个小时。因此，我们希望 humanDate() 的输出结果（UTC 时区）是 2022 年 3 月 17 日 09:15 时，而不是 2022 年 3 月 17 日 10:15 时。
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	for _, tt := range tests {
		// 使用 t.Run() 函数为每个测试用例运行一个子测试。第一个参数是测试名称（用于在任何日志输出中标识子测试），第二个参数是匿名函数，其中包含每个案例的实际测试内容
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
