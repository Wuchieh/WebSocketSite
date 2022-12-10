package Server

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

/*func SocketHandler(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ok, group := wsLoginAuthentication(ws)
	if ok {
	} else {
		_ = ws.WriteMessage(1, []byte("登入失敗，請重新連線"))
		time.Sleep(5 * time.Second)
		ws.Close()
		return
	}
	_ = ws.WriteMessage(1, []byte("登入成功"))
	defer func(ws *websocket.Conn) {
		wsLogout(ws, group)
		err := ws.Close()
		if err != nil {
			return
		}
	}(ws)
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil || msgType == -1 {
			break
		}
		sprintf := fmt.Sprintf("Message Type: %d, Message: %s", msgType, string(msg))
		fmt.Println(sprintf)
		err = ws.WriteJSON(struct {
			Reply string `json:"reply"`
		}{
			Reply: sprintf,
		})
		if err != nil {
			log.Println(err)
		}
	}
}*/

func SocketHandler(c *gin.Context) {
	ws, err := genWebSocket(c)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(msgType, string(msg))
	}
}
