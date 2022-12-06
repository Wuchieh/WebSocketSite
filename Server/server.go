package Server

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
	file, err := os.ReadFile("setting.json")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(file, &setting)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := os.ReadFile("tokens.json")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(bytes, &userTokens)
	if err != nil {
		log.Println(err)
		return
	}
	tokenRecycle()
	tokenSave()
}

func Server() error {
	// 建立 store
	store := cookie.NewStore([]byte("secret"))
	// session 的名稱會在 browser 變成 cookie 的 key
	r.Use(sessions.Sessions("mysession", store))

	r.Static("/statics", "statics")
	r.HTMLRender = createMyRender() //設定渲染模板
	initRouter()                    //設定路由
	r.StaticFile("/favicon.ico", "statics/imgs/websocket-logo-1280x720.png")
	err := r.Run(setting.ServerIP + ":" + setting.Port)
	return err
}
