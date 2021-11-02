package api

import (
	"MyEnvelope/algo"
	"MyEnvelope/dao"
	entity "MyEnvelope/model"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
)

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("snatched by %s", uid)

	// logic start
	// 当前剩余红包数，剩余红包总金额应该保存到哪里，用的时候去哪里取

	// 一定概率是否抢到
	algo.Init()

	// 获取该用户红包数，可以使用redis缓存进行优化
	currCount := dao.GetCurrentCount(uid)

	if rand.Float64() < algo.SnatchRatio ||
		algo.TotalAmountOfEnvelope == 0 ||
		algo.TotalAmountOfMoney == 0 {

		// 判断该用户红包次数是否用尽
		if currCount >= algo.MaxSnatchCount {
			c.JSON(200, gin.H{
				"code": 2,
				"msg":  "snatch count used up",
				"data": gin.H{
					"envelope_id": "",
					"max_count":   algo.MaxSnatchCount,
					"cur_count":   currCount,
				},
			})
		} else {
			// 根据random返回值判断是否还有足够金额
			value := algo.GetRandomMoney()
			if value == 0 {
				// 没有抢到红包
				c.JSON(200, gin.H{
					"code": 1,
					"msg":  "no more envelope",
					"data": gin.H{
						"envelope_id": "",
						"max_count":   algo.MaxSnatchCount,
						"cur_count":   0,
					},
				})
			} else {
				// 用得到的金额生成一个红包，写入到数据库
				envelope := entity.Envelope{
					EnvelopeID: uuid.New(),
					UserID:     uid,
					Opened:     false,
					Value:      value,
					SnatchTime: entity.UnixTime(time.Now()),
				}

				dao.InsertEnvelope(envelope)

				// logic end
				c.JSON(200, gin.H{
					"code": 0,
					"msg":  "success",
					"data": gin.H{
						"envelope_id": envelope.EnvelopeID,
						"max_count":   algo.MaxSnatchCount,
						"cur_count":   currCount + 1,
					},
				})
			}
		}
	} else { // 没有抢到红包
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "miss the envelope!",
			"data": gin.H{
				"envelope_id": "",
				"max_count":   algo.MaxSnatchCount,
				"cur_count":   currCount,
			},
		})
	}
}
