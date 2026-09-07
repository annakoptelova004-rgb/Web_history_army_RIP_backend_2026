package api

import (
	"github.com/gin-gonic/gin"

	"army/internal/app/handler"
)

func Server() {
	r := gin.Default()

	r.Static("/static", "../../../Web_history_army_RIP_frontend_2026/resources")

	r.LoadHTMLGlob("../../../Web_history_army_RIP_frontend_2026/templates/*")

	r.GET("/hello", handler.GetOrders)

	r.GET("/lenta", handler.GetLenta)
	r.GET("/add", handler.GetAdd)
	r.GET("/lenta/:id", handler.GetLenta)

	r.Run(":8081")
}
