package constants

const (
	// 上电消息topic
	TopicPowerOn = "/v1/%s/p/j"
	// 采集数据topic
	TopicGather = "/v1/%s/d/j"
	// 心跳状态topic
	TopicHeartbeat = "/v1/%s/s/j"
	// 指令下发topic
	TopicCommand = "/v1/%s/c/j"
	// 执行反馈topic
	TopicFeedback = "/v1/%s/r/j"
)

const (
	// 重启
	CommandRestart = "restart"
	// 对时
	CommandTime = "time"
	// 上报
	CommandReport = "report"
	// 修改配置
	CommandConf = "conf"
	// 写值
	CommandSet = "set"
)

const (
	// 反馈成功
	FeedBackSuccess = "00"
	// 反馈异常
	FeedBackError = "10"
	// 反馈指令异常
	FeedBackOrderError = "11"
)
