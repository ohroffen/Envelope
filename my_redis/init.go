package my_redis

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

//全局的rdb变量
var Rdb *redis.Client

//初始化redis
func InitRedis() error {
	Rdb = redis.NewClient(&redis.Options{
		//Addr: "redis-cn02pak6xbbkuiu9x.redis.volces.com:6379",
		Addr: viper.GetString("redis.host"),
		//Addr: "localhost:6379",
		//Password: "Group12345678",
		Password: viper.GetString("redis.password"),
		//Password: "",
		DB:       0,   //数据库的索引
		PoolSize: 100, //连接池大小
	})
	_, err := Rdb.Ping().Result()
	return err
}
