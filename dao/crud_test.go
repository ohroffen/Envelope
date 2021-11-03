package dao

import (
	entity "MyEnvelope/model"
	"fmt"
	"testing"
	"time"

	"github.com/go-basic/uuid"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// TestGetEnvelope 按照红包ID获取红包
func TestGetEnvelopesByUserId(t *testing.T) {
	InitDB()
	uid := "12345"
	envelopes, _ := GetEnvelopesByUserId(uid)

	for _, envelope := range envelopes {
		fmt.Println(envelope)
	}

}

func TestGetCurrentCount(t *testing.T) {
	uid := "2332"
	fmt.Println(GetCurrentCount(uid))
}

func TestGetEnvelopeByEnvelopeId(t *testing.T) {
	InitDB()
	uid := "1"
	envelopeId := "c4f780d6-c0d6-cab4-6374-38f9b8f93f3e"
	envelope := GetEnvelopeByUserIdAndEnvelopeId(uid, envelopeId)
	fmt.Println(envelope)
}

func TestUpdateOpenState(t *testing.T) {
	InitDB()
	uid := "18"
	envelopeId := "3ebdaef9-7211-cfe7-e85d-7f29da6327f2"
	envelope := GetEnvelopeByUserIdAndEnvelopeId(uid, envelopeId)
	fmt.Println(envelope)
	envelope.Opened = true
	UpdateOpenState(&envelope)
}

func TestInsertEnvelope(t *testing.T) {
	InitDB()
	envelope := entity.Envelope{
		EnvelopeID: uuid.New(),
		UserID:     "18",
		Opened:     false,
		Value:      80,
		SnatchTime: entity.UnixTime(time.Now()),
	}

	InsertEnvelope(envelope)
}

// func insert()  {
// 	////db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/envelope_rains?charset=utf8mb4&parseTime=True&loc=Local")
// 	////if err!= nil{
// 	////	panic(err)
// 	////}
// 	////defer db.Close()
// 	////配置MySQL连接参数
// 	//username := "root"  //账号
// 	//password := "root" //密码
// 	//host := "127.0.0.1" //数据库地址，可以是Ip或者域名
// 	//port := 3306 //数据库端口
// 	//Dbname := "envelope_rains" //数据库名
// 	//
// 	////通过前面的数据库参数，拼接MYSQL DSN， 其实就是数据库连接串（数据源名称）
// 	////MYSQL dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
// 	////类似{username}使用花括号包着的名字都是需要替换的参数
// 	//dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
// 	//	username, password, host, port, Dbname)
// 	//
// 	////连接MYSQL
// 	//db, err := gorm.Open("mysql", dsn)
// 	db, err := GetDB()

// 	if err != nil {
// 		panic("连接数据库失败, error=" + err.Error())
// 	}

// 	defer db.Close()

// 	//创建数据行
// 	uuid := uuid.New()
// 	u := entity.Envelope{
// 		EnvelopeID: uuid,
// 		UserID:     "1",
// 		Opened: true,
// 		Value:      100,
// 		SnatchTime: entity.UnixTime(time.Now()),
// 	}

// 	//db.Create(&u1)
// 	if err := db.Create(&u).Error; err != nil {
// 		fmt.Println("插入失败", err)
// 		return
// 	}

// 	fmt.Println(uuid)
// }