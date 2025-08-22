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
			logger.Fatalf("productionSDYRead error: %v", err)
			continue
		}

		logger.Infof("modbus read points result: %+v", result)

		// 结构体转json
		resultJson, err := json.Marshal(result)
		if err != nil {
			logger.Errorf("resultJson json.Marshal error: %v", err)
			logger.Fatalf("resultJson json.Marshal error: %v", err)
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

	modbusclient.SetUnitId(31)
	startAddr := uint16(0) // 起始地址
	quantity := uint16(13) // 点位数量，从 0 到 12
	// 使用 ReadCoils 读取多个点位
	flags, err := modbusclient.ReadCoils(startAddr, quantity)
	if err != nil {
		logger.Errorf("ReadCoils error: %v", err)
	} else {
		// 按点位逐一处理结果
		result.Exhaust1 = flags[0]
		result.Exhaust2 = flags[1]
		result.Assist1 = flags[2]
		result.Assist2 = flags[3]
		result.Supply1 = flags[4]
		result.Supply2 = flags[5]
		result.Emergency1 = flags[6]
		result.Emergency2 = flags[7]
		result.WasteHeat1 = flags[8]
		result.WasteHeat2 = flags[9]
		result.Tail = flags[10]
		result.Burner = flags[11]
	}

	// 假设起始地址对应 0001，数量 33
	startAddr = uint16(1)
	quantity = uint16(33)

	// 使用 ReadRegisters 读取多个 16 位寄存器
	values, readErr := modbusclient.ReadRegisters(startAddr, quantity, modbus.HOLDING_REGISTER)
	if readErr != nil {
		// 记录错误并返回
		logger.Errorf("ReadRegisters error: %v", readErr)
		return result, readErr
	}

	// 逐一处理结果，给到结构体
	// 注意：values[0] 对应地址 1, values[1] 对应地址 2, ..., values[32] 对应地址 33
	result.T1 = values[0]
	result.T2 = values[1]
	result.T3 = values[2]
	result.T4 = values[3]
	result.T5 = values[4]
	result.T6 = values[5]
	result.T7 = values[6]
	result.T8 = values[7]
	result.T9 = values[8]
	result.T10 = values[9]
	result.T11 = values[10]
	result.T12 = values[11]
	result.T13 = values[12]
	result.T14 = values[13]
	result.T15 = values[14]
	result.T16 = values[15]
	result.T17 = values[16]
	result.T18 = values[17]
	result.T19 = values[18]
	result.T20 = values[19]
	result.T21 = values[20]
	result.T22 = values[21]
	result.T23 = values[22]
	result.T24 = values[23]
	result.T25 = values[24]
	result.T26 = values[25]
	result.T27 = values[26]
	result.T28 = values[27]
	result.T29 = values[28]
	result.T30 = values[29]
	result.T31 = values[30]
	result.T32 = values[31]
	result.T33 = values[32]

	startAddr = uint16(41) // 假设是 0041
	quantity = uint16(22)  // T5~T26 共 22 个点

	values, readErr = modbusclient.ReadRegisters(startAddr, quantity, modbus.HOLDING_REGISTER)
	if readErr != nil {
		logger.Errorf("ReadRegisters error (SettingValues): %v", readErr)
		return result, readErr
	}
	result.T5Set = values[0]
	result.T7Set = values[2]
	result.T9Set = values[4]
	result.T10Set = values[5]
	result.T11Set = values[6]
	result.T12Set = values[7]
	result.T13Set = values[8]
	result.T14Set = values[9]
	result.T15Set = values[10]
	result.T16Set = values[11]
	result.T17Set = values[12]
	result.T18Set = values[13]
	result.T19Set = values[14]
	result.T20Set = values[15]
	result.T21Set = values[16]
	result.T22Set = values[17]
	result.T23Set = values[18]
	result.T24Set = values[19]
	result.T25Set = values[20]
	result.T26Set = values[21]

	startAddr = uint16(80) // 假设是 0080
	quantity = uint16(22)  // T5~T26 同 22 个点

	values, readErr = modbusclient.ReadRegisters(startAddr, quantity, modbus.HOLDING_REGISTER)
	if readErr != nil {
		logger.Errorf("ReadRegisters error (OutputPercents): %v", readErr)
		return result, readErr
	}
	result.T5Pct = values[0]
	result.T7Pct = values[2]
	result.T9Pct = values[4]
	result.T10Pct = values[5]
	result.T11Pct = values[6]
	result.T12Pct = values[7]
	result.T13Pct = values[8]
	result.T14Pct = values[9]
	result.T15Pct = values[10]
	result.T16Pct = values[11]
	result.T17Pct = values[12]
	result.T18Pct = values[13]
	result.T19Pct = values[14]
	result.T20Pct = values[15]
	result.T21Pct = values[16]
	result.T22Pct = values[17]
	result.T23Pct = values[18]
	result.T24Pct = values[19]
	result.T25Pct = values[20]
	result.T26Pct = values[21]

	startAddr = uint16(100) // 新的起始地址
	quantity = uint16(11)   // 点位数量，从 100 到 110

	// 使用 ReadRegisters 读取多个点位
	long, err := modbusclient.ReadFloat32s(startAddr, quantity, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadUint32s error: %v", err)
		logger.Fatalf("ReadUint32s error: %v", err)
	} else {
		result.ExhaustFrequency = long[0]
		result.ExhaustSetFrequency = long[1]
		result.AssistFrequency = long[2]
		result.AssistSetFrequency = long[3]
		result.EmergencyFrequency = long[4]
		result.WasteHeatFrequency = long[5]
		result.WasteHeatSetFrequency = long[6]
		result.TailFrequency = long[7]
		result.TailSetFrequency = long[8]
		result.BurnerFrequency = long[9]
		result.BurnerSetFrequency = long[10]
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
