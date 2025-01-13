package main

import (
	"IM-Backend/configs"
	"IM-Backend/route"
	"flag"
)

var (
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "configs/conf.yaml", "config file path")
}

func main() {
	nc := configs.NewNacosConfig(flagConf)
	ncClient := configs.NewNacosClient(nc)
	var ac configs.AppConf
	ac.AddNotifyer()        //添加配置通知
	ac.InitConfig(ncClient) //初始化应用配置并开启监听

	app := route.NewApp()
	app.Run() //运行应用
}
