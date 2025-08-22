package extra_config

import (
	"encoding/json"
	"fmt"
	"go-app/conf"
	"go-app/domains/bos"
	"go-app/libs/httpclient"
	"go-app/logger"
	"net/http"
)

func ExtraConfig() {

	// 上研需求，获取mqtt配置信息
	if conf.Conf.ServeName == "iotmqtt" {
		err := getIotMqttConfig()
		if err != nil {
			logger.Fatalf("GetIotMqttConfig error:", err)
			return
		}
	}

	logger.Infoln("get extra_config success !")
	return
}

func getIotMqttConfig() error {

	url := "http://" + conf.Conf.IotMqtt.IP + ":" + conf.Conf.IotMqtt.Port + conf.Conf.IotMqtt.Router
	logger.Infof("GetIotMqttConfig:url = %v", url)

	// 设置请求头
	headers := make(map[string]string)
	headers["userName"] = conf.Conf.IotMqtt.UserName
	headers["password"] = conf.Conf.IotMqtt.Password

	req := bos.IotMqttGetMqttInfoReq{
		DeviceId: conf.Conf.IotMqtt.DeviceId,
		Password: conf.Conf.IotMqtt.Password2,
	}

	reqData, _ := json.Marshal(req)

	bodyBytes, err := httpclient.DoHttpRequest(http.MethodPost, url, headers, reqData)
	if err != nil {
		logger.Errorf("GetIotMqttConfig:DoHttpRequest error! err = %v", err)
		return err
	}

	logger.Infof("GetIotMqttConfig:bodyBytes = %v", string(bodyBytes))
	resp := bos.IotMqttGetMqttInfoResp{}

	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		logger.Errorf("GetIotMqttConfig:Unmarshal error! err = %v", err)
		return err
	}

	if resp.Code != "001_0000_0000" {
		err := fmt.Errorf("GetIotMqttConfig:resp.Code error! resp.Code = %v, reap.message = %v", resp.Code, resp.Message)
		return err
	}

	conf.Conf.MqttCloud.ClientId = resp.Result.ClientId
	conf.Conf.MqttCloud.Host = resp.Result.MqttHost
	conf.Conf.MqttCloud.Port = resp.Result.MqttPort
	conf.Conf.MqttCloud.UserName = resp.Result.UserName
	conf.Conf.MqttCloud.PassWord = resp.Result.Password

	return nil
}
