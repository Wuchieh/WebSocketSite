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
	ServerIP string `json:"serverIP"`
	Port     string `json:"port"`
}

type UserToken struct {
	CreateTime  time.Time
	UpdateTime  time.Time
	ExpiredTime time.Time
	token       string
}
