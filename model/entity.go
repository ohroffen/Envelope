package entity

import (
	"fmt"
	"time"
)

// Envelope 获取用户当前红包数，可以单独设置一个该用户当前红包数（需要分表）
// 或者使用select count(*)
type Envelope struct {
	EnvelopeID string   `gorm:"column:envelope_id" json:"envelope_id"`
	UserID     string   `gorm:"column:user_id" json:"-"`
	Opened     bool     `gorm:"column:opened" json:"opened"`
	Value      int64    `gorm:"column:value" json:"value,omitempty"`
	SnatchTime UnixTime `gorm:"column:snatch_time" json:"snatch_time"`
	//CurrentCount int32 `gorm:"current_count"`
}

type UnixTime time.Time

// MarshalJSON 序列化的时候将time.Time转化为UnixTime
func (t UnixTime) MarshalJSON() ([]byte, error) {
	// do your serializing here
	stamp := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(stamp), nil
}

//type Envelope struct {
//	EnvelopeID string `gorm:"column:envelope_id" json:"envelope_id"`
//	UserID string `gorm:"column:user_id" json:"-"`
//	Opened bool `gorm:"column:opened" json:"opened"`
//	Value int64 `gorm:"column:value" json:"value"`
//	SnatchTime time.Time `gorm:"column:snatch_time" json:"snatch_time"`
//	//CurrentCount int32 `gorm:"current_count"`
//}

// TableName 设置表名，可以通过给struct类型定义 TableName函数，
// 返回当前struct绑定的mysql表名是什么
func (e Envelope) TableName() string {
	//绑定MYSQL表名为users
	return "users"
}
