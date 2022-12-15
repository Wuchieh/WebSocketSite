package Server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

var (
	r               *gin.Engine
	setting         Setting
	userTokens      = make(map[string]*UserToken)
	wsConnectGroups = make(map[string][]*UserToken)
	ScheduleChannel = make(chan int, 1)
)

type Setting struct { // setting.json 結構
	ServerIP     string `json:"serverIP"`     // IP
	Port         string `json:"port"`         // Port
	ScheduleTime int    `json:"scheduleTime"` // 排程執行時間
	ExpiredTime  int    `json:"expiredTime"`  // 過期時間 單位分鐘
	Mode         int    `json:"mode"`         // 啟動模式 0 Debug, 1 Release, 2 Test
	AdminPWD     string `json:"adminPWD"`     // 管理員密碼
}

type UserToken struct {
	CreateTime  time.Time
	UpdateTime  time.Time
	ExpiredTime time.Time
	Token       string
	Ws          *websocket.Conn `json:"-"`
	InGroup     bool            `json:"-"`
	Group       string          `json:"-"`
	Admin       bool
}
