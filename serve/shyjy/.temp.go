package web

import (
	"encoding/json"
	"fmt"
	"go-app/conf"
	"go-app/domains/bos"
	"go-app/domains/vos"
	"go-app/libs/constants"
	"go-app/libs/ip"
	"go-app/logger"
	m "go-app/modbus"
	"go-app/mqtt/mqttcloud"
	"time"

	"github.com/simonvetter/modbus"
)

var (
	lastData *vos.ProductionSDY
	inStart  bool
)

func ProductionSDYDeal() {
	// 初始化modbus
	m.InitModbusClient()

	IotMqttPowerOn()
	go IotMqttHeartBeat()
	go ProductionSDYReadTcpModbus2Cloud()
}

func ProductionSDYReadTcpModbus2Cloud() {
	logger.Infof("productionSDY read tcpmodbus start")

	for {
		time.Sleep(time.Millisecond * time.Duration(1000))

		// 读取数据
		result, err := productionSDYRead(m.ModbusClient)
		if err != nil {
			logger.Errorf("productionSDYRead error: %v", err)
			continue
		}

		logger.Infof("modbus read points result: %+v", result)

		// 结构体转json
		resultJson, err := json.Marshal(result)
		if err != nil {
			logger.Errorf("resultJson json.Marshal error: %v", err)
			continue
		}
		// 发送到上研
		IotMqttPublish(resultJson)
	}
}

func productionSDYRead(modbusclient *modbus.ModbusClient) (vos.ProductionSDY, error) {
	var (
		result vos.ProductionSDY
		//resultUint16 uint16
	)

	// 使用 ReadRegisters 读取多个 16 位寄存器
	values, readErr := modbusclient.ReadRegister(33, modbus.HOLDING_REGISTER)
	if readErr != nil {
		// 记录错误并返回
		logger.Errorf("ReadRegisters error: %v", readErr)
		return result, readErr
	}

	// 逐一处理结果，给到结构体
	// 注意：values[0] 对应地址 1, values[1] 对应地址 2, ..., values[32] 对应地址 33

	result.T33 = values

	return result, nil
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

// IotMqttHeartBeat 心跳
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
func IotMqttPublish(resultJson []byte) {

	// 间隔上报时间
	time.Sleep(time.Second * 1)

	logger.Infof("IotMqttPublish start:Deviceid = %v", conf.Conf.IotMqtt.DeviceId)
	// 定义上报数据
	resp := vos.IotMqttGatherSY{}

	// 获取当前时间 yyyy-MM-dd HH:mm:ss
	resp.Et = time.Now().Format("2006-01-02 15:04:05")

	// 定义上报数据的DeviceData
	deviceData := vos.IotMqttGatherDeviceDataSY{}
	deviceData.Id = conf.Conf.IotMqtt.SimpCode
	err := json.Unmarshal(resultJson, &deviceData.Da)
	if err != nil {
		logger.Errorf("IotMqttPublish json.Unmarshal error: %v", err)
		return
	}

	// 将DeviceData添加到resp.Da中
	resp.Da = append(resp.Da, deviceData)

	// 序列化数据
	dataBytes, _ := json.Marshal(resp)
	logger.Infof("IotMqttPublish data:dataBytes = %v", string(dataBytes))
	// 上报数据
	topic := fmt.Sprintf(constants.TopicGather, conf.Conf.IotMqtt.GatewayId)
	err = mqttcloud.Publish2Cloud(topic, 1, false, dataBytes)
	if err != nil {
		logger.Errorf("mqttcloud.Publish2Cloud error: %v", err)
	}
	logger.Infof("IotMqttPublish finish:topic = %v", topic)
	time.Sleep(time.Second * 5)

}

// 获取本地ip
func GetLocalIp() string {
	return ip.GetPrivateIP()
}
