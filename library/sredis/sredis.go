package sredis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

const (
	// DEFAULT_NAME Default group name for instance usage.
	DEFAULT_NAME     = "default"
	gREDIS_NODE_NAME = "redis"
)

var (
	// Instances map containing configuration instances.
	instances = gmap.NewStrAnyMap(true)
)

//Client 返回一个 go-redis 的单例
//UniversalClient is an abstract client which - based on the provided options - represents either a ClusterClient, a FailoverClient, or a single-node Client. This can be useful for testing cluster-specific applications locally or having different clients in different environments.
//
//NewUniversalClient returns a new multi client. The type of the returned client depends on the following conditions:
//
//If the MasterName option is specified, a sentinel-backed FailoverClient is returned.
//if the number of Addrs is two or more, a ClusterClient is returned.
//Otherwise, a single-node Client is returned.
//For example:
//
//
//// rdb is *redis.Client.
//rdb := NewUniversalClient(&redis.UniversalOptions{
//    Addrs: []string{":6379"},
//})
//
//// rdb is *redis.ClusterClient.
//rdb := NewUniversalClient(&redis.UniversalOptions{
//    Addrs: []string{":6379", ":6380"},
//})
//
//// rdb is *redis.FailoverClient.
//rdb := NewUniversalClient(&redis.UniversalOptions{
//    Addrs: []string{":6379"},
//    MasterName: "mymaster",
//})
func Client(name ...string) redis.UniversalClient {
	config := gins.Config()
	key := DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	var opts *redis.UniversalOptions
	return instances.GetOrSetFuncLock(key, func() interface{} {
		var m map[string]interface{}
		if _, v := gutil.MapPossibleItemByKey(gins.Config().GetMap("."), gREDIS_NODE_NAME); v != nil {
			m = gconv.Map(v)
		}
		if len(m) > 0 {
			if v, ok := m[key]; ok {
				err := gconv.Struct(v, &opts)
				if err != nil {
					panic(err)
				}
				return redis.NewUniversalClient(opts)
			} else {
				panic(fmt.Sprintf(`configuration for redis not found for group "%s"`, key))
			}
		} else {
			panic(fmt.Sprintf(`incomplete configuration for redis: "redis" node not found in config file "%s"`, config.GetFileName()))
		}
		return nil
	}).(redis.UniversalClient)
}
