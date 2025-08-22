package modbus

import (
	"encoding/json"
	"fmt"
	"go-app/conf"
	"go-app/domains/bos"
	"go-app/libs/constants"
	"go-app/logger"
	"os"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	modbus "github.com/simonvetter/modbus"
)

var ModbusClient *modbus.ModbusClient

func InitModbusClient() {
	var err error
	ModbusClient, err = TcpModbusClient()
	if err != nil {
		logger.Fatalf("modbusClient create error!")
		os.Exit(1)
	}
}

func TcpModbusClient() (*modbus.ModbusClient, error) {
	logger.Infof("TCPModbus start")
	// for a TCP endpoint
	// (see examples/tls_client.go for TLS usage and options)
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://" + conf.Conf.Tcpmodbus.Host + ":" + strconv.Itoa(conf.Conf.Tcpmodbus.Port),
		Timeout: 1 * time.Second,
	})
	if err != nil {
		// error out if client creation failed
		logger.Fatalf("tcpmodbus connect error!")
		return nil, err
	}
	err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.LOW_WORD_FIRST)
	if err != nil {
		logger.Fatalf("tcpmodbus SetEncoding error!")
		return nil, err
	}
	// now that the client is created and configured, attempt to connect
	err = client.Open()
	if err != nil {
		// error out if we failed to connect/open the device
		// note: multiple Open() attempts can be made on the same client until
		// the connection succeeds (i.e. err == nil), calling the constructor again
		// is unnecessary.
		// likewise, a client can be opened and closed as many times as needed.
		logger.Fatalf("tcpmodbus connect error!")
		return nil, err
	}

	logger.Infof("tcpmodbus connect to " + conf.Conf.Tcpmodbus.Host + " successful")

	return client, nil
}

/* 使用modbus slave模拟器测试 */

func ReadTcpModbus(client MQTT.Client, config conf.Config, modbusclient *modbus.ModbusClient) {
	logger.Infof("read tcpmodbus start")
	for {
		time.Sleep(time.Second * time.Duration(config.Tcpmodbus.Interval))
		for _, item := range config.Tcpmodbus.Devices {
			if item.Register == "holding" {
				var sendmsg string
				for _, read := range item.RegisterTable {
					var val string
					switch read.Type {
					case "int":
						results, err := modbusclient.ReadRegister(uint16(read.StartAddr), modbus.HOLDING_REGISTER)
						if err != nil {
							logger.Errorf("read holding register int error: %v", err)
						}
						val = strconv.Itoa(int(results))
						sendmsg = "{\"key\":" + "\"" + read.Name + "\"," + "\"val\":" + val + "}"
					case "float":
						results, err := modbusclient.ReadFloat32(uint16(read.StartAddr), modbus.HOLDING_REGISTER)
						if err != nil {
							fmt.Println(err.Error())
						}
						val = strconv.FormatFloat(float64(results), 'f', 2, 64)
						sendmsg = "{\"key\":" + "\"" + read.Name + "\"," + "\"val\":" + val + "}"
					}
					publish := client.Publish(item.Topic, 1, false, sendmsg)
					publish.Wait()
					logger.Infof("send message on topic: %s ; Message: \x1b[%dm%s\x1b[0m", item.Topic, constants.Cyan, sendmsg)
				}
			}
		}
	}
}

// 读取mqtt信息并根据配置文件写入到指定的modbus地址中
func WriteTcpModbus(client MQTT.Client, config conf.Config, modbusclient *modbus.ModbusClient) {
	logger.Infof("write tcpmodbus start")
	for {
		time.Sleep(time.Second * time.Duration(config.Tcpmodbus.Interval))
		// 读取mqtt消息
		for _, sub := range config.MqttCloud.SubList {
			token := client.Subscribe(sub, 1, func(client MQTT.Client, msg MQTT.Message) {
				logger.Infof("Received message on topic: %s\nMessage: \n%s", msg.Topic(), msg.Payload())
				enable := bos.CmdBo{}
				// 根据配置文件写入到modbus地址中
				for _, item := range config.Tcpmodbus.Devices {
					switch item.Register {
					case "write":
						if item.Topic == msg.Topic() {
							err := json.Unmarshal(msg.Payload(), &enable)
							if err != nil {
								logger.Errorf("enable json unmarshal error: %v", err)
								return
							}
							if enable.Enable == nil {
								logger.Errorf("enable json unmarshal nil")
								return
							}
							logger.Infof("enable: %d", *enable.Enable)
							for _, write := range item.RegisterTable {
								switch *enable.Enable {
								case 1: // 启动
									if write.Name == "start" {
										err := startTcp(modbusclient, uint16(write.StartAddr))
										if err != nil {
											logger.Errorf("start cgxi error: %v", err)
											return
										}
										logger.Infof("----------start cgxi successful----------")
									}
								case 0: // 停止
									if write.Name == "stop" {
										err := stopTcp(modbusclient, uint16(write.StartAddr))
										if err != nil {
											logger.Errorf("stop cgxi error: %v", err)
											return
										}
										logger.Infof("----------stop cgxi successful----------")
									}
								}

							}
						}
					}
				}
			})
			token.Wait()
		}
	}
}

func startTcp(modbusclient *modbus.ModbusClient, addr uint16) error {
	// 启动
	err := modbusclient.WriteRegister(addr, uint16(1))
	if err != nil {
		err := fmt.Errorf("write start register error: %v", err)
		return err
	}
	time.Sleep(time.Second * time.Duration(1))
	// 回写
	err = modbusclient.WriteRegister(addr, uint16(0))
	if err != nil {
		err := fmt.Errorf("write start register recover error: %v", err)
		return err
	}

	return nil
}

func stopTcp(modbusclient *modbus.ModbusClient, addr uint16) error {
	// 停止
	err := modbusclient.WriteRegister(addr, uint16(1))
	if err != nil {
		err := fmt.Errorf("write stop register error: %v", err)
		return err
	}
	time.Sleep(time.Second * time.Duration(1))
	// 回写
	err = modbusclient.WriteRegister(addr, uint16(0))
	if err != nil {
		err := fmt.Errorf("write stop register recover error: %v", err)
		return err
	}

	return nil
}
