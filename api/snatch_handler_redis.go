package api

import (
	"MyEnvelope/algo"
	"MyEnvelope/my_redis"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func SnatchHandlerRedis(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("snatched by %s", uid)
	if uid == ""{
		c.JSON(200, gin.H{
			"code": 4,
			"msg":  "uid is empty",
		})
		return
	}

	//对当前用户增加一个锁，如果没有拿到当前锁，则再次发送请求
	iCount := 5
	for i:=0 ; i< iCount; i++{
		lockSuccess, err := my_redis.Rdb.SetNX(uid+"valid",1,time.Second*2).Result()
		if err != nil || lockSuccess != true{
			log.Printf("%v get lock fail %v",uid,err)
			if i== iCount-1{
				my_redis.UnLock(uid)
				c.JSON(200, gin.H{
					"code": 5,
					"msg":  "too many requests",
				})
				return
			}
		} else{
			log.Printf("%v get lock success",uid)
			break
		}
	}

	// logic start
	//1、在redis中校验用户是否还有剩余次数
	curCountStr, err := my_redis.Rdb.Get(uid).Result()
	if err != nil{
		my_redis.Rdb.Set(uid,0,0)
		my_redis.Rdb.HSet("user_money",uid,0)
	}
	curCount, _ := strconv.Atoi(curCountStr)
	//log.Printf("%v",curCount)不存在的话为0
	if curCount >= algo.MaxSnatchCount {
		//pipe.Exec()
		my_redis.UnLock(uid)
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "Sorry, you have used up your snatch count",
		})
		return
	}

	//2、判断用户是否在一定概率能抢到
	algo.Init()
	if rand.Float64() > algo.SnatchRatio {
		//pipe.Exec()
		my_redis.UnLock(uid)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Sorry, you didn't catch the envelope. Good luck next time!",
		})
		return
	}

	//3、判断队列中是否还有剩余的红包，存在或者直接出队列
	money, err := my_redis.Rdb.LPop("envelope_list").Result()
	if err != nil {
		//pipe.Exec()
		my_redis.UnLock(uid)
		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "Sorry, There is no red envelope left!",
		})
		return
	}

	//增加用户抢红包次数
	curCount = int(my_redis.Rdb.Incr(uid).Val())
	//释放锁，因为后面的逻辑与红包剩余数量无关，所以可以提前释放，没有必要等到函数返回
	my_redis.UnLock(uid)

	envelopeId := uuid.New()
	snatchTime := time.Now().Unix()//获得当前时间戳，单位为s
	//将红包id添加到用户的未拆set中
	my_redis.Rdb.ZAdd(uid+"closed", redis.Z{
		Score: float64(snatchTime),
		Member: envelopeId,
	})
	//设置红包与对应金额之间的联系hash
	my_redis.Rdb.HSet("envelope_money",envelopeId,money)

	//todo 4、红包入消息队列
	//4、用户已经抢到红包，需要插入到数据库中
	//生成红包实体
	//value, err := strconv.ParseInt(money, 10, 64)
	//envelope := entity.Envelope{
	//	EnvelopeID: envelopeId,
	//	UserID:     uid,
	//	Opened:     false,
	//	Value:      value,
	//	SnatchTime: entity.UnixTime(snatchTime),
	//}
	////todo 插入失败的情况还没有做处理，不是很确定应该在这里进行插入
	//dao.InsertEnvelope(envelope)

	//pipe.Exec()
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"envelope_id": envelopeId,
			"max_count":   algo.MaxSnatchCount,
			"cur_count":   curCount,
		},
	})
}
