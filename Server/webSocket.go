package Server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
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

func wsClose(ws *websocket.Conn) {
	_ = ws.Close()
}

func wsLogin(ws *websocket.Conn) {

}
