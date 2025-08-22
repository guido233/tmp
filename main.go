package main

import (
	"go-app/conf"
	"go-app/logger"
	"go-app/mqtt/mqttcloud"
	"go-app/serve"
	"go-app/serve/extra_config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化日志
	defer logger.Start(
		logger.LogFilePath("./log/"),
		logger.LogSize(1),
		logger.LogMaxCount(3)).
		Stop()

	logger.Infof("go-app start")

	// 初始化配置
	conf.InitConfig()
	// 获取额外配置
	extra_config.ExtraConfig()

	//初始化db
	//db.Init()

	/** mqtt **/
	// 初始化mqtt edge
	//mqttedge.MQTTEdgeManager = mqttedge.NewMQTTEdgeManager()
	//mqttedge.MQTTEdgeManager.Init()

	// 初始化mqtt cloud
	mqttcloud.MQTTCloudManager = mqttcloud.NewMQTTCloudManager()
	mqttcloud.MQTTCloudManager.Init()

	// serve
	serve.StartServe()

	// 保持存活
	c := initSignal()
	handleSignal(c)
}

// initSignal register signals handler.
func initSignal() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	return c
}

// handleSignal fetch signal from chan then do exit or reload.
func handleSignal(c chan os.Signal) {
	// Block until a signal is received.
	for {
		s := <-c
		logger.Infof("get a signal %s", s.String())
		switch s {
		case os.Interrupt:
			return
		case syscall.SIGHUP:
			// TODO reload
			//return
		default:
			return
		}
	}
}
