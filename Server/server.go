package Server

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
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
	switch setting.Mode {
	case 0: // debug
		gin.SetMode(gin.DebugMode)
	case 1: // release
		gin.SetMode(gin.ReleaseMode)
	case 2: // test
		gin.SetMode(gin.TestMode)
	}
	tokenSchedule()
}

func Server() error {
	r = gin.Default()
	err := r.SetTrustedProxies([]string{setting.ServerIP})
	if err != nil {
		return err
	}
	// 建立 store
	store := cookie.NewStore([]byte("secret"))
	// session 的名稱會在 browser 變成 cookie 的 key
	r.Use(sessions.Sessions("mysession", store))

	r.Static("/statics", "statics")
	r.HTMLRender = createMyRender() //設定渲染模板
	initRouter()                    //設定路由
	r.StaticFile("/favicon.ico", "statics/imgs/websocket-logo-1280x720.png")
	err = r.Run(setting.ServerIP + ":" + setting.Port)
	return err
}
