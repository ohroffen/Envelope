package api

import (
	"MyEnvelope/algo"
	"MyEnvelope/entity"
	"MyEnvelope/mq"
	"MyEnvelope/my_redis"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
)

var node *snowflake.Node

func Init_snowflake_node() {
	id, err := my_redis.Rdb.Incr("node_id").Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("node id: %v", id)
	node, _ = snowflake.NewNode(id)
}

//定义红包信息结构体
type EnvelopeInfo struct {
	SnatchTime int64
	Money      int64
	Opened     bool
}

//结构体序列化
func (e EnvelopeInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func SnatchHandlerRedis(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("snatched by %s", uid)
	if uid == "" {
		c.JSON(200, gin.H{
			"code": 4,
			"msg":  "uid is empty",
		})
		return
	}

	// logic start
	//1、判断用户是否在一定概率能抢到
	algo.Init()
	if rand.Float64() > algo.SnatchRatio {
		//pipe.Exec()
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Sorry, you didn't catch the envelope. Good luck next time!",
		})
		return
	}

	/*
	//程序退出时解锁
	defer my_redis.UnLock(uid)
	//对当前用户增加一个锁，如果没有拿到当前锁，则间隔5ms的时间再次发送请求
	iCount := 5
	for i := 0; i < iCount; i++ {
		lockSuccess, err := my_redis.Rdb.SetNX(uid+"valid", 1, time.Second*2).Result()
		if err != nil || lockSuccess != true {
			log.Printf("%v get lock fail %v", uid, err)
			if i == iCount-1 {
				c.JSON(200, gin.H{
					"code": 5,
					"msg":  "too many requests",
				})
				return
			}
		} else {
			log.Printf("%v get lock success", uid)
			break
		}
		time.Sleep(3 * time.Millisecond)
	}
	*/

	//2、在redis中校验用户是否还有剩余次数
	//curCountStr, _ := my_redis.Rdb.HGet("user_count", uid).Result()
	//curCountStr, _ := my_redis.Rdb.HGet("user_count", uid).Result()
	curCount, _ := my_redis.Rdb.HIncrBy("user_count", uid, int64(1)).Result()
	//curCountStr, _ := curCount, _ := strconv.Atoi(curCountStr)
	//log.Printf("%v",curCount)不存在的话为0
	if curCount > int64(algo.MaxSnatchCount) {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "Sorry, you have used up your snatch count",
		})
		return
	}

	//3、判断队列中是否还有剩余的红包，存在或者直接出队列
	money, err := my_redis.Rdb.LPop("envelope_list").Result()
	if err != nil {
		//将用户增加的抢红包数要减1
		my_redis.Rdb.HSet("user_count", uid, curCount-1).Result()
		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "Sorry, There is no red envelope left!",
		})
		return
	}

	//增加用户抢红包次数
	//curCount++
	//my_redis.Rdb.HSet("user_count", uid, curCount)
	//my_redis.Rdb.HIncrBy("user_count", uid, int64(1))
	envelopeId := int64(node.Generate())
	snatchTime := time.Now().Unix() //获得当前时间戳，单位为s
	amount, _ := strconv.ParseInt(money, 10, 64)
	//将红包id添加到用户的未拆set中
	envelopeInfo := EnvelopeInfo{
		SnatchTime: snatchTime,
		Money:      amount,
		Opened:     false,
	}
	result, errInfo := my_redis.Rdb.HSet(uid+"list", fmt.Sprint(envelopeId), envelopeInfo).Result()
	if errInfo != nil {
		log.Printf("list set error,%v", errInfo)
	} else {
		log.Printf("%v insert envelope %v is %v", uid, envelopeId, result)
	}
	num_uid, _ := strconv.ParseInt(uid, 10, 64)
	//4、红包入消息队列
	mq.Send_message(&entity.Envelope{
		EnvelopeID: envelopeId,
		UserID:     num_uid,
		Opened:     false,
		Value:      amount,
		SnatchTime: snatchTime,
	})

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
