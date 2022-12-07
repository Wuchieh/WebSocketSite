package Server

func initRouter() {
	r.GET("/", index)
	r.GET("/myToken", myToken)
	r.GET("/ws", SocketHandler)
	api := r.Group("/api")
	api.POST("/GetNewToken", getNewToken)
}
