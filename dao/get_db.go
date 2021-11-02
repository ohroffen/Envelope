package dao

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// func GetDB() (db *gorm.DB, err error){
// 	//db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/envelope_rains?charset=utf8mb4&parseTime=True&loc=Local")
// 	//if err!= nil{
// 	//	panic(err)
// 	//}
// 	//defer db.Close()
// 	//配置MySQL连接参数
// 	// username := "root"  //账号
// 	// password := "root" //密码
// 	// host := "127.0.0.1" //数据库地址，可以是Ip或者域名
// 	// port := 3306 //数据库端口
// 	// Dbname := "envelope_rains" //数据库名

// 	// //通过前面的数据库参数，拼接MYSQL DSN， 其实就是数据库连接串（数据源名称）
// 	// //MYSQL dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
// 	// //类似{username}使用花括号包着的名字都是需要替换的参数
// 	// dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
// 	// 	username, password, host, port, Dbname)
// 	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
// 		viper.GetString("db.username"), viper.GetString("db.password"), viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.name"))

// 	fmt.Println(dsn)
// 	//连接MYSQL
// 	return gorm.Open("mysql", dsn)
// }
