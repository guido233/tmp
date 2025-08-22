package web

import (
	"encoding/json"
	"fmt"
	"go-app/data/db"
	"go-app/data/model"
	"go-app/domains/bos"
	"go-app/domains/vos"
	"go-app/libs/constants"
	"go-app/libs/utils"
	"go-app/logger"
	m "go-app/modbus"
	"go-app/mqtt/mqttedge"
	"time"

	"github.com/simonvetter/modbus"
)

var (
	lastData *vos.ProductionLineWebVo
	inStart  bool
)

func ProductionLineWebDeal() {
	// 初始化modbus
	m.InitModbusClient()
	go ProductionLineWebReadTcpModbus()
}

func ProductionLineWebReadTcpModbus() {
	logger.Infof("productionLineWen read tcpmodbus start")

	for {
		time.Sleep(time.Millisecond * time.Duration(1000))

		// 获取数据库最后一条数据
		res, err := getLastData()
		if err != nil {
			logger.Errorf("getLastData error: %v", err)
		}
		if res == nil {
			logger.Warnln("lastData is nil")
		} else {
			lastData = res
			//logger.Infof("lastData: %v", res)
		}

		// 读取数据
		result, err := productionLineWebRead(m.ModbusClient)
		if err != nil {
			logger.Errorf("productionLineWebRead error: %v", err)
			continue
		}

		logger.Infof("modbus read points result: %+v", result)

		// 是否需要发送
		needSend := false
		if lastData == nil {
			needSend = true
		} else {
			needSend = (*lastData != result)
		}

		if !needSend {
			continue
		}

		// 判断结果
		if result.Start == 1 {
			inStart = true
		}
		// 运行结束
		if inStart && result.Start == 0 {
			inStart = false
			// 保存历史数据
			result.HistoryTims++
			if result.DateResult {
				result.OKHistoryTims++
			} else {
				result.NGHistoryTims++
			}
		}

		// 结构体转json
		resultJson, err := json.Marshal(result)
		if err != nil {
			logger.Errorf("resultJson json.Marshal error: %v", err)
			continue
		}
		// 保存到数据库
		deviceData := model.Device_data{
			Serial: "CX001",
			Type:   "12003",
			Data:   string(resultJson),
			Time:   time.Now().Format("2006-01-02 15:04:05"),
		}
		flag := db.AddDeviceData(deviceData)
		if !flag {
			logger.Errorf("AddDeviceData error: %v ; data: %v", err, deviceData)
		}
		// 发送到网关侧
		sendmsg := bos.HubCloudMqttReport{
			Serial:     "CX001",
			DeviceType: "12003",
			Data:       resultJson,
		}
		// 结构体转json
		sendmsgJson, err := json.Marshal(sendmsg)
		if err != nil {
			logger.Errorf("sendmsgJson json.Marshal error: %v", err)
			continue
		}
		// 发送消息
		topic := "/api/gateway/data"
		err = mqttedge.Publish2Edge(topic, 1, false, sendmsgJson)
		if err != nil {
			logger.Errorf("mqttcloud.Publish2Cloud error: %v", err)
		}
	}
}

func productionLineWebRead(modbusclient *modbus.ModbusClient) (vos.ProductionLineWebVo, error) {
	var (
		result        vos.ProductionLineWebVo
		resultUint16  uint16
		resultUint16s []uint16
		resultInt16s  []int16
		flag          bool
		err           error
	)
	/*
		寄存器读D0-D20479，寄存器读0-20479
		线圈读M0-M20479，线圈读0-20479
		线圈读X21-X23，线圈读20497-20499
	*/

	// 获取历史运行记录
	if lastData != nil {
		result.HistoryTims = lastData.HistoryTims
		result.NGHistoryTims = lastData.NGHistoryTims
		result.OKHistoryTims = lastData.OKHistoryTims
	}

	// 上料位前后初始位 M130
	flag, err = modbusclient.ReadCoil(130)
	if err != nil {
		logger.Errorf("ReadCoil UpInitPos error: %v", err)
	} else {
		result.UpInitPos = flag
	}
	// 上料位前后取料位 M153
	flag, err = modbusclient.ReadCoil(153)
	if err != nil {
		logger.Errorf("ReadCoil UpTakePos error: %v", err)
	} else {
		result.UpTakePos = flag
	}
	// 上料位前后皮带位 M170
	flag, err = modbusclient.ReadCoil(170)
	if err != nil {
		logger.Errorf("ReadCoil UpBeltPos error: %v", err)
	} else {
		result.UpBeltPos = flag
	}
	// 传感器检测
	// X21
	flag, err = modbusclient.ReadCoil(20497)
	if err != nil {
		logger.Errorf("ReadCoil X21 error: %v", err)
	} else {
		if flag { // X21为true时
			flag, err = modbusclient.ReadCoil(200)
			if err != nil {
				logger.Errorf("ReadCoil BeltSpeed error: %v", err)
			} else {
				result.BeltSpeed = flag
			}
		}
	}
	// X22
	flag, err = modbusclient.ReadCoil(20498)
	if err != nil {
		logger.Errorf("ReadCoil X22 error: %v", err)
	} else {
		if flag { // X22为true时
			flag, err = modbusclient.ReadCoil(200)
			if err != nil {
				logger.Errorf("ReadCoil UpPickHeight error: %v", err)
			} else {
				result.OutBeltLinePos = flag
			}
		}
	}
	// X23
	flag, err = modbusclient.ReadCoil(20499)
	if err != nil {
		logger.Errorf("ReadCoil X23 error: %v", err)
	} else {
		if flag { // X23为true时
			flag, err = modbusclient.ReadCoil(200)
			if err != nil {
				logger.Errorf("ReadCoil UpMateHeight error: %v", err)
			} else {
				result.OutBeltLinePos = flag
			}
		}
	}
	// 检测工位抬起高度（上料） M201
	flag, err = modbusclient.ReadCoil(201)
	if err != nil {
		logger.Errorf("ReadCoil UpPickHeight error: %v", err)
	} else {
		result.UpPickHeight = flag
	}
	// 检测工位等料高度（上料） M202
	flag, err = modbusclient.ReadCoil(202)
	if err != nil {
		logger.Errorf("ReadCoil UpMateHeight error: %v", err)
	} else {
		result.UpMateHeight = flag
	}
	// 皮带线速度（脉冲）（出料） M204
	flag, err = modbusclient.ReadCoil(204)
	if err != nil {
		logger.Errorf("ReadCoil OutInitPos error: %v", err)
	} else {
		result.OutInitPos = flag
	}
	// 检测工位抬起高度（出料） M206
	flag, err = modbusclient.ReadCoil(206)
	if err != nil {
		logger.Errorf("ReadCoil OutTakePos error: %v", err)
	} else {
		result.OutTakePos = flag
	}
	// 检测工位等料高度（出料） M220
	flag, err = modbusclient.ReadCoil(220)
	if err != nil {
		logger.Errorf("ReadCoil OutBeltPos error: %v", err)
	} else {
		result.OutBeltPos = flag
	}
	// 启动 （1表示1层启动） D3
	resultUint16, err = modbusclient.ReadRegister(3, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister Start error: %v", err)
	} else {
		result.Start = int(resultUint16)
	}
	// 机器人开始检测标志 D18
	resultUint16, err = modbusclient.ReadRegister(18, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister RobotStart error: %v", err)
	} else {
		result.RobotStart = getBool(resultUint16)
	}
	// 机器人检测完成标志 D19
	resultUint16, err = modbusclient.ReadRegister(19, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister RobotOver error: %v", err)
	} else {
		result.RobotOver = getBool(resultUint16)
	}
	// 检测结果(OK/NG) D2
	resultUint16, err = modbusclient.ReadRegister(2, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister DateResult error: %v", err)
	} else {
		result.DateResult = getNGOK(resultUint16)
	}
	// 检测结果true为OK，false为NG
	if result.DateResult {
		// OK路线执行条件 M223
		flag, err = modbusclient.ReadCoil(223)
		if err != nil {
			logger.Errorf("ReadCoil OutOk error: %v", err)
		} else {
			result.OutOk = flag
		}
		// OK方向前进 M225
		flag, err = modbusclient.ReadCoil(225)
		if err != nil {
			logger.Errorf("ReadCoil OutOkLift error: %v", err)
		} else {
			result.OutOkLift = flag
		}
	} else {
		// NG路线执行条件 M222
		flag, err = modbusclient.ReadCoil(222)
		if err != nil {
			logger.Errorf("ReadCoil OutOk error: %v", err)
		} else {
			result.OutOk = flag
		}
		// NG方向前进 M224
		flag, err = modbusclient.ReadCoil(224)
		if err != nil {
			logger.Errorf("ReadCoil OutOkLift error: %v", err)
		} else {
			result.OutOkLift = flag
		}
	}
	// 复位（1表示复位） D0
	resultUint16, err = modbusclient.ReadRegister(0, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister Reset error: %v", err)
	} else {
		result.Reset = int(resultUint16)
	}
	// 急停（停止）（1表示急停） D4
	resultUint16, err = modbusclient.ReadRegister(4, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister Stop error: %v", err)
	} else {
		result.Stop = int(resultUint16)
	}
	// 手自动模式（1表示自动，0表示手动模式） D5
	resultUint16, err = modbusclient.ReadRegister(5, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister HandAuto error: %v", err)
	} else {
		result.HandAuto = int(resultUint16)
	}
	// 相机状态 D6
	resultUint16, err = modbusclient.ReadRegister(6, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister CameraStatus error: %v", err)
	} else {
		result.CameraStatus = int(resultUint16)
	}
	// 检测工位（1表示下降，2表示上升） D7
	resultUint16, err = modbusclient.ReadRegister(7, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister DelectPos error: %v", err)
	} else {
		result.DelectPos = int(resultUint16)
	}
	// 机械臂位置 D17
	resultUint16, err = modbusclient.ReadRegister(17, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegister RobotArmPos error: %v", err)
	} else {
		result.RobotArmPos = int(resultUint16)
	}
	// 机械臂位置姿态 D50-D61
	resultUint16s, err = modbusclient.ReadRegisters(50, 12, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegisters X/Y/Z/RX/RY/RZ error: %v", err)
	} else {
		resultInt16s = utils.Uint16sToInt16s(resultUint16s)
		if len(resultInt16s) != 12 {
			logger.Errorf("failed to read X/Y/Z/RX/RY/RZ data")
		} else {
			resultInt16s = utils.Uint16sToInt16s(resultUint16s)
			result.X = utils.TwoIntToFloat(int(resultInt16s[0]), int(resultInt16s[1]))
			result.Y = utils.TwoIntToFloat(int(resultInt16s[2]), int(resultInt16s[3]))
			result.Z = utils.TwoIntToFloat(int(resultInt16s[4]), int(resultInt16s[5]))
			result.Rx = utils.TwoIntToFloat(int(resultInt16s[6]), int(resultInt16s[7]))
			result.Ry = utils.TwoIntToFloat(int(resultInt16s[8]), int(resultInt16s[9]))
			result.Rz = utils.TwoIntToFloat(int(resultInt16s[10]), int(resultInt16s[11]))
		}
	}
	//关节角度 D62-D73
	resultUint16s, err = modbusclient.ReadRegisters(62, 12, modbus.HOLDING_REGISTER)
	if err != nil {
		logger.Errorf("ReadRegisters Joint1-6 error: %v", err)
	} else {
		resultInt16s = utils.Uint16sToInt16s(resultUint16s)
		if len(resultInt16s) != 12 {
			logger.Errorf("failed to read Joint1-6 data")
		} else {
			result.Joint1 = utils.TwoIntToFloat(int(resultInt16s[0]), int(resultInt16s[1]))
			result.Joint2 = utils.TwoIntToFloat(int(resultInt16s[2]), int(resultInt16s[3]))
			result.Joint3 = utils.TwoIntToFloat(int(resultInt16s[4]), int(resultInt16s[5]))
			result.Joint4 = utils.TwoIntToFloat(int(resultInt16s[6]), int(resultInt16s[7]))
			result.Joint5 = utils.TwoIntToFloat(int(resultInt16s[8]), int(resultInt16s[9]))
			result.Joint6 = utils.TwoIntToFloat(int(resultInt16s[10]), int(resultInt16s[11]))
		}
	}

	return result, nil
}

func getBool(result uint16) bool {
	return result != 0
}

func getNGOK(result uint16) bool {
	if result == 1 {
		return true
	}
	if result == 2 {
		return false
	}
	if lastData == nil {
		return false
	}
	return lastData.DateResult
}

func getLastData() (*vos.ProductionLineWebVo, error) {
	// 数据库获取最后一条数据
	lastDataResult, err := db.GetLastDeviceData()
	if err != nil {
		logger.Errorf("GetLastDeviceData error: %v", err)
	}
	// 数据库没有数据
	if lastDataResult == nil {
		return nil, nil
	}
	// 数据库有数据
	//logger.Infof("lastDataResult: %v", lastDataResult)
	// string转结构体
	var lastDataStruct vos.ProductionLineWebVo
	err = json.Unmarshal([]byte(lastDataResult.Data), &lastDataStruct)
	if err != nil {
		logger.Errorf("json.Unmarshal error: %v", err)
	}
	return &lastDataStruct, err
}

// ProductionLineDealCmd 处理接收到的云平台消息
func ProductionLineDealCmd(cmd string) error {

	switch cmd {
	case "start":
		err := start(m.ModbusClient)
		if err != nil {
			logger.Errorf("start error: %v", err)
			return err
		}
		logger.Infof("\x1b[%dm----------start successful----------\x1b[0m", constants.Blue)
	case "stop":
		err := stop(m.ModbusClient)
		if err != nil {
			logger.Errorf("stop error: %v", err)
			return err
		}
		logger.Infof("\x1b[%dm----------stop successful----------\x1b[0m", constants.Blue)
	case "reset":
		err := reset(m.ModbusClient)
		if err != nil {
			logger.Errorf("reset error: %v", err)
			return err
		}
		logger.Infof("\x1b[%dm----------reset successful----------\x1b[0m", constants.Blue)
	default:
		logger.Errorf("CloudCmdRequestProcess:Unmarshal error! err = %v,msg = %v", "event not found", string("msg"))
	}

	return nil
}

func start(modbusclient *modbus.ModbusClient) error {
	// 启动 D3
	err := modbusclient.WriteRegister(3, 1)
	if err != nil {
		err := fmt.Errorf("write start register error: %v", err)
		return err
	}
	return nil
}

func stop(modbusclient *modbus.ModbusClient) error {
	// 停止 D4
	err := modbusclient.WriteRegister(4, 1)
	if err != nil {
		err := fmt.Errorf("write stop register error: %v", err)
		return err
	}
	return nil
}

func reset(modbusclient *modbus.ModbusClient) error {
	// 复位 D0
	err := modbusclient.WriteRegister(0, 1)
	if err != nil {
		err := fmt.Errorf("write reset register error: %v", err)
		return err
	}
	// 将启动停止和线圈回写
	err = modbusclient.WriteRegister(3, 0)
	err = modbusclient.WriteRegister(4, 0)
	err = modbusclient.WriteCoil(120, false)
	err = modbusclient.WriteCoil(130, false)
	err = modbusclient.WriteCoil(150, false)
	err = modbusclient.WriteCoil(151, false)
	err = modbusclient.WriteCoil(152, false)
	err = modbusclient.WriteCoil(153, false)
	err = modbusclient.WriteCoil(154, false)
	err = modbusclient.WriteCoil(155, false)
	err = modbusclient.WriteCoil(156, false)
	err = modbusclient.WriteCoil(157, false)
	err = modbusclient.WriteCoil(158, false)
	err = modbusclient.WriteCoil(159, false)
	err = modbusclient.WriteCoil(160, false)
	err = modbusclient.WriteCoil(161, false)
	err = modbusclient.WriteCoil(170, false)
	err = modbusclient.WriteCoil(171, false)
	err = modbusclient.WriteCoil(200, false)
	err = modbusclient.WriteCoil(201, false)
	err = modbusclient.WriteCoil(202, false)
	err = modbusclient.WriteCoil(203, false)
	err = modbusclient.WriteCoil(204, false)
	err = modbusclient.WriteCoil(205, false)
	err = modbusclient.WriteCoil(206, false)
	err = modbusclient.WriteCoil(220, false)
	err = modbusclient.WriteCoil(221, false)
	err = modbusclient.WriteCoil(222, false)
	err = modbusclient.WriteCoil(223, false)
	err = modbusclient.WriteCoil(224, false)
	err = modbusclient.WriteCoil(225, false)
	err = modbusclient.WriteCoil(226, false)
	err = modbusclient.WriteCoil(227, false)
	err = modbusclient.WriteCoil(228, false)
	err = modbusclient.WriteCoil(229, false)
	err = modbusclient.WriteCoil(230, false)
	err = modbusclient.WriteCoil(231, false)
	err = modbusclient.WriteCoil(232, false)
	err = modbusclient.WriteCoil(233, false)
	if err != nil {
		logger.Errorf("write reset register error: %v", err)
	}

	return nil
}

func handAutoOn(modbusclient *modbus.ModbusClient) error {
	// 手自动模式（1表示自动，0表示手动模式） D5
	err := modbusclient.WriteRegister(5, 1)
	if err != nil {
		err := fmt.Errorf("write handAuto register error: %v", err)
		return err
	}
	return nil
}

func handAutoOff(modbusclient *modbus.ModbusClient) error {
	// 手自动模式（1表示自动，0表示手动模式） D5
	err := modbusclient.WriteRegister(5, 0)
	if err != nil {
		err := fmt.Errorf("write handAuto register error: %v", err)
		return err
	}
	return nil
}
