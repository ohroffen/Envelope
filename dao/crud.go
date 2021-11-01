package dao

import (
	entity "MyEnvelope/model"
	"fmt"
	"sort"
	"time"
)

// GetEnvelopesByUserId 获取当前用户的所有红包
func GetEnvelopesByUserId(uid string) ([]entity.Envelope, int64) {
	var envelopes []entity.Envelope

	envelopes = make([]entity.Envelope, 0, 5)
	db, err := GetDB()
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	defer db.Close()

	db.Where("user_id = ?", uid).Find(&envelopes)

	less := func(i, j int) bool {
		return time.Time(envelopes[i].SnatchTime).Unix() < time.Time(envelopes[j].SnatchTime).Unix()
	}
	sort.Slice(envelopes, less)
	totalAmount := int64(0)

	// 预处理返回的所有红包，将所有未开状态红包的金额置0，防止在前端显示
	// 对应entity中的json的omitempty， 为0则不序列化该字段，也不会返回到前端
	if len(envelopes) > 0 {
		// for i, enve := envelopes 不能用这种方法，因为结构体是值传递，
		// 整个复制一份到enve中，改变enve中的值并没有改变红包切片中的值
		for i, _ := range envelopes {
			totalAmount += envelopes[i].Value
			if !envelopes[i].Opened {
				envelopes[i].Value = int64(0)
			}
		}
	}

	return envelopes, totalAmount
}

func GetEnvelopeByUserIdAndEnvelopeId(uid, envId string) entity.Envelope  {
	var envelope entity.Envelope
	db, err := GetDB()
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	defer db.Close()

	db.Where("user_id = ? and envelope_id = ?",uid, envId).Find(&envelope)

	return envelope
}


func UpdateOpenState(envelope *entity.Envelope) {
	db, err := GetDB()
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	defer db.Close()

	db.Model(&envelope).
		Where("envelope_id = ?", envelope.EnvelopeID).
		Update("opened", envelope.Opened)
}

// GetCurrentCount 当前遍历可以缓存到内存
// 找到用户的红包数量，可以直接select count();
func GetCurrentCount(uid string) int {
	db, err := GetDB()
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	defer db.Close()

	var count int = 0
	db.Model(entity.Envelope{}).Where("user_id = ?", uid).Count(&count)

	return count
}

// InsertEnvelope 插入红包
func InsertEnvelope(envelope entity.Envelope) {
	db, err := GetDB()
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	defer db.Close()

	if err := db.Create(&envelope).Error; err != nil {
		fmt.Println("插入失败", err)
		return
	}
}