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