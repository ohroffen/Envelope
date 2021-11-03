package main

import (
	"MyEnvelope/algo"
	"MyEnvelope/api"
	"MyEnvelope/dao"
	"MyEnvelope/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	conf = pflag.StringP("config", "c", "", "config filepath")
)

func main() {
	pflag.Parse()

	if err := utils.Run(*conf); err != nil {
		panic(err)
	}
	fmt.Println(viper.AllSettings())
	dao.InitDB()
	algo.InitConfig()

	r := gin.Default()

	// router
	r.POST("/snatch", api.SnatchHandler)
	r.POST("/open", api.OpenHandler)
	r.POST("/get_wallet_list", api.WalletListHandler)
	err := r.Run(":9090")
	if err != nil {
		return
	} //监听并在127.0.0.1:9090. 上启动服务
}
