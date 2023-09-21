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




