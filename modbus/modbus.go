package modbus

import (
	config "go-app/conf"
	"go-app/logger"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func DealModbus(config config.Config, mqttClient MQTT.Client) {
	modbusClient, err := TcpModbusClient()
	if err != nil {
		logger.Errorf("modbusClient create error!")
	} else {
		go ReadTcpModbus(mqttClient, config, modbusClient)
		go WriteTcpModbus(mqttClient, config, modbusClient)
	}
}
