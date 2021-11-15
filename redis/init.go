package redis

import (
	"os"

	"github.com/go-redis/redis"
)

//全局的rdb变量
var Rdb *redis.Client

//初始化redis
func InitRedis() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,   //数据库的索引
		PoolSize: 100, //连接池大小
	})
	_, err := Rdb.Ping().Result()
	return err
}
