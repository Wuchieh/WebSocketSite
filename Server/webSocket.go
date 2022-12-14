package Server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

func genWebSocket(c *gin.Context) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	return ws, err
}

func wsClose(ws *websocket.Conn, token *UserToken) {
	_ = ws.Close()
	wsLogout(token.Group, token)
}

func wsLogout(group string, token *UserToken) {
	for i, v := range wsConnectGroups[group] {
		if v.Token == token.Token {
			wsConnectGroups[group] = append(wsConnectGroups[group][:i], wsConnectGroups[group][i+1:]...)
			break
		}
	}

	token.Ws = nil
	token.InGroup = false
	token.Group = ""
}

func wsLogin(ws *websocket.Conn) (userToken *UserToken, b bool) {
	var strMap map[string]string
	var group string
	for i := 3; i > 0; i-- {
		err := ws.ReadJSON(&strMap)
		if err != nil {
			log.Println(err)
		}

		if groupName, ok := strMap["group"]; ok {
			group = groupName
		} else {
			_ = ws.WriteMessage(1, []byte("登入失敗 尚未輸入Group"))
			continue
		}

		if token, ok := strMap["token"]; ok { // 該user的驗證 若已經加入群組則return nil,false
			if userToken, b = tokenCheck(token); userToken.InGroup {
				b = false
				return nil, b
			}
		} else {
			if i == 1 {
				break
			}
			_ = ws.WriteMessage(1, []byte("token 錯誤\n登入失敗 還剩下 "+strconv.Itoa(i-1)+" 次嘗試的機會"))
			continue
		}

		wsConnectGroups[group] = append(wsConnectGroups[group], userToken)
		userToken.InGroup = true
		userToken.Group = group
		userToken.Ws = ws
		fmt.Println("success")
		break
	}
	return
}

func wsProcess(ws *websocket.Conn, userToken *UserToken) {
	for {
		_, msg, err := ws.ReadMessage()
		if string(msg) == "/close" {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}
		wsBroadcast(msg, userToken)
	}
}

func wsBroadcast(msg []byte, userToken *UserToken) {
	if userToken.InGroup {
		for _, tokens := range wsConnectGroups[userToken.Group] {
			err := tokens.Ws.WriteMessage(1, msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
