package Server

import (
	"encoding/json"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"os"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("adminPage", "templates/base.html", "templates/adminTemplates/adminBase.html")
	r.AddFromFiles("adminIndex", "templates/adminTemplates/adminIndex.html")
	r.AddFromFiles("adminChart", "templates/adminTemplates/adminChart.html")
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
	if a == nil {
		c.Redirect(302, "/")
		return
	}
	data := gin.H{
		"CreateTime":  a.CreateTime,
		"UpdateTime":  a.UpdateTime,
		"ExpiredTime": a.ExpiredTime,
		"Token":       a.Token,
		"InGroup":     a.InGroup,
		"Group":       a.Group,
	}
	c.JSON(200, data)
}

func SocketHandler(c *gin.Context) {
	// 建立連線
	ws, err := genWebSocket(c)
	if err != nil {
		log.Println(err)
		return
	}

	var userToken *UserToken
	var ok bool

	// 登入
	if userToken, ok = wsLogin(ws); !ok {
		return
	}

	// 連線關閉
	defer func(ws *websocket.Conn) {
		wsClose(ws, userToken)
	}(ws)

	// 開始工作
	wsProcess(ws, userToken)
}

func adminPage(c *gin.Context) {
	session := sessions.Default(c)
	u := session.Get("userToken")
	if u == nil {
		c.Redirect(302, "/")
		return
	}

	userToken := u.(string)
	user := userTokens[userToken]

	if user == nil {
		c.Redirect(302, "/")
		return
	}

	if !func() bool {
		if user.Admin == false {
			if c.Query("adminPWD") != setting.AdminPWD {
				return false
			} else {
				user.Admin = true
				c.Redirect(302, "/admin")
				return true
			}
		} else {
			return true
		}
	}() {
		c.Redirect(302, "/")
	}
	if len(c.Request.URL.Query()) != 0 {
		c.Redirect(302, "/admin")
		return
	}
	c.HTML(200, "adminPage", gin.H{"title": "管理端"})
}

func getContent(c *gin.Context) {
	readAll, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	path := string(readAll)
	switch path {
	case "/admin":
		c.HTML(200, "adminIndex", nil)
	case "/admin/chart":
		c.HTML(200, "adminChart", nil)
	default:
		log.Println("出現例外情況 path:", path)
		c.String(404, "接收到錯誤Path")
	}
}

func readAll(c *gin.Context) {
	param := c.Param("PassWorld")
	if param != setting.AdminPWD {
		c.AbortWithStatus(404)
		return
	}
	func() {
		func() {
			for i := 0; i < 2; i++ {
				file, err := os.ReadFile("setting.json")
				if err != nil {
					log.Println(err)
					err := os.WriteFile("setting.json", []byte("{\n  \"serverIP\": \"127.0.0.1\",\n  \"port\": \"8080\",\n  \"scheduleTime\": 60,\n  \"expiredTime\": 120,\n  \"mode\": 0,\n  \"adminPWD\": \"adminPWD\"\n}"), 0666)
					if err != nil {
						log.Println(err)
					}
					continue
				}
				err = json.Unmarshal(file, &setting)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}()
		func() {
			for i := 0; i < 2; i++ {
				bytes, err := os.ReadFile("tokens.json")
				if err != nil {
					log.Println(err)
					err := os.WriteFile("tokens.json", []byte("{}"), 0666)
					if err != nil {
						log.Println(err)
					}
					continue
				}
				err = json.Unmarshal(bytes, &userTokens)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}()
	}()
	c.AbortWithStatus(200)
}

func saveAll(c *gin.Context) {
	param := c.Param("PassWorld")
	if param != setting.AdminPWD {
		c.AbortWithStatus(404)
		return
	}
	ScheduleChannel <- 1
	c.AbortWithStatus(200)
}
