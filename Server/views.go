package Server

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	return r
}

func index(c *gin.Context) {
	session := sessions.Default(c)
	userToken := session.Get("userToken")
	if userToken == nil {
		userToken = RandomString(32)
		session.Set("userToken", userToken)
		genNewToken(userToken.(string))
		err := session.Save()
		if err != nil {
			log.Println(err)
			return
		}
	}
	user := userTokens[userToken.(string)]
	if user == nil {
		user = genNewToken(userToken.(string))
	}
	c.HTML(200, "index", gin.H{"title": "首頁", "token": user.Token})
}

func getNewToken(c *gin.Context) {
	session := sessions.Default(c)
	userToken := session.Get("userToken")
	statusMap := make(map[string]any)
	if userToken == nil {
		statusMap["status"] = false
		statusMap["msg"] = "更新失敗"
		c.JSON(401, statusMap)
		return
	}
	if updateToken(userToken.(string)) {
		statusMap["status"] = true
		statusMap["msg"] = "更新成功"
		c.JSON(200, statusMap)
		return
	}
	statusMap["status"] = false
	statusMap["msg"] = "每分鐘只能請求一次新的Token"
	c.JSON(401, statusMap)
}

func myToken(c *gin.Context) {
	session := sessions.Default(c)
	userToken := session.Get("userToken")
	a := userTokens[userToken.(string)]
	c.JSON(200, a)
}
