package Server

import (
	"github.com/gin-gonic/gin"
	"time"
)

var (
	r          = gin.Default()
	setting    Setting
	userTokens = make(map[string]*UserToken)
)

type Setting struct {
	ServerIP         string `json:"serverIP"`
	Port             string `json:"port"`
	SaveTime         int    `json:"saveTime"`
	CheckExpiredTime int    `json:"checkExpiredTime"`
}

type UserToken struct {
	CreateTime  time.Time
	UpdateTime  time.Time
	ExpiredTime time.Time
	Token       string
}
