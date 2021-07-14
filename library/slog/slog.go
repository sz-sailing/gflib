package slog

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

type log struct {
	*glog.Logger
	Ctx context.Context
}

// LogStd 赛凌标准日志输出格式
type LogStd struct {
	Prefix     string      `json:"prefix,omitempty"`
	Time       string      `json:"time,omitempty"`
	RequestId  string      `json:"request_id,omitempty"`
	Level      string      `json:"level,omitempty"`
	Domain     string      `json:"domain,omitempty"`
	ShopId     int64       `json:"shop_id,omitempty"`
	Referer    string      `json:"referer,omitempty"`
	RemoteAddr string      `json:"remote_addr,omitempty"`
	Message    interface{} `json:"message,omitempty"`
	Trace      string      `json:"trace,omitempty"`
}

func Init(ctx context.Context, name ...string) *log {
	return &log{g.Log(name...).Async(true), ctx}
}

//SetCtxVar 根据http请求信息，在http请求中设置日志使用的上下文变量
func SetCtxVar(r *ghttp.Request) {
	r.SetCtxVar("remote_addr", r.GetRemoteIp())
	r.SetCtxVar("domain", r.GetHost())
	r.SetCtxVar("referer", r.GetReferer())
	r.SetCtxVar("shop_id", r.GetHeader("shop_id"))
	r.SetCtxVar("request_id", r.GetHeader("request_id"))
}

// GetCtxVar 从请求中获取日志需要的上下文内容
func GetCtxVar(r *ghttp.Request) context.Context {
	Ctx := context.Background()
	Ctx = context.WithValue(Ctx, "remote_addr", r.GetRemoteIp())
	Ctx = context.WithValue(Ctx, "domain", r.GetHost())
	Ctx = context.WithValue(Ctx, "referer", r.GetReferer())
	Ctx = context.WithValue(Ctx, "shop_id", r.GetHeader("shop_id"))
	Ctx = context.WithValue(Ctx, "request_id", r.GetHeader("request_id"))
	return Ctx
}

func (l *log) Error(e error) {
	var los = LogStd{
		Prefix:     "[" + gconv.String(l.Ctx.Value("remote_addr")) + "][-][-]",
		Time:       gtime.Datetime(),
		RequestId:  gconv.String(l.Ctx.Value("request_id")),
		Level:      l.GetLevelPrefix(glog.LEVEL_ERRO),
		Domain:     gconv.String(l.Ctx.Value("domain")),
		Trace:      l.GetStack(1),
		ShopId:     gconv.Int64(l.Ctx.Value("shop_id")),
		Referer:    gconv.String(l.Ctx.Value("referer")),
		RemoteAddr: gconv.String(l.Ctx.Value("remote_addr")),
	}
	if e == nil {
		los.Message = nil
	} else {
		los.Message = e.Error()
	}
	l.Logger.Error(los)
}

func (l *log) Redis() *glog.Logger {
	l.Logger.SetPrefix("[redis]")
	return l.Logger
}

func (l *log) Mongodb() *glog.Logger {
	l.Logger.SetPrefix("[mongodb]")
	return l.Logger
}

func (l *log) Mysql() *glog.Logger {
	l.Logger.SetPrefix("[mysql]")
	return l.Logger
}

func (l *log) Cache() *glog.Logger {
	l.Logger.SetPrefix("[cache]")
	return l.Logger
}

func (l *log) S3() *glog.Logger {
	l.Logger.SetPrefix("[s3]")
	return l.Logger
}
