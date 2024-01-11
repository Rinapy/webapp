package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	getInfo := router.Group("/api/v1/")
	{
		getInfo.GET("/:UID", h.getInfoToUID)
	}
	return router
}

func (h *Handler) getInfoToUID(c *gin.Context) {

}
