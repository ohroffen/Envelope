package api

import (
	"MyEnvelope/dao"
	"github.com/gin-gonic/gin"
	"log"
)

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelopeId, _ := c.GetPostForm("envelope_id")
	log.Printf("envelope %s opened by %s", envelopeId, uid)

	// logic start
	// 直接查询到该红包ID，然后返回
	envelope := dao.GetEnvelopeByUserIdAndEnvelopeId(uid, envelopeId)

	value := envelope.Value

	// 修改红包状态为打开状态
	envelope.Opened = true
	dao.UpdateOpenState(&envelope)
	// logic end

	if value > 0 {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{
				"value": value,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Given user don't have such envelope",
			"data": gin.H{
				"value": value,
			},
		})
	}

}
