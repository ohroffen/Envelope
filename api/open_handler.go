package api

import (
	"MyEnvelope/entity"
	"MyEnvelope/mq"
	"MyEnvelope/my_redis"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")
	log.Printf("envelope %s opened by %s", envelopeId, uid)
	//参数的合法
	if uid == "" || envelopeId == "" {
		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "invalid input",
		})
		return
	}

	resultStr, err := my_redis.Rdb.HGet(uid+"list", envelopeId).Result()
	//1、判断用户是否有相应的红包
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "envelope doesn't exist",
		})
		return
	}
	//2、更改红包状态，直接返回金额，不更新用户总金额
	envelopeInfo := EnvelopeInfo{}
	json.Unmarshal([]byte(resultStr), &envelopeInfo)
	if envelopeInfo.Opened {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "already opened",
		})
		return
	} else {
		envelopeInfo.Opened = true
		my_redis.Rdb.HSet(uid+"list", envelopeId, envelopeInfo)
		num_envelopeId, _ := strconv.ParseInt(envelopeId, 10, 64)
		num_uid, _ := strconv.ParseInt(uid, 10, 64)
		mq.Send_message(&entity.Envelope{
			EnvelopeID: num_envelopeId,
			UserID:     num_uid,
			Opened:     true,
			Value:      envelopeInfo.Money,
			SnatchTime: envelopeInfo.SnatchTime,
		})
	}
	//Attention: 在拆红包接口没有将金额入账，而是选择在获取红包列表的时候更新

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "Success",
		"data": gin.H{
			"value": envelopeInfo.Money,
		},
	})
}
