package Server

func initRouter() {
	r.GET("/", index)
	api := r.Group("/api")
	api.POST("/GetNewToken", getNewToken)
}
