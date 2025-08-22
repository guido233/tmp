package cgxi

import (
	"fmt"
	"go-app/conf"
	"go-app/domains/bos"
	"go-app/domains/vos"
	"go-app/logger"
	m "go-app/modbus"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/simonvetter/modbus"
)

/*
	需要修改modbus包源码:
	在读float的时候，有高低位的需求
	mc.endianness = BIG_ENDIAN
	mc.wordOrder  = LOW_WORD_FIRST
	cgxi读取使用大端低位
*/

func CgxiDealModbus(config conf.Config, mqttClient MQTT.Client) {
	modbusClient, err := m.RtuModbusClient(config)
	if err != nil {
		logger.Errorf("modbusClient create error!")
	} else {
		go CgxiReadRtuModbus(mqttClient, config, modbusClient)
	}
}

func CgxiReadRtuModbus(client MQTT.Client, config conf.Config, modbusclient *modbus.ModbusClient) {
	logger.Infof("cgxi read rtumodbus start")
	cgxiJxbDataVo := vos.CgxiJxbDataVo{
		JxbBo: bos.JxbBo{
			Serial:     "JXB001",
			DeviceType: "42001",
		},
		Data: vos.CgxiJxbVo{},
	}
	logger.Infof("cgxiJxbDataVo: %+v", cgxiJxbDataVo)
	for {
		time.Sleep(time.Second * time.Duration(config.Rtumodbus.Interval))
		for _, item := range config.Rtumodbus.Devices {
			if item.Register == "holding" {
				var sendmsg string
				for _, read := range item.RegisterTable {
					switch read.Type {
					case "cgxi-tcp":
						results, err := modbusclient.ReadFloat32s(uint16(read.StartAddr), uint16(read.DataLen), modbus.HOLDING_REGISTER)
						if err != nil {
							logger.Errorf("××××××××××read holding register cgxi-tcp error: %v", err)
						}
						if len(results) != 6 {
							logger.Errorf("failed to read cgxi-tcp data")
							continue
						}
						cgxiJxbDataVo.Data.X = fmt.Sprintf("%.2f", results[0])
						cgxiJxbDataVo.Data.Y = fmt.Sprintf("%.2f", results[1])
						cgxiJxbDataVo.Data.Z = fmt.Sprintf("%.2f", results[2])
						cgxiJxbDataVo.Data.Rx = fmt.Sprintf("%.2f", results[3])
						cgxiJxbDataVo.Data.Ry = fmt.Sprintf("%.2f", results[4])
						cgxiJxbDataVo.Data.Rz = fmt.Sprintf("%.2f", results[5])
					case "cgxi-joint":
						results, err := modbusclient.ReadFloat32s(uint16(read.StartAddr), uint16(read.DataLen), modbus.HOLDING_REGISTER)
						if err != nil {
							logger.Errorf("××××××××××read holding register cgxi-joint error: %v", err)
						}
						if len(results) != 6 {
							logger.Errorf("failed to read cgxi-joint data")
							continue
						}
						cgxiJxbDataVo.Data.Joint1 = fmt.Sprintf("%.2f", results[0])
						cgxiJxbDataVo.Data.Joint2 = fmt.Sprintf("%.2f", results[1])
						cgxiJxbDataVo.Data.Joint3 = fmt.Sprintf("%.2f", results[2])
						cgxiJxbDataVo.Data.Joint4 = fmt.Sprintf("%.2f", results[3])
						cgxiJxbDataVo.Data.Joint5 = fmt.Sprintf("%.2f", results[4])
						cgxiJxbDataVo.Data.Joint6 = fmt.Sprintf("%.2f", results[5])
					}
				}
				sendmsg = "{\"serial\":\"" + cgxiJxbDataVo.Serial + "\",\"deviceType\":\"" + cgxiJxbDataVo.DeviceType + "\",\"data\":\"{" +
					"\\\"x\\\":" + cgxiJxbDataVo.Data.X + ",\\\"y\\\":" + cgxiJxbDataVo.Data.Y + ",\\\"z\\\":" + cgxiJxbDataVo.Data.Z + ",\\\"rx\\\":" + cgxiJxbDataVo.Data.Rx + ",\\\"ry\\\":" + cgxiJxbDataVo.Data.Ry + ",\\\"rz\\\":" + cgxiJxbDataVo.Data.Rz + "," +
					"\\\"joint1\\\":" + cgxiJxbDataVo.Data.Joint1 + ",\\\"joint2\\\":" + cgxiJxbDataVo.Data.Joint2 + ",\\\"joint3\\\":" + cgxiJxbDataVo.Data.Joint3 + ",\\\"joint4\\\":" + cgxiJxbDataVo.Data.Joint4 + ",\\\"joint5\\\":" + cgxiJxbDataVo.Data.Joint5 + ",\\\"joint6\\\":" + cgxiJxbDataVo.Data.Joint6 + "}\"}"
				client.Publish(item.Topic, 1, false, sendmsg)
			}
		}
	}
}
