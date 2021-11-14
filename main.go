package main

import (
	"MyEnvelope/api"
	"MyEnvelope/utils"
	"github.com/gin-gonic/gin"

	"github.com/spf13/pflag"

)
var (
    conf = pflag.StringP("config", "c", "", "config filepath")
)
func main() {
    pflag.Parse()

    if err := utils.Run(*conf); err != nil {
        panic(err)
    }
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
