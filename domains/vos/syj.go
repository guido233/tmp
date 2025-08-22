package vos

type ProductionSYJ struct {
	// 运行状态
	RunState bool `json:"运行状态"`
	// 运行时长
	RunTime uint32 `json:"运行时长"`
	// 产量
	ProduceAmount uint32 `json:"产量"`
	// 上坯当前值
	SPCurrentValue uint32 `json:"上坯当前值"`
	// 翻转当前值
	FZCurrentValue uint32 `json:"翻转当前值"`
	// 升降当前值
	SJCurrentValue uint32 `json:"升降当前值"`
	// 取坯当前值
	QPCurrentValue uint32 `json:"取坯当前值"`
	// 横移转速
	HYSpeed uint32 `json:"横移转速"`
	// 翻转转速
	FZSpeed uint32 `json:"翻转转速"`
	// 升降转速
	SJSpeed uint32 `json:"升降转速"`
	// 浸釉时间
	JYTime uint32 `json:"浸釉时间"`
	// 上釉时间
	SYTime uint32 `json:"上釉时间"`
	// 皮带运行时间
	PDYXTime uint32 `json:"皮带运行时间"`
	// 负压罐压力
	FYGStress bool `json:"负压罐压力"`
	// 当前配方号
	RecipeID uint32 `json:"当前配方号"`
	// 横移距离
	HYDistance uint32 `json:"横移距离"`
	// 翻转距离
	FZDistance uint32 `json:"翻转距离"`
	// 升降距离
	SJDistance    uint32 `json:"升降距离"`
	NGHistoryTims int    `json:"NGHistoryTims"`
	// ▒~N~F▒~O▒OK次▒~U▒
	OKHistoryTims int `json:"OKHistoryTims"`
	// ▒~P▒~L次▒~U▒
	HistoryTims int `json:"HistoryTims"`
}

type ProductionSDY struct {
	// 排烟 1#机
	Exhaust1 bool `json:"排烟 1#机"`
	// 排烟 2#机
	Exhaust2 bool `json:"排烟 2#机"`
	// 助燃 1#机
	Assist1 bool `json:"助燃 1#机"`
	// 助燃 2#机
	Assist2 bool `json:"助燃 2#机"`
	// 供气 1 电磁阀
	Supply1 bool `json:"供气 1 电磁阀"`
	// 供气 2 电磁阀
	Supply2 bool `json:"供气 2 电磁阀"`
	// 急冷 1#机
	Emergency1 bool `json:"急冷 1#机"`
	// 急冷 2#机
	Emergency2 bool `json:"急冷 2#机"`
	// 余热 1#机
	WasteHeat1 bool `json:"余热 1#机"`
	// 余热 2#机
	WasteHeat2 bool `json:"余热 2#机"`
	// 尾冷
	Tail bool `json:"尾冷"`
	// 燃子
	Burner bool `json:"燃子"`
	// T1 ~ T33 温度
	T1  uint16 `json:"T1温度"`
	T2  uint16 `json:"T2温度"`
	T3  uint16 `json:"T3温度"`
	T4  uint16 `json:"T4温度"`
	T5  uint16 `json:"T5温度"`
	T6  uint16 `json:"T6温度"`
	T7  uint16 `json:"T7温度"`
	T8  uint16 `json:"T8温度"`
	T9  uint16 `json:"T9温度"`
	T10 uint16 `json:"T10温度"`
	T11 uint16 `json:"T11温度"`
	T12 uint16 `json:"T12温度"`
	T13 uint16 `json:"T13温度"`
	T14 uint16 `json:"T14温度"`
	T15 uint16 `json:"T15温度"`
	T16 uint16 `json:"T16温度"`
	T17 uint16 `json:"T17温度"`
	T18 uint16 `json:"T18温度"`
	T19 uint16 `json:"T19温度"`
	T20 uint16 `json:"T20温度"`
	T21 uint16 `json:"T21温度"`
	T22 uint16 `json:"T22温度"`
	T23 uint16 `json:"T23温度"`
	T24 uint16 `json:"T24温度"`
	T25 uint16 `json:"T25温度"`
	T26 uint16 `json:"T26温度"`
	T27 uint16 `json:"T27温度"`
	T28 uint16 `json:"T28温度"`
	T29 uint16 `json:"T29温度"`
	T30 uint16 `json:"T30温度"`
	T31 uint16 `json:"T31温度"`
	T32 uint16 `json:"T32温度"`
	T33 uint16 `json:"T33温度"`
	// T5~T26 设定值
	T5Set  uint16 `json:"T5设定值"`
	T7Set  uint16 `json:"T7设定值"`
	T9Set  uint16 `json:"T9设定值"`
	T10Set uint16 `json:"T10设定值"`
	T11Set uint16 `json:"T11设定值"`
	T12Set uint16 `json:"T12设定值"`
	T13Set uint16 `json:"T13设定值"`
	T14Set uint16 `json:"T14设定值"`
	T15Set uint16 `json:"T15设定值"`
	T16Set uint16 `json:"T16设定值"`
	T17Set uint16 `json:"T17设定值"`
	T18Set uint16 `json:"T18设定值"`
	T19Set uint16 `json:"T19设定值"`
	T20Set uint16 `json:"T20设定值"`
	T21Set uint16 `json:"T21设定值"`
	T22Set uint16 `json:"T22设定值"`
	T23Set uint16 `json:"T23设定值"`
	T24Set uint16 `json:"T24设定值"`
	T25Set uint16 `json:"T25设定值"`
	T26Set uint16 `json:"T26设定值"`
	// T5~T26 输出百分比
	T5Pct  uint16 `json:"T5输出百分比"`
	T7Pct  uint16 `json:"T7输出百分比"`
	T9Pct  uint16 `json:"T9输出百分比"`
	T10Pct uint16 `json:"T10输出百分比"`
	T11Pct uint16 `json:"T11输出百分比"`
	T12Pct uint16 `json:"T12输出百分比"`
	T13Pct uint16 `json:"T13输出百分比"`
	T14Pct uint16 `json:"T14输出百分比"`
	T15Pct uint16 `json:"T15输出百分比"`
	T16Pct uint16 `json:"T16输出百分比"`
	T17Pct uint16 `json:"T17输出百分比"`
	T18Pct uint16 `json:"T18输出百分比"`
	T19Pct uint16 `json:"T19输出百分比"`
	T20Pct uint16 `json:"T20输出百分比"`
	T21Pct uint16 `json:"T21输出百分比"`
	T22Pct uint16 `json:"T22输出百分比"`
	T23Pct uint16 `json:"T23输出百分比"`
	T24Pct uint16 `json:"T24输出百分比"`
	T25Pct uint16 `json:"T25输出百分比"`
	T26Pct uint16 `json:"T26输出百分比"`

	ExhaustFrequency      float32 `json:"排烟频率"`     // 排烟频率
	ExhaustSetFrequency   float32 `json:"排烟设定频率"` // 排烟设定频率
	AssistFrequency       float32 `json:"助燃频率"`     // 助燃频率
	AssistSetFrequency    float32 `json:"助燃设定频率"` // 助燃设定频率
	EmergencyFrequency    float32 `json:"急冷频率"`     // 急冷频率
	WasteHeatFrequency    float32 `json:"余热频率"`     // 余热频率
	WasteHeatSetFrequency float32 `json:"余热设定频率"` // 余热设定频率
	TailFrequency         float32 `json:"尾冷频率"`     // 尾冷频率
	TailSetFrequency      float32 `json:"尾冷设定频率"` // 尾冷设定频率
	BurnerFrequency       float32 `json:"烘干频率"`     // 烘干频率
	BurnerSetFrequency    float32 `json:"烘干设定频率"` // 烘干设定频率
}
