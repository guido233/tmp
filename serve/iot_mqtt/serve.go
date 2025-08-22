package iot_mqtt

import (
	"encoding/json"
	"fmt"
	"go-app/conf"
	"go-app/domains/bos"
	"go-app/domains/vos"
	"go-app/libs/constants"
	"go-app/libs/ip"
	"go-app/logger"
	"go-app/mqtt/mqttcloud"
	"time"
)

//var (
//	// IP http get mqttinfo
//	IP = "dmp.cmsrict.com"
//	// PORT http get mqttinfo
//	PORT = "39997"
//	// Router http get mqttinfo
//	Router = "/connection/info"
//	// UserName
//	UserName = "dmpHarborCon"
//	// Password
//	Password = "dmp@123"
//	// Deviceid
//	//Deviceid = "4968596387315335168"
//	Deviceid = "5449779746298150912"
//	// Password2
//	//Password2 = "@.47DsWP"
//	Password2 = "iaX0#7Gq"
//	// gatawayId
//	//gatawayId = "4968596387315335168"
//	gatawayId = "5449779746298150912"
//	GatawayId = "5449779746298150912"
//	// simpCode 简码
//	//simpCode = "QD1GaU"
//	simpCode = "5449779746298150912"
//)

type IotMqtt struct {
	ExtraConfig bos.ExtraConfig
}

func IotMqttImpl() {
	IotMqttPowerOn()
	go IotMqttHeartBeat()
	//time.Sleep(time.Second * 1)
	go IotMqttPublish()
}

// IotMqttPowerOn 上电
func IotMqttPowerOn() {

	logger.Infof("IotMqttPowerOn start:Deviceid = %v", conf.Conf.IotMqtt.DeviceId)

	// 定义上报数据
	resp := vos.IotMqttPowerOnVo{
		Et: time.Now().Format("2006-01-02 15:04:05"),
		Ip: GetLocalIp(),
	}

	dataBytes, _ := json.Marshal(resp)
	logger.Infof("IotMqttPowerOn data:dataBytes = %v", string(dataBytes))
	// 上报数据
	topic := fmt.Sprintf(constants.TopicPowerOn, conf.Conf.IotMqtt.GatewayId)
	err := mqttcloud.Publish2Cloud(topic, 1, false, dataBytes)
	if err != nil {
		logger.Errorf("mqttcloud.Publish2Cloud error: %v", err)
	}
	logger.Infof("IotMqttPowerOn finish:topic = %v", topic)

}

func IotMqttHeartBeat() {
	for {
		logger.Infof("IotMqttHeartBeat start:Deviceid = %v", conf.Conf.IotMqtt.DeviceId)

		// 定义上报数据
		resp := vos.IotMqttHeartbeatVo{
			Et: time.Now().Format("2006-01-02 15:04:05"),
			Ip: GetLocalIp(),
		}

		// 定义上报数据的DeviceData
		deviceData := vos.IotMqttHeartbeatDeviceDataVo{}
		deviceData.Id = conf.Conf.IotMqtt.SimpCode
		deviceData.Ds = bos.Ds

		// 将DeviceData添加到resp.Da中
		resp.Da = append(resp.Da, deviceData)

		dataBytes, _ := json.Marshal(resp)
		logger.Infof("IotMqttHeartBeat data:dataBytes = %v", string(dataBytes))
		// 上报数据
		topic := fmt.Sprintf(constants.TopicHeartbeat, conf.Conf.IotMqtt.GatewayId)
		err := mqttcloud.Publish2Cloud(topic, 1, false, dataBytes)
		if err != nil {
			logger.Errorf("mqttcloud.Publish2Cloud error: %v", err)
		}
		logger.Infof("IotMqttHeartBeat finish:topic = %v", topic)

		// 间隔上报时间
		//time.Sleep(time.Minute * 1)
		time.Sleep(time.Second * 10)
	}
}

// IotMqttPublish 上报数据
func IotMqttPublish() {

	data := vos.IotMqttVo{
		On:   bos.Data.On,
		SBCL: bos.Data.SBCL,
	}
	for {
		// 间隔上报时间
		time.Sleep(time.Second * 1)

		logger.Infof("IotMqttPublish start:Deviceid = %v", conf.Conf.IotMqtt.DeviceId)
		// 定义上报数据
		resp := vos.IotMqttGatherVo{}

		// 获取当前时间 yyyy-MM-dd HH:mm:ss
		resp.Et = time.Now().Format("2006-01-02 15:04:05")

		// 定义上报数据的DeviceData
		deviceData := vos.IotMqttGatherDeviceDataVo{}
		deviceData.Id = conf.Conf.IotMqtt.SimpCode
		deviceData.Da = data

		// 将DeviceData添加到resp.Da中
		resp.Da = append(resp.Da, deviceData)

		// 序列化数据
		dataBytes, _ := json.Marshal(resp)
		logger.Infof("IotMqttPublish data:dataBytes = %v", string(dataBytes))
		// 上报数据
		topic := fmt.Sprintf(constants.TopicGather, conf.Conf.IotMqtt.GatewayId)
		err := mqttcloud.Publish2Cloud(topic, 1, false, dataBytes)
		if err != nil {
			logger.Errorf("mqttcloud.Publish2Cloud error: %v", err)
		}
		logger.Infof("IotMqttPublish finish:topic = %v", topic)
		time.Sleep(time.Second * 5)
	}
}

// 获取本地ip
func GetLocalIp() string {
	return ip.GetPrivateIP()
}
