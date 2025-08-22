package conf

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go-app/logger"
)

var (
	Conf *Config
)

type MqttCloud struct {
	Host     string
	Port     int
	ClientId string
	UserName string
	PassWord string
	SubList  []string
	PubList  []string
	Qos      int
}

type MqttEdge struct {
	Host     string
	Port     int
	ClientId string
	UserName string
	PassWord string
	SubList  []string
	PubList  []string
	Qos      int
}

type RegisterTable struct {
	StartAddr int
	DataLen   int
	Type      string
	Name      string
}

type Device struct {
	Register      string // 要读取的寄存器 holding或 coil
	Topic         string
	RegisterTable []RegisterTable
}

type TcpModbus struct {
	Enable   bool
	Host     string // modbus服务器地址
	Port     int    // modbus服务器端口
	SlaveID  int    // 从机地址
	Interval int    // 读取频率
	Devices  []Device
}

type RtuModbus struct {
	Enable   bool
	Device   string // 0-232,1-485(具体根据当前设备)
	BaudRate int    // 波特率
	DataBits int    // 数据位
	Parity   string // 校验位
	StopBits int    // 停止位
	SlaveID  int    // 从机地址
	Interval int    // 读取频率
	Devices  []Device
}

type IotMqtt struct {
	IP        string
	Port      string
	Router    string
	UserName  string
	Password  string
	DeviceId  string
	Password2 string
	GatewayId string
	SimpCode  string
}

type Config struct {
	ServeName string
	Serial    string
	MqttCloud MqttCloud
	MqttEdge  MqttEdge
	Tcpmodbus TcpModbus
	Rtumodbus RtuModbus
	IotMqtt   IotMqtt
}

func NewConfig() *Config {
	return &Config{
		ServeName: "",
		Serial:    "",
		MqttCloud: MqttCloud{
			Host:     "127.0.0.1",
			Port:     1883,
			ClientId: "",
			UserName: "",
			PassWord: "",
			SubList:  nil,
			PubList:  nil,
			Qos:      0,
		},
		MqttEdge: MqttEdge{
			Host:     "127.0.0.1",
			Port:     1883,
			ClientId: "",
			UserName: "",
			PassWord: "",
			SubList:  nil,
			PubList:  nil,
			Qos:      0,
		},
		Tcpmodbus: TcpModbus{
			Enable:   true,
			Host:     "192.168.20.6",
			Port:     502,
			SlaveID:  1,
			Interval: 3,
			Devices:  nil,
		},
		Rtumodbus: RtuModbus{
			Enable:   false,
			Device:   "",
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
			SlaveID:  1,
			Interval: 3,
			Devices:  nil,
		},
		IotMqtt: IotMqtt{
			Password2: "@.47DsWP",
		},
	}
}

func InitConfig() {

	Conf = NewConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Fatalf("Config file not found")
		} else {
			fmt.Println(err.Error())
		}
	}

	err := viper.Unmarshal(Conf)
	if err != nil {
		logger.Fatalf("unable to decode into struct, %v", err)
	}

	// 打印配置信息
	configuration, _ := json.Marshal(Conf)
	logger.Infof("Using conf: %v", string(configuration))
}
