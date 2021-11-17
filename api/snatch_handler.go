package api

import (
	"MyEnvelope/entity"
	"MyEnvelope/mq"
	"MyEnvelope/redis"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
)

var node *snowflake.Node
var prob float64
var max_count int64
var amount_channel chan int64
var Terminate chan os.Signal
var noEnvelopeLeft bool

func FetchMoney() {
	amount_channel = make(chan int64, 100)
	go func() {
		for {
			select {
			case <-Terminate:
				amounts := make([]interface{}, 0)
				for amount := range amount_channel {
					amounts = append(amounts, amount)
				}
				if _, err := redis.Rdb.LPush("envelope_list", amounts...).Result(); err != nil {
					log.Fatal("[Fetcher] failed writing envelope amounts back to Redis")
				}
				return
			default:
				pipe := redis.Rdb.TxPipeline()
				amounts := pipe.LRange("envelope_list", 0, 49)
				pipe.LTrim("envelope_list", 50, -1)
				pipe.Exec()
				if len(amounts.Val()) == 0 {
					noEnvelopeLeft = true
					continue
				} else if noEnvelopeLeft {
					noEnvelopeLeft = false
				}
				for _, money := range amounts.Val() {
					amount, _ := strconv.ParseInt(money, 10, 64)
					amount_channel <- amount
				}
			}
		}
	}()
}

func RetrieveSnatchConfig() {
	// init snowflake node
	id, err := redis.Rdb.Incr("node_id").Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("node id: %v", id)
	node, _ = snowflake.NewNode(id)

	// retrieve `prob` and `max_count`
	str_prob, err := redis.Rdb.Get("prob").Result()
	prob, _ = strconv.ParseFloat(str_prob, 64)
	if err != nil {
		log.Fatal(err)
	}

	str_max_count, err := redis.Rdb.Get("max_count").Result()
	max_count, _ = strconv.ParseInt(str_max_count, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("prob: %v, max_count: %v", prob, max_count)
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

func SnatchHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	if uid == "" {
		c.JSON(200, gin.H{
			"code": 4,
			"msg":  "invalid input",
		})
		return
	}
	_, err := strconv.Atoi(uid)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 4,
			"msg":  "invalid input",
		})
		return
	}

	// logic start
	//1、判断用户是否在一定概率能抢到
	if rand.Float64() > prob {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "miss",
		})
		return
	}

	//2、在redis中校验用户是否还有剩余次数
	curCount, _ := redis.Rdb.HIncrBy("user_count", uid, int64(1)).Result()
	if curCount > int64(max_count) {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "snatch count used up",
		})
		return
	}

	//3、判断队列中是否还有剩余的红包，存在或者直接出队列
	var amount int64
	if noEnvelopeLeft && len(amount_channel) == 0 {
		//将用户多增加的抢红包数要减1
		redis.Rdb.HIncrBy("user_count", uid, -1)
		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "no more envelope",
		})
		return
	} else {
		amount = <-amount_channel
	}

	//增加用户抢红包次数
	envelopeId := int64(node.Generate())
	snatchTime := time.Now().Unix() //获得当前时间戳，单位为s
	//将红包id添加到用户的未拆set中
	envelopeInfo := EnvelopeInfo{
		SnatchTime: snatchTime,
		Money:      amount,
		Opened:     false,
	}
	_, errInfo := redis.Rdb.HSet(uid+"list", fmt.Sprint(envelopeId), envelopeInfo).Result()
	if errInfo != nil {
		log.Printf("list set error,%v", errInfo)
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
			"max_count":   max_count,
			"cur_count":   curCount,
		},
	})
}
