## 依赖管理

- 安装单个依赖 
```shell
go get github.com/go-sql-driver/mysql@v1
```
- 更新
```shell
go get -u github.com/foo/bar@v2.0.0
```
- 检验
```shell
go mod verify
```
- 安装项目所有依赖
```shell
 go mod download
```
- 删除
```shell
go get github.com/foo/bar@none
```
> 或者，如果你已经删除了代码中对软件包的所有引用，可以运行 go mod tidy，它会自动从 go.mod 和 go.sum 文件中删除任何未使用的软件包。
```shell
go mod tidy -v
```

## SameSite cookie 和 TLS 1.3
在本章前面，我说过我们不能完全依赖 SameSite cookie 属性来防止 CSRF 攻击，因为并非所有浏览器都完全支持它。
但这一规则有一个例外，那就是不存在支持 TLS 1.3 但不支持 SameSite cookie 的浏览器。
换句话说，如果您在服务器的 TLS 配置中将 TLS 1.3 设置为最低支持版本，那么所有能使用您的应用程序的浏览器都将支持 SameSite Cookie。
```
tlsConfig := &tls.Config{
 MinVersion: tls.VersionTLS13,
}
```

只要只允许向应用程序发出 HTTPS 请求，并将 TLS 1.3 作为最低 TLS 版本，就不需要针对 CSRF 攻击采取任何额外的缓解措施（如使用 justinas/nosurf 软件包）。
- 只需确保您始终在会话 cookie 上设置 SameSite=Lax 或 SameSite=Strict；
- 对任何改变状态的请求使用 "不安全 "HTTP 方法（即 POST、PUT 或 DELETE）。

## 使用嵌入式文件
#### 用法
```go
package ui
import (
 "embed"
)
//go:embed "html" "static"
var Files embed.FS
```
这看起来像一个注释，但实际上是一个特殊的注释指令。编译应用程序时，这条注释指令会指示 Go 将 ui/html 和 ui/static 文件夹中的文件存储到全局变量 Files 引用的 embed.FS 嵌入式文件系统中。
#### 注意事项：
- 注释指令必须紧贴在要存储嵌入文件的变量上方。
- 该指令的一般格式为 go:embed <paths>，可以在一个指令中指定多个路径（如上面的代码）。路径应相对于包含该指令的源代码文件。因此，在我们的例子中，go:embed "static" "html "嵌入了项目中的 ui/static 和 ui/html 目录。
- 您只能在包级别的全局变量上使用 go:embed 指令，而不能在函数或方法中使用。如果试图在函数或方法中使用该指令，编译时会出现 "go:embed cannot apply to var inside func "的错误。
- 路径不能包含.或.元素，也不能以.开头或结尾。 这基本上限制了你只能嵌入与带有 go:embed 指令的源代码位于同一目录（或子目录）中的文件。
- 如果路径指向一个目录，那么该目录中的所有文件都会被递归嵌入，名称以 . 或 .. 开头的文件除外。如果要包含这些文件，应使用 all:前缀，如 go:embed "all:static"。
- 路径分隔符应始终为正斜线，即使在 Windows 机器上也是如此
- 嵌入式文件系统的根目录总是包含 go:embed 指令的目录。因此，在上面的示例中，我们的 Files 变量包含一个 embed.FS 嵌入式文件系统，而该文件系统的根目录就是我们的 ui 目录。

## 打包程序
```shell
go build -o ./releases ./cmd/web 
cp -r ./tls releases 
cp -r ./log releases
cd releases
./web
```

## 测试

```shell
// 运行 ./cmd/web 下的测试
go test ./cmd/web  
// 显示详细信息
go test -v ./cmd/web
```

### 自定义测试运行
- 使用 ./... 通配符模式运行当前项目中的所有测试
```shell
go test ./...
```
-  只运行 TestPing 测试，如下所示
```shell
 go test -v -run="^TestPing$" ./cmd/web/
```

> 使用 -run 标志将测试限制在某些特定的子测试中，格式为 {test regexp}/{sub-test regexp} 。例如，要运行 TestHumanDate 测试的 UTC 子测试，我们可以这样做
    
```shell
go test -v -run="^TestHumanDate$/^UTC$" ./cmd/web
```
- 跳过指定测试运行 -skip 标志
```shell
go test -v -skip="^TestHumanDate$" ./cmd/web/
```

### 测试的缓存
如果运行两次完全相同的测试，而不对测试的软件包做任何更改，那么就会显示测试结果的缓存版本（由软件包名称旁边的（cached）注释表示）。
```shell
go test ./cmd/web
ok snippetbox.alexedwards.net/cmd/web (cached)
```
如果想避免 test 的缓存影响 可以使用 -count 标志 保证 每个测试至少被运行 count 次

```shell
go test -count=1 ./cmd/web
```

或者 使用 `go clean` 清除 test 缓存

```shell
go test -testcache
```

### 故障快速反馈
当我们使用 t.Errorf() 函数将测试标记为失败时，并不会导致测试立即退出。所有其他测试（和子测试）都将在失败后继续运行。
如果想第一次出现错误就直接中止运行，可以使用 -failfast 标签
```shell
go test -failfast ./cmd/web
```
> 需要注意的是，"-failfast "标志只停止发生故障的软件包中的测试。如果你在多个软件包中运行测试（例如使用 go test ./...），那么其他软件包中的测试将继续运行。

### 并行测试
```go
func TestPing(t *testing.T) {
 t.Parallel()
 ...
}
```
- 使用 t.Parallel() 标记的测试将与且仅与其他并行测试并行运行。
- 默认情况下，同时运行的最大测试次数为 GOMAXPROCS 的当前值。您可以通过 -parallel 标志设置一个特定值来覆盖该值。例如
```shell
 go test -parallel=4 ./...
```
- 并非所有测试都适合并行运行。例如，如果集成测试要求数据库表处于特定的已知状态，那么就不希望与其他操作相同数据库表的测试并行运行。

### 启用竞争检测
```shell
go test -race ./cmd/web/
```

### 跳过长期运行的测试
跳过长时间运行的测试的一种常见惯用方法是使用 testing.Short() 函数检查 go 测试命令中是否存在 -short 标记，如果存在该标记，则调用 t.Skip() 方法跳过测试。

```go
func TestUserModelExists(t *testing.T) {
 // Skip the test if the "-short" flag is provided when running the test.
 if testing.Short() {
 t.Skip("models: skipping integration test")
 }
 ...
}
```
go test 运行命令添加 -short 标志时 则 testing.Short() 为 true
`go test -v -short ./...` 

### 分析测试覆盖率

- 命令行输出覆盖率情况(以包为单位)

`go test -cover ./...`

```shell
➜ go test -cover ./...                                        
?       github.com/hlf2016/snippetbox/internal/assert   [no test files]
?       github.com/hlf2016/snippetbox/internal/models/mocks     [no test files]
?       github.com/hlf2016/snippetbox/ui        [no test files]
?       github.com/hlf2016/snippetbox/internal/validator        [no test files]
ok      github.com/hlf2016/snippetbox/cmd/web   1.234s  coverage: 43.8% of statements
ok      github.com/hlf2016/snippetbox/internal/models   1.261s  coverage: 6.8% of statements

```

- 命令行输出覆盖率情况(以包为单位)，并将覆盖率元数据输出到 当前目录的 profile.out 文件

`go test -coverprofile=./profile.out ./...`
```shell
➜ go test -coverprofile=./profile.out ./...
?       github.com/hlf2016/snippetbox/internal/assert   [no test files]
?       github.com/hlf2016/snippetbox/internal/models/mocks     [no test files]
?       github.com/hlf2016/snippetbox/internal/validator        [no test files]
?       github.com/hlf2016/snippetbox/ui        [no test files]
ok      github.com/hlf2016/snippetbox/cmd/web   0.877s  coverage: 43.8% of statements
ok      github.com/hlf2016/snippetbox/internal/models   0.980s  coverage: 6.8% of statements
```

- 将 profile.out 中的元数据 转换成 覆盖率情况展示出来(以 func 为单位)

`go tool cover -func=./profile.out`

```shell
➜ go tool cover -func=./profile.out
github.com/hlf2016/snippetbox/cmd/web/handler.go:13:            home                    0.0%
github.com/hlf2016/snippetbox/cmd/web/handler.go:34:            snippetView             92.9%
github.com/hlf2016/snippetbox/cmd/web/handler.go:68:            snippetCreate           0.0%
github.com/hlf2016/snippetbox/cmd/web/handler.go:87:            snippetCreatePost       0.0%
github.com/hlf2016/snippetbox/cmd/web/handler.go:144:           userSignup              100.0%
github.com/hlf2016/snippetbox/cmd/web/handler.go:149:           userSignupPost          88.5%
github.com/hlf2016/snippetbox/cmd/web/handler.go:193:           userLogin               0.0%
github.com/hlf2016/snippetbox/cmd/web/handler.go:198:           userLoginPost           0.0%
github.com/hlf2016/snippetbox/cmd/web/handler.go:241:           userLogoutPost          0.0%
github.com/hlf2016/snippetbox/cmd/web/handler.go:252:           ping                    100.0%
github.com/hlf2016/snippetbox/cmd/web/helper.go:15:             serverError             0.0%
github.com/hlf2016/snippetbox/cmd/web/helper.go:24:             clientError             100.0%
github.com/hlf2016/snippetbox/cmd/web/helper.go:29:             notFound                100.0%
github.com/hlf2016/snippetbox/cmd/web/helper.go:33:             render                  57.1%
github.com/hlf2016/snippetbox/cmd/web/helper.go:62:             newTemplateData         100.0%
github.com/hlf2016/snippetbox/cmd/web/helper.go:73:             decodePostForm          50.0%
github.com/hlf2016/snippetbox/cmd/web/helper.go:92:             isAuthenticated         75.0%
github.com/hlf2016/snippetbox/cmd/web/main.go:40:               main                    0.0%
github.com/hlf2016/snippetbox/cmd/web/main.go:135:              openDB                  0.0%
github.com/hlf2016/snippetbox/cmd/web/middleware.go:11:         secureHeaders           100.0%
github.com/hlf2016/snippetbox/cmd/web/middleware.go:23:         neuter                  0.0%
github.com/hlf2016/snippetbox/cmd/web/middleware.go:33:         logRequest              100.0%
github.com/hlf2016/snippetbox/cmd/web/middleware.go:40:         recoverPanic            66.7%
github.com/hlf2016/snippetbox/cmd/web/middleware.go:59:         requireAuthentication   16.7%
github.com/hlf2016/snippetbox/cmd/web/middleware.go:72:         authenticate            38.5%
github.com/hlf2016/snippetbox/cmd/web/middleware.go:95:         noSurf                  100.0%
github.com/hlf2016/snippetbox/cmd/web/routes.go:10:             routes                  100.0%
github.com/hlf2016/snippetbox/cmd/web/templates.go:23:          humanDate               100.0%
github.com/hlf2016/snippetbox/cmd/web/templates.go:37:          newTemplateCache        83.3%
github.com/hlf2016/snippetbox/internal/models/snippets.go:27:   Insert                  0.0%
github.com/hlf2016/snippetbox/internal/models/snippets.go:42:   Get                     0.0%
github.com/hlf2016/snippetbox/internal/models/snippets.go:60:   Latest                  0.0%
github.com/hlf2016/snippetbox/internal/models/users.go:30:      Insert                  0.0%
github.com/hlf2016/snippetbox/internal/models/users.go:51:      Authenticate            0.0%
github.com/hlf2016/snippetbox/internal/models/users.go:76:      Exists                  100.0%
total:                                                          (statements)            37.0%

```

- 自动打开浏览器，以 html 形式查看 测试覆盖情况

`go tool cover -html=./profile.out`

> 使用 -covermode=count 会使覆盖率配置文件记录测试过程中每条语句被执行的确切次数，而不是仅仅突出显示绿色和红色语句。

```shell
go test -covermode=count -coverprofile=./profile.out ./...
go tool cover -html=./profile.out
```


