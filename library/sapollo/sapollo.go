package sapollo

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/philchia/agollo/v4"
	"strconv"
	"strings"
)

type Config struct {
	Appid      string
	Cluster    string
	Namespaces []string
	Address    string
}

// Start 请求apollo获取配置文件
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
		MetaAddr:       ac.Address,
	}, agollo.SkipLocalCache())
	if err != nil {
		panic(err)
	}
	agollo.OnUpdate(func(event *agollo.ChangeEvent) {
		setConfig(ac.Namespaces)
	})
	setConfig(ac.Namespaces)

}

//设置配置
func setConfig(namespaces []string) {
	var yamlContents string
	//获取各个yaml命名空间的配置内容，写到文件里面
	for _, namespace := range namespaces {
		if strings.HasSuffix(namespace, ".yaml") {
			content := agollo.GetContent(agollo.WithNamespace(namespace))
			yamlContents += "\n\n" + content
		}
	}
	writeYamlFile(yamlContents)
	//获取properties命名空间的配置，写内存
	for _, namespace := range namespaces {
		if strings.HasSuffix(namespace, ".properties") {
			//获取其他的配置写入内存
			allKeys := agollo.GetAllKeys(agollo.WithNamespace(namespace))
			for _, key := range allKeys {
				if key == "yaml" {
					continue
				}
				setKeyValue(key, agollo.GetString(key))
			}
		}
	}
}

// 写入配置项
func setKeyValue(key string, value string) {
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
func writeYamlFile(contents string) {
	dir := "etc"
	err := gfile.Mkdir(dir)
	if err != nil {
		g.Log().Error(err)
		panic(err)
	}
	filename := "config-" + gstr.ToLower(genv.Get("ENV")) + ".yaml"
	path := dir + "/" + filename
	err = gfile.PutContents(path, contents)
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
