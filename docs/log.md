# log 方法扩展

slog 的出现是为了满足 log 的格式与 php 环境保持一致，内部的日志标准

## 配置示例

公司内部的 go 项目都是容器部署，日志只需要输出到终端就行，不用写文件，因此，请保持配置中的 Path 参数为空。同时 Stdout 参数为 true

```yaml
# Logger.
logger:
  Path: ""
  Level: "all"
  Stdout: true
```

其他配置完全遵循 goframe 官方的 log 配置。

## 使用示例

- 在中间件将请求的部分参数写入上下文

```go
func (m *middleware) Base(r *ghttp.Request) {
	r.SetCtxVar("remote_addr", r.GetRemoteIp())
	r.SetCtxVar("domain", r.GetHost())
	r.SetCtxVar("referer", r.GetReferer())
	r.SetCtxVar("shop_id", r.GetHeader("shop_id"))
	r.SetCtxVar("request_id", r.GetHeader("request_id"))
	
	//或者直接调用下面方法
	//slog.SetCtxVar(*ghttp.Request)

	r.Middleware.Next()
}
```

- 在控制器获取上下文，继续往下传，或者直接写日志
```go
//直接写日志的情况
var ctx = r.GetCtx()
slog.Init(ctx).Error("我是日志啊")

//继续往下传的情况	
err := PutThemePack(ctx,shopId,themeId)
func PutThemePack(ctx context.Context,shopId int64, themeId int64)() error  {
    ...
    slog.Init(ctx).Error("这里出错啦")
    return nil
}
slog.Init(ctx),Error("这里也出错啦")
```

## 日志示例

```json
{
    "prefix":"[192.168.3.246][-][-]",
    "time":"2021-05-06 17:28:58",
    "request_id":"5b47bfeb3f102f59f960214857b0a766",
    "level":"ERRO",
    "domain":"api.shopcity.vip",
    "shopid":"25",
    "referer":"https://jeff.shopcity.vip/",
    "remote_addr":"192.168.3.246",
    "message":"这里出错了啊",
    "trace":"1.  slog/app/library/slog.(*lo\ng).Error\n    D:/code/test/slog/app/library/slog/slog.go:40\n2.  slog/app/api.(*helloApi).Index\n    D:/code/test/slog/app/api/hello.go:15\n3.  github.com/gogf/gf/net/ghttp.(*middleware).Next.func1.6\n    C:/Users/dabuge-micro/go/pk\ng/mod/github.com/gogf/gf@v1.16.2/net/ghttp/ghttp_request_middleware.go:95\n4.  github.com/gogf/gf/net/ghttp.niceCallFunc\n    C:/Users/dabuge-micro/go/pkg/mod/github.com/gogf/gf@v1.16.2/net/ghttp/ghttp_func.go:46\n5.  github.com/gog\nf/gf/net/ghttp.(*middleware).Next.func1\n    C:/Users/dabuge-micro/go/pkg/mod/github.com/gogf/gf@v1.16.2/net/ghttp/ghttp_request_middleware.go:94\n6.  github.com/gogf/gf/util/gutil.TryCatch\n    C:/Users/dabuge-micro/go/pkg/mod/gith\nub.com/gogf/gf@v1.16.2/util/gutil/gutil.go:46\n7.  github.com/gogf/gf/net/ghttp.(*middleware).Next\n    C:/Users/dabuge-micro/go/pkg/mod/github.com/gogf/gf@v1.16.2/net/ghttp/ghttp_request_middleware.go:47\n8.  slog/app/middleware.(*\nmiddleware).Auth\n    D:/code/test/slog/app/middleware/auth.go:15\n9.  github.com/gogf/gf/net/ghttp.(*middleware).Next.func1.1\n    C:/Users/dabuge-micro/go/pkg/mod/github.com/gogf/gf@v1.16.2/net/ghttp/ghttp_request_middleware.go:53\n\n10. github.com/gogf/gf/net/ghttp.niceCallFunc\n    C:/Users/dabuge-micro/go/pkg/mod/github.com/gogf/gf@v1.16.2/net/ghttp/ghttp_func.go:46\n11. github.com/gogf/gf/net/ghttp.(*middleware).Next.func1\n    C:/Users/dabuge-micro/go/pkg\n/mod/github.com/gogf/gf@v1.16.2/net/ghttp/ghttp_request_middleware.go:52\n12. github.com/gogf/gf/util/gutil.TryCatch\n    C:/Users/dabuge-micro/go/pkg/mod/github.com/gogf/gf@v1.16.2/util/gutil/gutil.go:46\n13. github.com/gogf/gf/net\n/ghttp.(*middleware).Next\n    C:/Users/dabuge-micro/go/pkg/mod/github.com/gogf/gf@v1.16.2/net/ghttp/ghttp_request_middleware.go:47\n14. github.com/gogf/gf/net/ghttp.(*Server).ServeHTTP\n    C:/Users/dabuge-micro/go/pkg/mod/github.c\nom/gogf/gf@v1.16.2/net/ghttp/ghttp_server_handler.go:122\n"
}
```
