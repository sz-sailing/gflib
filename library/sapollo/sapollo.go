package sapollo

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/philchia/agollo/v4"
	"strconv"
)

type Config struct {
	Appid      string
	Cluster    string
	Namespaces []string
	Addrs      map[string]string
}

//请求apollo获取配置文件
func Start(ac Config) {
	//如果没有配置环境变量，默认是本地环境，不走apollo
	env := gstr.ToLower(genv.Get("ENV"))
	if env == "" {
		_ = g.Cfg().SetPath("etc")
		g.Cfg().SetFileName("config-local.yaml")
		g.Log().Info("环境变量未配置,读取配置文件: config-local.yaml")
		return
	}
	err := agollo.Start(&agollo.Conf{
		AppID:          ac.Appid,
		Cluster:        ac.Cluster,
		NameSpaceNames: ac.Namespaces,
		MetaAddr:       ac.Addrs[env],
	}, agollo.SkipLocalCache())
	if err != nil {
		panic(err)
	}
	agollo.OnUpdate(func(event *agollo.ChangeEvent) {
		for key, value := range event.Changes {
			setConfig(key, value.NewValue)
		}
	})
	//先将yaml的配置写文件，并且设置默认的配置文件目录和配置文件
	setConfig("yaml", agollo.GetString("yaml"))
	//获取其他的配置写入内存
	allKeys := agollo.GetAllKeys()
	for _, key := range allKeys {
		if key == "yaml" {
			continue
		}
		setConfig(key, agollo.GetString(key))
	}
}

// 写入配置项
func setConfig(key string, value string) {
	//yaml 直接写入文件
	if key == "yaml" {
		writeFile(key, value)
		return
	}
	if value == "true" {
		err := g.Cfg().Set(key, true)
		if err != nil {
			g.Log().Error(err)
		}
		return
	}
	if value == "false" {
		err := g.Cfg().Set(key, false)
		if err != nil {
			g.Log().Error(err)
		}
		return
	}
	// 整型的情况
	if i, err := strconv.Atoi(value); err == nil {
		err := g.Cfg().Set(key, i)
		if err != nil {
			g.Log().Error(err)
		}
		return
	}
	// 浮点数的情况
	if f, err := strconv.ParseFloat(value, 10); err == nil {
		err := g.Cfg().Set(key, f)
		if err != nil {
			g.Log().Error(err)
		}
		return
	}
	err := g.Cfg().Set(key, value)
	if err != nil {
		g.Log().Error(err)
	}
	return
}

// 将apollo的配置写入文件
func writeFile(key string, value string) {
	if key != "yaml" {
		return
	}
	dir := "etc"
	err := gfile.Mkdir(dir)
	if err != nil {
		g.Log().Error(err)
		panic(err)
	}
	filename := "config-" + gstr.ToLower(genv.Get("ENV")) + ".yaml"
	path := dir + "/" + filename
	err = gfile.PutContents(path, value)
	if err != nil {
		g.Log().Error(err)
		panic(err)
	}
	//设置默认的配置文件目录和配置文件
	err = g.Cfg().SetPath(dir)
	if err != nil {
		g.Log().Error(err)
		panic(err)
	}
	g.Cfg().SetFileName(filename)
	g.Log().Info("读取配置文件: %v", path)
	return
}
