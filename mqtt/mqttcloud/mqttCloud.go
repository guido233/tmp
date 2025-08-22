package mqttcloud

//连接 云端 mqtt
import (
	"encoding/json"
	"errors"
	"fmt"
	"go-app/conf"
	"go-app/data/model"
	"go-app/domains/bos"
	"go-app/domains/vos"
	"go-app/libs/constants"
	"go-app/libs/mqtt"
	"go-app/libs/uuid"
	"go-app/logger"
	"go-app/serve/productionline/web"
	"strconv"
	"sync"
	"time"
)

var (
	MQTTCloudManager *MQTTCloudsManager
)

// MQTTCloudsManager 云端mqtt
type MQTTCloudsManager struct {
	/*接收设备的 mqtt 实例*/
	deviceCli *mqtt.Client
	init      bool
	lock      sync.Mutex
}

func NewMQTTCloudManager() *MQTTCloudsManager {
	return &MQTTCloudsManager{}
}

// 初始化
func (dm *MQTTCloudsManager) Init() {
	if dm.init {
		return
	}
	InitCloudTopics()

	dm.init = true
	dm.deviceCli = nil
	dm.mqttReload()

	go dm.CheckDeviceConnection()
	//初始化mqtt
}

type Topic struct {
	topic string
	qos   byte
	cb    func(topic string, buff []byte) error
}

var cloudTopics []Topic
var (
	CloudRequestProcessFirst = true

	IotMqttFlag = true
)

func InitCloudTopics() {
	cloudTopics = []Topic{
		// 接收云平台消息
		//Topic{"/platform_command_push_request/", 0, CloudRequestProcess},
		// iotmatt
		Topic{fmt.Sprintf(constants.TopicCommand, conf.Conf.IotMqtt.GatewayId), 0, IotMqtt},
	}
}

func (dm *MQTTCloudsManager) mqttReload() {
	//mqtt 服务改变
	if dm.deviceCli != nil {
		dm.deviceCli.Disconnect(0)
	}

	logger.Infof("preparing CloudsConnect!")
	dm.lock.Lock()
	defer dm.lock.Unlock()
	if conf.Conf.MqttCloud.ClientId == "" {
		conf.Conf.MqttCloud.ClientId = string(uuid.NewV4().String())
	}
	dm.deviceCli = mqtt.Connect(&mqtt.Configuration{
		ClientId:               conf.Conf.MqttCloud.ClientId,
		UserName:               conf.Conf.MqttCloud.UserName,
		Password:               conf.Conf.MqttCloud.PassWord,
		BrokerAddr:             conf.Conf.MqttCloud.Host,
		BrokerPort:             conf.Conf.MqttCloud.Port,
		Timeout:                10,
		DefaultCallback:        dm.deviceCmdRequestProcess,
		ConnectedCallback:      dm.deviceConnected,
		ConnectionLostCallback: dm.deviceConnectionLost,
	}, false, false, false)

	if dm.deviceCli != nil {
		logger.Infof("CloudsConnect success!")
	} else {
		logger.Errorf("CloudsConnect error!")
		logger.Fatalf("CloudsConnect error!")
	}

}

func (dm *MQTTCloudsManager) deviceConnected(clientId string) error {
	logger.Infof("cbCloudsConnected!clientId = %v", clientId)

	// 订阅topic
	logger.Infof("cbCloudConnected!clientId = %v,cloudTopics len = %v", clientId, len(cloudTopics))
	dm.lock.Lock()
	defer dm.lock.Unlock()
	// 订阅topic
	if dm.deviceCli != nil {
		for _, t := range cloudTopics {
			go func(topic string, qos byte, cb func(topic string, buff []byte) error) {
				logger.Infof("cbCloudConnected:Subscribing topic:%v!", topic)

				for {
					if dm.deviceCli != nil && dm.deviceCli.IsConnected() {
						err := dm.deviceCli.Subscribe(topic, qos, cb)
						if err != nil {
							logger.Errorf("cbCloudConnected:Subscribe topic:%v Failed!", topic)

							time.Sleep(time.Second * time.Duration(10))
							continue
						} else {
							logger.Infof("cbCloudConnected:Subscribe topic:%v success!", topic)
						}
					}
					break
				}
			}(t.topic, t.qos, t.cb)
		}
	}

	return nil
}

func (dm *MQTTCloudsManager) deviceConnectionLost(clientId string) error {
	logger.Infof("cbCloudsConnectionLost!clientId = %v", clientId)

	// 尝试重新连接
	go func() {
		for {
			// 先等待一段时间再尝试重新连接
			time.Sleep(time.Second * 5)

			logger.Infof("Reconnecting to cloud MQTT...")

			// 重新加载并尝试连接
			dm.mqttReload()

			// 如果连接成功，则退出循环
			if dm.deviceCli != nil && dm.deviceCli.IsConnected() {
				logger.Infof("Reconnected successfully!")
				break
			}

			logger.Warnf("Reconnection attempt failed, retrying...")
		}
	}()

	return nil
}

func (dm *MQTTCloudsManager) CheckDeviceConnection() {

	connected := dm.deviceCli != nil && dm.deviceCli.IsConnected()
	if !connected {
		dm.mqttReload()
		logger.Infof("CheckCloudsConnection:CloudReconnect")
		logger.Fatalf("mqtt云端连接异常，退出程序。")
	}
	time.Sleep(time.Second * time.Duration(5))
}

/*注册TOPIC的回调*/
func (dm *MQTTCloudsManager) deviceCmdRequestProcess(topic string, msg []byte) error {
	return nil
}

// CloudRequestProcess 处理云平台消息
func CloudRequestProcess(topic string, msg []byte) error {

	// 第一次不处理
	if CloudRequestProcessFirst {
		CloudRequestProcessFirst = false
		return nil
	}

	MQTTCloudManager.lock.Lock()
	defer MQTTCloudManager.lock.Unlock()

	var req model.MQTTCommandRequest

	if err := json.Unmarshal(msg, &req); err != nil {
		logger.Errorf("CloudCmdRequestProcess:Unmarshal error! err = %v,msg = %v", err, string(msg))
		return err
	}

	if req.Serial == "CX001" {
		err := web.ProductionLineDealCmd(req.Event)
		if err != nil {
			logger.Errorf("CloudCmdRequestProcess:ProductionLineDealCmd error! err = %v", err)
			return err
		}
	}

	return nil
}

func IotMqtt(topic string, msg []byte) error {

	MQTTCloudManager.lock.Lock()
	defer MQTTCloudManager.lock.Unlock()

	var req bos.IotMqttCommand

	if err := json.Unmarshal(msg, &req); err != nil {
		logger.Errorf("CloudCmdRequestProcess:Unmarshal error! err = %v,msg = %v", err, string(msg))
		return err
	}

	fb := bos.IotMqttFeedBack{
		Et:   time.Now().Format("2006-01-02 15:04:05"),
		Id:   req.Id,
		Code: "",
		Msg:  "",
		Da:   "",
	}

	logger.Infof("IotMqtt data:req = %v", req)

	if req.Tp != constants.CommandRestart && req.Tp != constants.CommandTime && req.Tp != constants.CommandReport && req.Tp != constants.CommandConf && req.Tp != constants.CommandSet {
		logger.Errorf("CloudCmdRequestProcess:Unmarshal error! err = %v,msg = %v", errors.New("Command Error"), string(msg))

		iotmqttFeedBackSend(fb, constants.FeedBackOrderError, "无法识别到指令")

		return errors.New("Command Error")
	}

	switch req.Tp {
	case "restart":
		iotmqttFeedBackSend(fb, constants.FeedBackSuccess, "")
	case "time":
		iotmqttFeedBackSend(fb, constants.FeedBackSuccess, "")
	case "report":
		iotMqttPublish()
		iotmqttFeedBackSend(fb, constants.FeedBackSuccess, "")
	case "conf":
		iotmqttFeedBackSend(fb, constants.FeedBackSuccess, "")
	case "set":
		var set []bos.IotMqttSet
		if err := json.Unmarshal(req.Da, &set); err != nil {
			logger.Errorf("CloudCmdRequestProcess:Unmarshal error! err = %v,msg = %v", err, string(msg))
			iotmqttFeedBackSend(fb, constants.FeedBackError, "解析Da异常")
			return err
		}
		logger.Infof("IotMqtt data:set = %v", set)
		if len(set) == 0 {
			iotmqttFeedBackSend(fb, constants.FeedBackError, "Da没有数据")
			return errors.New("no da")
		}
		if set[0].Id != conf.Conf.IotMqtt.GatewayId {
			iotmqttFeedBackSend(fb, constants.FeedBackError, "设备id错误")
			return errors.New("Command Error: set.Id != conf.Conf.IotMqtt.GatewayId")
		}
		var da bos.IotMqttSetSub
		if err := json.Unmarshal(set[0].Da, &da); err != nil {
			logger.Errorf("CloudCmdRequestProcess:Unmarshal error! err = %v,msg = %v", err, string(msg))
			return err
		}
		logger.Infof("IotMqtt data:da = %v", da)
		if da.QGyO8O != "" {
			atoi, err := strconv.Atoi(da.QGyO8O)
			if err != nil {
				iotmqttFeedBackSend(fb, constants.FeedBackError, "字段错误")
				return err
			}
			bos.Data.On = atoi
			//fb.Da = da.QGyO8O
		}
		if da.QGySCO != "" {
			atoi, err := strconv.Atoi(da.QGyO8O)
			if err != nil {
				iotmqttFeedBackSend(fb, constants.FeedBackError, "字段错误")
				return err
			}
			bos.Data.SBCL = atoi
			//fb.Da = da.QGySCO
		}
		iotmqttFeedBackSend(fb, constants.FeedBackSuccess, "")
	}

	logger.Infof(string(msg))

	return nil
}

func Publish2Cloud(topic string, qos byte, retained bool, data []byte) error {
	logger.Infof("Publish2Cloud:topic = %v, qos = %v, retained = %v, data = \x1b[%dm%v\x1b[0m", topic, qos, retained, constants.Cyan, string(data))

	if MQTTCloudManager.deviceCli != nil {
		return MQTTCloudManager.deviceCli.Publish(topic, qos, retained, data)
	} else {
		return errors.New("not connected!please check connection!")
	}
}

func iotmqttFeedBackSend(fb bos.IotMqttFeedBack, code string, msg string) {
	logger.Infof("start feedback")

	fb.Code = code
	fb.Msg = msg
	dataBytes, _ := json.Marshal(fb)
	logger.Infof("IotMqttFeedBack data:dataBytes = %v", string(dataBytes))
	// 上报数据
	topic := fmt.Sprintf(constants.TopicFeedback, conf.Conf.IotMqtt.GatewayId)
	err := Publish2Cloud(topic, 1, false, dataBytes)
	if err != nil {
		logger.Errorf("mqttcloud.Publish2Cloud error: %v", err)
	}
	logger.Infof("IotMqttFeedBack finish:topic = %v", topic)
}

func iotMqttPublish() {
	data := vos.IotMqttVo{
		On:   bos.Data.On,
		SBCL: bos.Data.SBCL,
	}
	for {
		// 间隔上报时间
		time.Sleep(time.Second * 1)

		logger.Infof("IotMqttPublish start:deviceid = %v", conf.Conf.IotMqtt.GatewayId)
		// 定义上报数据
		resp := vos.IotMqttGatherVo{}

		// 获取当前时间 yyyy-MM-dd HH:mm:ss
		resp.Et = time.Now().Format("2006-01-02 15:04:05")

		// 定义上报数据的DeviceData
		deviceData := vos.IotMqttGatherDeviceDataVo{}
		deviceData.Id = conf.Conf.IotMqtt.GatewayId
		deviceData.Da = data

		// 将DeviceData添加到resp.Da中
		resp.Da = append(resp.Da, deviceData)

		// 序列化数据
		dataBytes, _ := json.Marshal(resp)
		logger.Infof("IotMqttPublish data:dataBytes = %v", string(dataBytes))
		// 上报数据
		topic := fmt.Sprintf(constants.TopicGather, conf.Conf.IotMqtt.GatewayId)
		err := Publish2Cloud(topic, 1, false, dataBytes)
		if err != nil {
			logger.Errorf("mqttcloud.Publish2Cloud error: %v", err)
		}
		logger.Infof("IotMqttPublish finish:topic = %v", topic)
	}
}
