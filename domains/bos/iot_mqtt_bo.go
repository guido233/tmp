package bos

import "encoding/json"

var Ds = 0
var Data *IotMqttData

type IotMqttData struct {
	On   int
	SBCL int
}

func NewData() *IotMqttData {
	return &IotMqttData{
		On:   1,
		SBCL: 800,
	}
}

type IotMqttGetMqttInfoReq struct {
	DeviceId string `json:"deviceId"`
	Password string `json:"password"`
}

type IotMqttGetMqttInfoResp struct {
	Code    string                       `json:"code"`
	Message string                       `json:"message"`
	Result  IotMqttGetMqttInfoRespResult `json:"result"`
}

type IotMqttGetMqttInfoRespResult struct {
	ClientId string `json:"clientId"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	MqttHost string `json:"mqttHost"`
	MqttPort int    `json:"mqttPort"`
}

type IotMqttCommand struct {
	Et string          `json:"et"`
	Id string          `json:"id"`
	Tp string          `json:"tp"`
	Da json.RawMessage `json:"da"`
}

type IotMqttFeedBack struct {
	Et   string `json:"et"`
	Id   string `json:"id"`
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Da   string `json:"da"`
}

type IotMqttSet struct {
	Id string          `json:"id"`
	Da json.RawMessage `json:"da"`
}

type IotMqttSetSub struct {
	QGyO8O string `json:"QGyO8O"`
	QGySCO string `json:"QGySCO"`
}
