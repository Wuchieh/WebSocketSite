package Server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("adminPage", "templates/base.html", "templates/adminTemplates/adminBase.html")
	r.AddFromFiles("adminSetting", "templates/adminTemplates/adminSetting.html")

	r.AddFromFilesFuncs("adminIndex", template.FuncMap{
		"timeFormat": timeFormat,
	}, "templates/adminTemplates/adminIndex.html")

	r.AddFromFiles("adminChart", "templates/adminTemplates/adminChart.html")
	return r
}
func timeFormat(t time.Time) string {
	return t.Format("2006-01-02-15-04-05")
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
	//time.Now().Format("2006-01-02-15-04-05")
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
		c.HTML(200, "adminIndex", gin.H{"userTokens": userTokens})
	case "/admin/chart":
		c.HTML(200, "adminChart", nil)
	case "/admin/setting":
		c.HTML(200, "adminSetting", gin.H{"setting": setting})
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

func adminEditExpiredTime(c *gin.Context) {
	// 驗證是否傭有管理員權限
	session := sessions.Default(c)
	userToken := session.Get("userToken").(string)
	if userTokens[userToken] == nil || !userTokens[userToken].Admin {
		c.JSON(404, gin.H{})
		return
	}

	// 讀取接收到的信息
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(200, gin.H{"status": "false", "msg": err.Error()})
		return
	}

	// 格式化
	var s map[string]string
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		c.JSON(200, gin.H{"status": "false", "msg": err.Error()})
		return
	}
	timeNumber, err := strconv.Atoi(s["time"][:len(s["time"])-3])
	if err != nil {
		c.JSON(200, gin.H{"status": "false", "msg": err.Error()})
		return
	}

	// 轉換為時間格式
	unixTimestamp := int64(timeNumber)
	t := time.Unix(unixTimestamp, 0)

	// 開始設置
	userTokens[s["id"]].ExpiredTime = t.UTC()
	c.JSON(200, gin.H{"status": "true", "msg": t.UTC().Format("2006-01-02-15-04-05")})
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

	// 處理回傳
	if updateToken(string(bytes), true) {
		ut := userTokens[string(bytes)]
		if ut != nil {
			c.JSON(500, gin.H{"status": "false", "msg": "請重新加載頁面"})
		} else {
			c.JSON(200, gin.H{"status": "true",
				"CreateTime":  ut.CreateTime.Format("2006-01-02-15-04-05"),
				"UpdateTime":  ut.UpdateTime.Format("2006-01-02-15-04-05"),
				"ExpiredTime": ut.ExpiredTime.Format("2006-01-02-15-04-05"),
				"Token":       ut.Token})
		}
	} else {
		c.JSON(500, gin.H{"status": "false"})
	}

}

func adminRemoveTokenBtn(c *gin.Context) {
	// 驗證是否傭有管理員權限
	session := sessions.Default(c)
	userToken := session.Get("userToken").(string)
	if userTokens[userToken] == nil || !userTokens[userToken].Admin {
		c.JSON(404, gin.H{})
		return
	}

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	// 檢查是否想刪除自己
	if string(bytes) == userToken {
		c.JSON(401, gin.H{"status": "false", "msg": "不可刪除自己的Token"})
		return
	}

	// 開始進行刪除邏輯
	token := userTokens[string(bytes)]
	if group, ok := wsConnectGroups[token.Group]; ok {
		for i, u := range group {
			if u.Token == token.Token {
				err := u.Ws.Close()
				if err != nil {
					log.Println(err)
				}
				group = append(group[:i], group[i+1:]...)
			}
		}
		if len(group) < 1 {
			delete(wsConnectGroups, token.Group)
		}
	}
	delete(userTokens, string(bytes))
	c.JSON(200, gin.H{"status": "true"})
}

func adminSetSetting(c *gin.Context) {
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
		log.Println(err)
		c.JSON(401, gin.H{"status": "false", "msg": "接收到錯誤的資料"})
		return
	}

	// 格式化
	var s Setting
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{"status": "false", "msg": "Json格式化錯誤"})
		return
	}
	if s.ScheduleTime <= 0 || s.ExpiredTime <= 0 {
		c.JSON(401, gin.H{"status": "false", "msg": "接收到錯誤的值\nScheduleTime ExpiredTime只能為整數"})
		return
	}
	if s.ServerIP != setting.ServerIP || s.Port != setting.Port || s.Mode != setting.Mode {
		c.JSON(200, gin.H{"status": "true", "msg": "已成功修改，須重啟伺服器"})
	} else {
		c.JSON(200, gin.H{"status": "true", "msg": "已成功修改"})
	}
	setting = s
	bytes, err = json.MarshalIndent(setting, "", "  ")
	os.WriteFile("setting.json", bytes, 0666)
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
