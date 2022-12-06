package Server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func initRouter() {
	r.GET("/", index)
	r.GET("/myToken", myToken)
	r.GET("/ws", SocketHandler)
	api := r.Group("/api")
	api.POST("/GetNewToken", getNewToken)
}

func SocketHandler(c *gin.Context) {
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

	defer func(ws *websocket.Conn) {
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
}
