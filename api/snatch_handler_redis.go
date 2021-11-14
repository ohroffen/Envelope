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
			"msg":  "invalid input",
		})
		return
	}

	// logic start
	//1、判断用户是否在一定概率能抢到
	algo.Init()
	if rand.Float64() > algo.SnatchRatio {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "miss",
		})
		return
	}

	//2、在redis中校验用户是否还有剩余次数
	curCount, _ := my_redis.Rdb.HIncrBy("user_count", uid, int64(1)).Result()
	if curCount > int64(algo.MaxSnatchCount) {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "snatch count used up",
		})
		return
	}

	//3、判断队列中是否还有剩余的红包，存在或者直接出队列
	money, err := my_redis.Rdb.LPop("envelope_list").Result()
	if err != nil {
		//将用户多增加的抢红包数要减1
		my_redis.Rdb.HSet("user_count", uid, curCount-1).Result()
		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "no more envelope",
		})
		return
	}

	//增加用户抢红包次数
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
