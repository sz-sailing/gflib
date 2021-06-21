# Apollo 配置中心连接文档

支持 goframe 框架自动从 Apollo 获取配置文件

## 配置项说明

在 apollo 后台的配置项：

- KEY 名的层级用 “.” 分隔
- KEY 名为 `yaml` 的配置项会被写入本地文件，文件路径为：`etc/config-[env].yaml`，由于会原样将 value 写入文件，因此，配置 `yaml` 这个 key 的时候，value 请严格保持 yaml 格式
- `yaml` 这个 KEY 的值在 apollo 中被修改后，会被重新写入文件，覆盖之前的文件，gofram 框架文档上面说框架会自动检测配置文件的变动，然后重新加载配置，我没有做测试，建议修改了 `yaml` 的值后重启程序，以获取最新的配置
- 其他的 key-value 从 apollo 读取后会写入内存，值被修改后会实时获取，并重新写入内存
- value 支持 bool （true/false）、整型、浮点型
- 更复杂的数据类型请用 yaml 格式，直接写入 `yaml` 这个 KEY 的 value 中

## 使用示例

服务器需要配置环境变量`ENV`，`ENV` 的值有 `dev`、`fat`、`uat`、`pro`，如果没有配置环境变量，则默认是本地环境，会直接读取本地 `etc/config-local.yaml` 文件，不会请求 apollo

在项目的入口文件填写好 apollo 的配置，然后启动就行，程序会根据服务器的环境请求对应的 apollo 地址，获取配置

```go
//_ = genv.Set("ENV", "FAT")
var ac = sapollo.Config{
    Appid:      "example-appid",
    Cluster:    "example-cluster",
    Namespaces: []string{"application"},
    Addrs: map[string]string{
        "dev": "http://dev-internal.apollo-config-center.com",
        "fat": "http://test-internal.apollo-config-center.com",
        "uat": "http://gray-internal.apollo-config-center.com",
        "pro": "http://prod-internal.apollo-config-center.com",
    },
}
sapollo.Start(ac)
//g.Cfg().Dump()
```

