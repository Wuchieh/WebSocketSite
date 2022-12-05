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
	token := session.Get("userToken")
	if token == nil {
		token = RandomString(32)
		session.Set("userToken", token)
		err := session.Save()
		if err != nil {
			log.Println(err)
			return
		}
	}
	c.HTML(200, "index", gin.H{"title": "首頁", "token": token})
}

func getNewToken(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("userToken")
	if token == nil {
		c.JSON(200, gin.H{"status": false, "msg": "更新失敗"})
		return
	}
	session.Set("userToken", nil)
	err := session.Save()
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"status": false, "msg": "更新失敗"})
		return
	}
	c.JSON(200, gin.H{"status": true, "msg": "成功更新Token"})
}
