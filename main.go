package main

import (
	"MyEnvelope/algo"
	"MyEnvelope/api"
	"MyEnvelope/mq"
	"MyEnvelope/my_redis"
	"MyEnvelope/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	conf = pflag.StringP("config", "c", "", "config filepath")
)

func main() {
	pflag.Parse()

	if err := utils.Run(*conf); err != nil {
		panic(err)
	}
	fmt.Println(viper.AllSettings())
	algo.InitConfig()

	//初始化Redis
	if err := my_redis.InitRedis(); err != nil {
		log.Printf("init my_redis client failed, err:%v\n", err)
	} else {
		log.Printf("connect my_redis success...")
	}
	defer my_redis.Rdb.Close()
	//清空之前的缓存，云redis没有权限执行FlushDB()，必须手动删除
	//my_redis.Rdb.FlushDB()
	//缓存预热
	my_redis.PreAllocated()

	// 初始化 snowflake 节点, 从 Redis 获取节点 Id
	api.Init_snowflake_node()
	// 初始化 kafka writer
	mq.Mq_init()

	r := gin.Default()

	// router
	r.POST("/snatch", api.SnatchHandlerRedis)
	r.POST("/open", api.OpenHandler)
	r.POST("/get_wallet_list", api.WalletListHandler)
	err := r.Run(":9090")
	if err != nil {
		return
	} //监听并在127.0.0.1:9090. 上启动服务
}
