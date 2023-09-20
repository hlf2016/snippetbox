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



