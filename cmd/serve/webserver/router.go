package webserver

func (ws *webServer) muxRouter() {
	ws.router.GET("/", ws.getAllData())
	ws.router.POST("/add_data", ws.addData())
	ws.router.POST("/remove_data", ws.removeData())
	ws.router.POST("/update_data", ws.updateData())
}
