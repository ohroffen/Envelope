package main

import (
	"MyEnvelope/api"
	"MyEnvelope/mq"
	"MyEnvelope/redis"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化Redis
	if err := redis.InitRedis(); err != nil {
		log.Printf("init my_redis client failed, err:%v\n", err)
	} else {
		log.Printf("connect my_redis success...")
	}
	defer redis.Rdb.Close()

	// 从 Redis 获取抢红包配置数据
	api.RetrieveSnatchConfig()
	// 开始获取红包金额
	api.FetchMoney()
	// 监听 kubernetes 发出的 SIGTERM 信号
	api.Terminate = make(chan os.Signal, 1)
	signal.Notify(api.Terminate, syscall.SIGTERM)
	// 初始化 kafka writer
	mq.Mq_init()

	r := gin.Default()

	// router
	r.POST("/snatch", api.SnatchHandler)
	r.POST("/open", api.OpenHandler)
	r.POST("/get_wallet_list", api.WalletListHandler)
	err := r.Run(":9090")

	if err != nil {
		return
	} //监听并在127.0.0.1:9090. 上启动服务
}
