package api

import (
	"MyEnvelope/redis"
	"encoding/json"
	"log"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pair struct {
	Key   string
	Value EnvelopeInfo
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool {
	return p[i].Value.SnatchTime < p[j].Value.SnatchTime
}
func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func rankByValueTime(envelopes map[string]string, length int) PairList {
	pl := make(PairList, length)
	i := 0
	for k, v := range envelopes {
		envelopeInfo := EnvelopeInfo{}
		json.Unmarshal([]byte(v), &envelopeInfo)
		pl[i] = Pair{k, envelopeInfo}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func WalletListHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	log.Printf("query %s's wallet", uid)
	if uid == "" {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "invalid input",
		})
		return
	}

	_, err := strconv.Atoi(uid)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "invalid input",
		})
		return
	}

	// logic start
	// 通过用户ID查询到所有红包，将金额相加，并返回
	// 按时间排序
	//version 2
	allAmount := int64(0)
	userEnvelopeList, _ := redis.Rdb.HGetAll(uid + "list").Result()
	length := len(userEnvelopeList)
	if length == 0 {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{
				"amount":        0,
				"envelope_list": "",
			},
		})
		return
	}
	pairList := rankByValueTime(userEnvelopeList, length)
	//组装返回值
	results := make([]gin.H, 0)
	for i := 0; i < length; i++ {
		envelopeId := pairList[i].Key
		snatchTime := pairList[i].Value.SnatchTime
		money := pairList[i].Value.Money
		opened := pairList[i].Value.Opened
		if !opened {
			temp := gin.H{
				"envelope_id": envelopeId,
				"opened":      false,
				"snatch_time": snatchTime,
			}
			results = append(results, temp)
		} else {
			allAmount += money
			temp := gin.H{
				"envelope_id": envelopeId,
				"value":       money,
				"opened":      true,
				"snatch_time": snatchTime,
			}
			results = append(results, temp)
		}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        allAmount,
			"envelope_list": results,
		},
	})

}
