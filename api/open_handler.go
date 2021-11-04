package api

import (
	"MyEnvelope/my_redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"strconv"
)

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")
	log.Printf("envelope %s opened by %s", envelopeId, uid)
	//参数的合法
	if uid == "" || envelopeId == ""{
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "uid or envelopeid is empty",
		})
		return
	}

	time, err := my_redis.Rdb.ZScore(uid+"closed",envelopeId).Result()
	//1、判断用户是否有相应的红包
	if err != nil{
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Given user don't have such envelope",
		})
	}else{
		//2、未拆红包出集合
		my_redis.Rdb.ZRem(uid+"closed",envelopeId)
		my_redis.Rdb.ZAdd(uid+"open", redis.Z{
			Score: time,
			Member: envelopeId,
		})
		valueStr := my_redis.Rdb.HGet("envelope_money",envelopeId).Val()
		value, _ := strconv.ParseInt(valueStr, 10, 64)
		//红包金额入账，更新个人总金额
		my_redis.Rdb.HIncrBy("user_money",uid,value).Result()

		//todo 时间戳的格式与数据库不一致
		//valueStr, _ := strconv.ParseInt(value, 10, 64)
		//envelope := entity.Envelope{
		//	EnvelopeID: envelopeId,
		//	UserID:     uid,
		//	Opened:     true,
		//	Value:      valueStr,
		//	SnatchTime: time,
		//}
		//dao.updateEnvelope(envelope)

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{
				"value": value,
			},
		})
	}
	// logic start
	// 直接查询到该红包ID，然后返回
	//envelope := dao.GetEnvelopeByUserIdAndEnvelopeId(uid, envelopeId)
	//
	//value := envelope.Value
	//
	//// 修改红包状态为打开状态
	//envelope.Opened = true
	//dao.UpdateOpenState(&envelope)
	//// logic end
	//
	//if value > 0 {
	//	c.JSON(200, gin.H{
	//		"code": 0,
	//		"msg":  "Success",
	//		"data": gin.H{
	//			"value": value,
	//		},
	//	})
	//} else {
	//	c.JSON(200, gin.H{
	//		"code": 1,
	//		"msg":  "Given user don't have such envelope",
	//		"data": gin.H{
	//			"value": 0,
	//		},
	//	})
	//}

}
