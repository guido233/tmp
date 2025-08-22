package serve

import (
	"go-app/logger"
	"go-app/serve/shyjy"
)

func StartServe() {

	// serve
	//switch conf.Conf.ServeName {
	switch "shyjy" {
	//case "cgxi":
	//	cgxi.CgxiDealModbus(config, mqttClient)
	//case "test":
	//	modbus.DealModbus(config, mqttClient)
	//case "xinjie":
	//	xinjie.XinJieDealModbus(config, mqttClient)
	//case "digitalTwin":
	//	digitaltwin.DigitalTwinDealModbus(config, mqttClient)
	//case "iotmqtt":
	//	iot_mqtt.IotMqttImpl()
	//case "productionline_web":
	//	web.ProductionLineWebDeal()
	case "shyjy":
		web.ProductionSDYDeal()
	default:
		logger.Errorf("config error, no serve name")
	}

}
