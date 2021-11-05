package my_redis

import (
	"log"
)

//对当前用户增加一个锁，如果没有拿到当前锁，则再次发送请求
func UnLock(uid string) {
	_, err := Rdb.Del(uid + "valid").Result()
	if err != nil {
		log.Printf("%v, unlock fail,%v", uid, err)
	} else {
		log.Printf("%v unlock success", uid)
	}
}
