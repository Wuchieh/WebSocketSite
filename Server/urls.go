package Server

func initRouter() {
	r.GET("/", index)
	r.GET("/myToken", myToken)
	api := r.Group("/api")
	api.POST("/GetNewToken", getNewToken)
}
