package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"army/internal/app/repository"
)

func GetOrders(c *gin.Context) {
	speedText := c.Query("speed")

	var orders []repository.Order

	if speedText == "" {
		allOrders := repository.GetOrders()

		for _, order := range allOrders {
			if order.Status == "опубликован" {
				orders = append(orders, order)
			}
		}
	} else {
		speed, err := strconv.Atoi(speedText)

		if err != nil {
			orders = []repository.Order{}
		} else {
			orders = repository.GetOrdersBySpeed(speed)
		}
	}

	type OrderView struct {
		Order      repository.Order
		LikesCount int
	}

	ordersView := []OrderView{}

	for _, order := range orders {
		ordersView = append(ordersView, OrderView{
			Order:      order,
			LikesCount: len(order.Likes),
		})
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"orders": ordersView,
		"speed":  speedText,
	})
}

func GetOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	order := repository.GetOrder(id)

	c.HTML(http.StatusOK, "order.html", gin.H{
		"order": order,
	})
}

func GetLenta(c *gin.Context) {
	idText := c.Param("id")

	id := 1

	if idText != "" {
		id, _ = strconv.Atoi(idText)
	}

	if c.Query("next") == "true" {
		id++

		if id > 4 {
			id = 1
		}
	}

	order := repository.GetOrder(id)

	c.HTML(http.StatusOK, "lenta.html", gin.H{
		"order":      order,
		"likesCount": len(order.Likes),
	})
}

func GetAdd(c *gin.Context) {
	draft := repository.GetDraft()

	c.HTML(http.StatusOK, "add.html", gin.H{
		"order": draft,
	})
}
