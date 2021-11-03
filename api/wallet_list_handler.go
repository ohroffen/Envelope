package api

import (
	"MyEnvelope/my_redis"
	"github.com/gin-gonic/gin"
	"log"
)

func WalletListHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("query %s's wallet", uid)
	if uid == ""{
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "uid is empty",
		})
		return
	}

	// logic start
	// 通过用户ID查询到多个红包模型，然后直接返回
	// 按时间排序
	closed_len := my_redis.Rdb.ZCard(uid+"closed").Val()
	open_len := my_redis.Rdb.ZCard(uid+"open").Val()
	if closed_len == 0 && open_len == 0{
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "envelope list is empty",
			"data": gin.H{
				"amount": 0,
			},
		})
		return
	}
	closed, err := my_redis.Rdb.ZRangeWithScores(uid+"closed",0,closed_len).Result()
	open, err := my_redis.Rdb.ZRangeWithScores(uid+"open",0,open_len).Result()
	money, err := my_redis.Rdb.HGet("user_money",uid).Result()

	if err != nil{
		log.Printf("error")
	}

	envelopes := make([]gin.H, 0)
	i := int64(0)
	j := int64(0)
	for i<closed_len && j<open_len {
		if closed[i].Score <= open[j].Score {
			temp := gin.H{
				"envelope_id": closed[i].Member,
				"opened":      false,
				"snatch_time": closed[i].Score,
			}
			i++
			envelopes = append(envelopes, temp)
		}else{
			temp := gin.H{
				"envelope_id": open[j].Member,
				"value":       my_redis.Rdb.HGet("envelope_money",open[j].Member.(string)).Val(),
				"opened":      true,
				"snatch_time": open[j].Score,
			}
			j++
			envelopes = append(envelopes, temp)
		}
	}
	for i<closed_len {
		temp := gin.H{
			"envelope_id": closed[i].Member,
			"opened":      false,
			"snatch_time": closed[i].Score,
		}
		i++
		envelopes = append(envelopes, temp)
	}
	for j < open_len {
		temp := gin.H{
			"envelope_id": open[j].Member,
			"value":       my_redis.Rdb.HGet("envelope_money",open[j].Member.(string)).Val(),
			"opened":      true,
			"snatch_time": open[j].Score,
		}
		j++
		envelopes = append(envelopes, temp)
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":   money,
			"envelope_list": envelopes,
		},
	})
	/*
		envelopes, totalAmount := dao.GetEnvelopesByUserId(uid)
		amountOfEnvelopes := len(envelopes)

		// 各种特殊情况应该返回什么状态码
		if amountOfEnvelopes == 0 {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "Amount of envelope is 0",
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
	*/
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
