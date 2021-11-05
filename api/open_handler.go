package api

import (
	"MyEnvelope/my_redis"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")
	log.Printf("envelope %s opened by %s", envelopeId, uid)
	//参数的合法
	if uid == "" || envelopeId == "" {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "uid or envelopeid is empty",
		})
		return
	}

	resultStr, err := my_redis.Rdb.HGet(uid+"list", envelopeId).Result()
	//1、判断用户是否有相应的红包
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Given user don't have such envelope",
		})
		return
	}
	//2、更改红包状态，直接返回金额，不更新用户总金额
	envelopeInfo := EnvelopeInfo{}
	json.Unmarshal([]byte(resultStr), &envelopeInfo)

	//todo 不确定是否需要增加另一个状态，重复拆红包
	if envelopeInfo.Opened == false {
		envelopeInfo.Opened = true
		my_redis.Rdb.HSet(uid+"list", envelopeId, envelopeInfo)
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
