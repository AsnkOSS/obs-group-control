package ui

import "obs-group-control/internal/state"

// Lang identifies a UI language.
type Lang int

const (
	LangZH Lang = iota
	LangEN
)

type texts struct {
	windowTitle    string
	subtitle       string
	start          string
	stop           string
	devices        string
	failed         string
	configErrTitle string
	configErrHint  string
	quit           string
	phases         map[state.Phase]string
}

var translations = map[Lang]texts{
	LangZH: {
		windowTitle:    "OBS 群控",
		subtitle:       "一键控制所有设备的录制",
		start:          "开始录制",
		stop:           "停止录制",
		devices:        "设备",
		failed:         "失败",
		configErrTitle: "配置加载失败",
		configErrHint:  "请检查设备配置文件是否存在且格式正确:",
		quit:           "退出",
		phases: map[state.Phase]string{
			state.Idle:      "空闲",
			state.Starting:  "启动中...",
			state.Recording: "录制中",
			state.Stopping:  "停止中...",
			state.Failed:    "出错",
		},
	},
	LangEN: {
		windowTitle:    "OBS Group Control",
		subtitle:       "Control recording on every device at once",
		start:          "Start",
		stop:           "Stop",
		devices:        "Devices",
		failed:         "Failed",
		configErrTitle: "Failed to load configuration",
		configErrHint:  "Check that the device config file exists and is valid:",
		quit:           "Quit",
		phases: map[state.Phase]string{
			state.Idle:      "Idle",
			state.Starting:  "Starting...",
			state.Recording: "Recording",
			state.Stopping:  "Stopping...",
			state.Failed:    "Error",
		},
	},
}
