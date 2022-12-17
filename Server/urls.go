package Server

func initRouter() {
	r.GET("/", index)
	r.GET("/myToken", myToken)
	r.GET("/ws", SocketHandler)
	r.GET("/saveAll/:PassWorld", saveAll)
	r.GET("/readAll/:PassWorld", readAll)
	r.GET("/admin/*id", adminPage)
	r.GET("/admin", adminPage)
	api := r.Group("/api")
	api.POST("/GetNewToken", getNewToken)
	api.POST("/getContent", getContent)
	api.POST("/adminEditExpiredTime", adminEditExpiredTime)
	api.POST("/adminUpdateTokenBtn", adminUpdateTokenBtn)
	api.POST("/adminRemoveTokenBtn", adminRemoveTokenBtn)
	api.POST("/adminSetSetting", adminSetSetting)
	api.POST("/adminReboot", adminReboot)
}
