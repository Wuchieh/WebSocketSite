package Server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
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
}

func adminUpdateTokenBtn(c *gin.Context) {
	// 驗證是否傭有管理員權限
	session := sessions.Default(c)
	userToken := session.Get("userToken").(string)
	if userTokens[userToken] == nil || !userTokens[userToken].Admin {
		c.JSON(404, gin.H{})
		return
	}

	// 接收信息
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	if updateToken(string(bytes), true) {
		ut := userTokens[string(bytes)]
		c.JSON(200, gin.H{"status": "true",
			"CreateTime":  ut.CreateTime.Format("2006-01-02-15-04-05"),
			"UpdateTime":  ut.UpdateTime.Format("2006-01-02-15-04-05"),
			"ExpiredTime": ut.ExpiredTime.Format("2006-01-02-15-04-05"),
			"Token":       ut.Token})
	} else {
		c.JSON(500, gin.H{"status": "false"})
	}

}
