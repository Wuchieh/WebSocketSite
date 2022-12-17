package Server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func initRouter() {
	r.GET("/", index)
	r.GET("/myToken", myToken)
	r.GET("/ws", SocketHandler)
	r.GET("/saveAll/:PassWorld", saveAll)
	r.GET("/readAll/:PassWorld", readAll)
	r.GET("/admin/*id", adminPage)
	r.GET("/admin", adminPage)
	api := r.Group("/api")
	api.POST("/GetNewToken", getNewToken)
	api.POST("/getContent", getContent)
	api.POST("/adminEditExpiredTime", adminEditExpiredTime)
	api.POST("/adminUpdateTokenBtn", adminUpdateTokenBtn)
	api.POST("/adminRemoveTokenBtn", adminRemoveTokenBtn)
	api.POST("/adminSetSetting", adminSetSetting)
	api.POST("/adminReboot", adminReboot)
}

func adminReboot(c *gin.Context) {
	// 驗證是否傭有管理員權限
	session := sessions.Default(c)
	userToken := session.Get("userToken").(string)
	if userTokens[userToken] == nil || !userTokens[userToken].Admin {
		c.JSON(404, gin.H{})
		return
	}
	// 重新啟動伺服器
	<-boot
	boot <- 1
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	fmt.Println("===================================================")
	fmt.Println("===================reboot==========================")
	fmt.Println("===================================================")
	if err != nil {
		log.Println(err)
		return
	}
}
