package handler

import "github.com/fasthttp/router"

func (h *Handler) Register(r *router.Router) {
	r.GET("/", h.getAllData)
	r.POST("/add_data", h.addData)
	r.POST("/remove_data", h.removeData)
	r.POST("/update_data", h.updateData)
}
