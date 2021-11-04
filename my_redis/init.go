package my_redis

import (
	"MyEnvelope/algo"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"log"
)

//全局的rdb变量
var Rdb *redis.Client

//初始化redis
func InitRedis() error{
	Rdb = redis.NewClient(&redis.Options{
		//Addr: "redis-cn02pak6xbbkuiu9x.redis.volces.com:6379",
		Addr: viper.GetString("redis.host"),
		//Addr: "localhost:6379",
		//Password: "Group12345678",
		Password: viper.GetString("redis.password"),
		//Password: "",
		DB:0,//数据库的索引
		PoolSize: 100, //连接池大小
	})
	_, err := Rdb.Ping().Result()
	return err
}

//预热，预先将大红包分配成小红包
func PreAllocated() {
	for algo.TotalAmountOfEnvelope > int64(0) {
		money := algo.GetRandomMoney()
		n, err := Rdb.LPush("envelope_list",money).Result()
		if err != nil{
			log.Printf("insert failed, %v",err)
		}else{
			log.Printf("insert success %v,the value is %v",n,money)
		}
	}
}
