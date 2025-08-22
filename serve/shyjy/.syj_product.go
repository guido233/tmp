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
	lastData *vos.ProductionSYJ
	inStart  bool
)

func ProductionSYJDeal() {
	// 初始化modbus
	m.InitModbusClient()

	IotMqttPowerOn()
	go IotMqttHeartBeat()
	go ProductionSYJReadTcpModbus2Cloud()
}

func ProductionSYJReadTcpModbus2Cloud() {
	logger.Infof("productionSYJ read tcpmodbus start")

	for {
		time.Sleep(time.Millisecond * time.Duration(1000))

		// 读取数据
		result, err := productionSYJRead(m.ModbusClient)
		if err != nil {
			logger.Errorf("productionSYJRead error: %v", err)
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

func productionSYJRead(modbusclient *modbus.ModbusClient) (vos.ProductionSYJ, error) {
	var (
		result vos.ProductionSYJ
		//resultUint16 uint16
		resultUint32 uint32
		flag         bool
		err          error
	)
	/*
		寄存器读D0-D20479，寄存器读0-20479
		线圈读M0-M20479，线圈读0-20479
		线圈读X21-X23，线圈读20497-20499
	*/

	// 运行状态
	flag, err = modbusclient.ReadCoil(0)
	if err != nil {
		logger.Errorf("ReadCoil RunState error: %v", err)
	} else {
		result.RunState = flag
	}
	// 运行时长
	resultUint32, err = modbusclient.ReadUint32(47285, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister RunTime error: %v", err)
	} else {
		result.RunTime = uint32(int(resultUint32))
	}
	// 产量
	resultUint32, err = modbusclient.ReadUint32(41095, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister ProduceAmount error: %v", err)
	} else {
		result.ProduceAmount = uint32(int(resultUint32))
	}
	// 上坯当前值
	resultUint32, err = modbusclient.ReadUint32(47232, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister SPCurrentValue error: %v", err)
	} else {
		result.SPCurrentValue = uint32(int(resultUint32))
	}
	// 翻转当前值
	resultUint32, err = modbusclient.ReadUint32(47239, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister FZCurrentValue error: %v", err)
	} else {
		result.FZCurrentValue = uint32(int(resultUint32))
	}
	// 升降当前值
	resultUint32, err = modbusclient.ReadUint32(47243, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister SJCurrentValue error: %v", err)
	} else {
		result.SJCurrentValue = uint32(int(resultUint32))
	}
	// 取坯当前值
	resultUint32, err = modbusclient.ReadUint32(47251, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister QPCurrentValue error: %v", err)
	} else {
		result.QPCurrentValue = uint32(int(resultUint32))
	}
	// 横移转速
	resultUint32, err = modbusclient.ReadUint32(41211, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister HYSpeed error: %v", err)
	} else {
		result.HYSpeed = uint32(int(resultUint32))
	}
	// 翻转转速
	resultUint32, err = modbusclient.ReadUint32(41311, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister FZSpeed error: %v", err)
	} else {
		result.FZSpeed = uint32(int(resultUint32))
	}
	// 升降转速
	resultUint32, err = modbusclient.ReadUint32(41411, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister SJSpeed error: %v", err)
	} else {
		result.SJSpeed = uint32(int(resultUint32))
	}
	// 浸釉时间
	resultUint32, err = modbusclient.ReadUint32(41131, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister JYTime error: %v", err)
	} else {
		result.JYTime = uint32(int(resultUint32))
	}
	// 上釉时间
	resultUint32, err = modbusclient.ReadUint32(41133, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister SYTime error: %v", err)
	} else {
		result.SYTime = uint32(int(resultUint32))
	}
	// 皮带运行时间
	resultUint32, err = modbusclient.ReadUint32(41163, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister PDYXTime error: %v", err)
	} else {
		result.PDYXTime = uint32(int(resultUint32))
	}
	// 负压罐压力
	flag, err = modbusclient.ReadCoil(20510)
	if err != nil {
		logger.Errorf("ReadCoil FYGStress error: %v", err)
	} else {
		result.FYGStress = flag
	}
	// 当前配方号
	resultUint32, err = modbusclient.ReadUint32(42433, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister RecipeID error: %v", err)
	} else {
		result.RecipeID = uint32(int(resultUint32))
	}
	// 横移距离
	resultUint32, err = modbusclient.ReadUint32(41213, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister HYDistance error: %v", err)
	} else {
		result.HYDistance = uint32(int(resultUint32))
	}
	// 翻转距离
	resultUint32, err = modbusclient.ReadUint32(41313, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister FZDistance error: %v", err)
	} else {
		result.FZDistance = uint32(int(resultUint32))
	}
	// 升降距离
	resultUint32, err = modbusclient.ReadUint32(41413, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister SJDistance error: %v", err)
	} else {
		result.SJDistance = uint32(int(resultUint32))
	}

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
