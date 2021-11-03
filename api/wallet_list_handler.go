package api

import (
	"MyEnvelope/dao"
	"log"

	"github.com/gin-gonic/gin"
)

func WalletListHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("query %s's wallet", uid)

	// logic start
	// 通过用户ID查询到多个红包模型，然后直接返回
	// TODO 按时间排序
	envelopes, totalAmount := dao.GetEnvelopesByUserId(uid)
	amountOfEnvelopes := len(envelopes)

	// 各种特殊情况应该返回什么状态码
	if amountOfEnvelopes == 0 {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "no envelope",
			"data": gin.H{
				"amount":        amountOfEnvelopes,
				"envelope_list": envelopes,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"amount":        totalAmount,
				"envelope_list": envelopes,
			},
		})
	}
	//fmt.Println(envelopes[0])

	// 红包字段

	//envelopes := []gin.H{
	//	{
	//		"envelope_id": 123,
	//		"value":       50,
	//		"opened":      true,
	//		"snatch_time": 1634551711,
	//	},
	//	{
	//		"envelope_id": 234,
	//		"opened":      false,
	//		"snatch_time": 1634551711,
	//	},
	//}
}
