package vos

type ProductionSYJ struct {
	// 运行状态
	RunState bool `json:"QKPuaA"`
	// 运行时长
	RunTime uint32 `json:"QKPlyN"`
	// 产量
	ProduceAmount uint32 `json:"QKQ0Vf"`
	// 上坯当前值
	SPCurrentValue uint32 `json:"QKQ1Dt"`
	// 翻转当前值
	FZCurrentValue uint32 `json:"QKPlF5"`
	// 升降当前值
	SJCurrentValue uint32 `json:"QKPpf9"`
	// 取坯当前值
	QPCurrentValue uint32 `json:"QKPnyy"`
	// 横移转速
	HYSpeed uint32 `json:"QKPwL7"`
	// 翻转转速
	FZSpeed uint32 `json:"QKQ0qp"`
	// 升降转速
	SJSpeed uint32 `json:"QKPzyq"`
	// 浸釉时间
	JYTime uint32 `json:"QKPgk4"`
	// 上釉时间
	SYTime uint32 `json:"QKPpDp"`
	// 皮带运行时间
	PDYXTime uint32 `json:"QKPhXJ"`
	// 负压罐压力
	FYGStress bool `json:"QKPzNc"`
	// 当前配方号
	RecipeID uint32 `json:"QKPjOl"`
	// 横移距离
	HYDistance uint32 `json:"QKPtGt"`
	// 翻转距离
	FZDistance uint32 `json:"QKPhvE"`
	// 升降距离
	SJDistance    uint32 `json:"QKPsMZ"`
	NGHistoryTims int    `json:"NGHistoryTims"`
	// ▒~N~F▒~O▒OK次▒~U▒
	OKHistoryTims int `json:"OKHistoryTims"`
	// ▒~P▒~L次▒~U▒
	HistoryTims int `json:"HistoryTims"`
}

type ProductionSDY struct {
	// 排烟 1#机
	Exhaust1 bool `json:"RIfoKo"`
	// 排烟 2#机
	Exhaust2 bool `json:"RIfr9U"`
	// 助燃 1#机
	Assist1 bool `json:"RIfpQH"`
	// 助燃 2#机
	Assist2 bool `json:"RIfpDS"`
	// 供气 1 电磁阀
	Supply1 bool `json:"RIg16t"`
	// 供气 2 电磁阀
	Supply2 bool `json:"RIfuDI"`
	// 急冷 1#机
	Emergency1 bool `json:"RIfmYv"`
	// 急冷 2#机
	Emergency2 bool `json:"RIg2az"`
	// 余热 1#机
	WasteHeat1 bool `json:"RIg3UK"`
	// 余热 2#机
	WasteHeat2 bool `json:"RIfmWc"`
	// 尾冷
	Tail bool `json:"RIfvC2"`
	// 燃子
	Burner bool `json:"RIfoFG"`
	// T1 ~ T33 温度
	T1  uint16 `json:"RIfhl5"`
	T2  uint16 `json:"RIfrYt"`
	T3  uint16 `json:"RIg0E4"`
	T4  uint16 `json:"RIfivW"`
	T5  uint16 `json:"RIfneO"`
	T6  uint16 `json:"RIg0w7"`
	T7  uint16 `json:"RIfrFn"`
	T8  uint16 `json:"RIfsgW"`
	T9  uint16 `json:"RIfzmG"`
	T10 uint16 `json:"RIftyb"`
	T11 uint16 `json:"RIg2Iy"`
	T12 uint16 `json:"RIfzwl"`
	T13 uint16 `json:"RIfo1t"`
	T14 uint16 `json:"RIfzAm"`
	T15 uint16 `json:"RIg4Ts"`
	T16 uint16 `json:"RIfza4"`
	T17 uint16 `json:"RIfhzj"`
	T18 uint16 `json:"RIfxjA"`
	T19 uint16 `json:"RIfqnq"`
	T20 uint16 `json:"RIfsXW"`
	T21 uint16 `json:"RIg0Li"`
	T22 uint16 `json:"RIfsrI"`
	T23 uint16 `json:"RIfj7N"`
	T24 uint16 `json:"RIg0TF"`
	T25 uint16 `json:"RIfxNj"`
	T26 uint16 `json:"RIg0Sp"`
	T27 uint16 `json:"RIfoAo"`
	T28 uint16 `json:"RIfr2r"`
	T29 uint16 `json:"RIfj4q"`
	T30 uint16 `json:"RIg3Vx"`
	T31 uint16 `json:"RIfnEE"`
	T32 uint16 `json:"RIg4Wg"`
	T33 uint16 `json:"RIflrc"`
	// T5~T26 设定值
	T5Set  uint16 `json:"RIfqAp"`
	T7Set  uint16 `json:"RIfxgR"`
	T9Set  uint16 `json:"RIfqNx"`
	T10Set uint16 `json:"RIfjKs"`
	T11Set uint16 `json:"RIftV6"`
	T12Set uint16 `json:"RIfmWs"`
	T13Set uint16 `json:"RIfr9S"`
	T14Set uint16 `json:"RIfsBI"`
	T15Set uint16 `json:"RIg535"`
	T16Set uint16 `json:"RIg2hj"`
	T17Set uint16 `json:"RIfsS3"`
	T18Set uint16 `json:"RIfoN6"`
	T19Set uint16 `json:"RIfnCc"`
	T20Set uint16 `json:"RIfpk1"`
	T21Set uint16 `json:"RIfyyU"`
	T22Set uint16 `json:"RIg48A"`
	T23Set uint16 `json:"RIg1nr"`
	T24Set uint16 `json:"RIfjy9"`
	T25Set uint16 `json:"RIg0PG"`
	T26Set uint16 `json:"RIfpIg"`
	// T5~T26 输出百分比
	T5Pct  uint16 `json:"RIftHk"`
	T7Pct  uint16 `json:"RIftgA"`
	T9Pct  uint16 `json:"RIg5Hg"`
	T10Pct uint16 `json:"RIfyc0"`
	T11Pct uint16 `json:"RIfmDT"`
	T12Pct uint16 `json:"RIfxXk"`
	T13Pct uint16 `json:"RIfvMy"`
	T14Pct uint16 `json:"RIfqPZ"`
	T15Pct uint16 `json:"RIfrzo"`
	T16Pct uint16 `json:"RIfqKo"`
	T17Pct uint16 `json:"RIfwY7"`
	T18Pct uint16 `json:"RIg0mW"`
	T19Pct uint16 `json:"RIfkUp"`
	T20Pct uint16 `json:"RIfiJh"`
	T21Pct uint16 `json:"RIfzQP"`
	T22Pct uint16 `json:"RIfyA9"`
	T23Pct uint16 `json:"RIg3QC"`
	T24Pct uint16 `json:"RIfz9x"`
	T25Pct uint16 `json:"RIg3Vu"`
	T26Pct uint16 `json:"RIfsxD"`

	ExhaustFrequency      float32 `json:"RIfmeU"` // 排烟频率
	ExhaustSetFrequency   float32 `json:"RIfqHE"` // 排烟设定频率
	AssistFrequency       float32 `json:"RIfpf6"` // 助燃频率
	AssistSetFrequency    float32 `json:"RIfszi"` // 助燃设定频率
	EmergencyFrequency    float32 `json:"RIfkFW"` // 急冷频率
	WasteHeatFrequency    float32 `json:"RIg37x"` // 余热频率
	WasteHeatSetFrequency float32 `json:"RIfkaR"` // 余热设定频率
	TailFrequency         float32 `json:"RIfpzD"` // 尾冷频率
	TailSetFrequency      float32 `json:"RIfq7J"` // 尾冷设定频率
	BurnerFrequency       float32 `json:"RIfuXY"` // 烘干频率
	BurnerSetFrequency    float32 `json:"RIg04t"` // 烘干设定频率
}
