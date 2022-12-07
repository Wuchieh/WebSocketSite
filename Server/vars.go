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
)

type Setting struct {
	ServerIP         string `json:"serverIP"`
	Port             string `json:"port"`
	SaveTime         int    `json:"saveTime"`
	CheckExpiredTime int    `json:"checkExpiredTime"`
	Mode             int    `json:"mode"`
}

type UserToken struct {
	CreateTime  time.Time
	UpdateTime  time.Time
	ExpiredTime time.Time
	Token       string
	ws          *websocket.Conn
}
