package handler

import (
	"army/internal/app/ds"
	"army/internal/app/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository repository.MilitaryBranchRepository
}

func New(repo repository.MilitaryBranchRepository) *Handler {
	return &Handler{
		repository: repo,
	}
}
func (h *Handler) GetMilitaryBranches(c *gin.Context) {
	speedText := c.Query("speed")

	var militaryBranches []ds.MilitaryBranch
	var err error

	if speedText == "" {
		militaryBranches, err = h.repository.GetMilitaryBranches()
	} else {
		speed, parseErr := strconv.Atoi(speedText)

		if parseErr != nil {
			c.HTML(400, "military_branch.html", gin.H{
				"error": "Неверное значение скорости",
			})
			return
		}

		militaryBranches, err = h.repository.GetMilitaryBranchesBySpeed(speed)
	}

	if err != nil {
		c.HTML(500, "military_branch.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	type MilitaryBranchView struct {
		MilitaryBranch ds.MilitaryBranch
		LikesCount     int
	}

	militaryBranchesView := []MilitaryBranchView{}

	for _, militaryBranch := range militaryBranches {
		if militaryBranch.Status != "опубликован" {
			continue
		}

		likesCount, err := h.repository.GetLikesCount(militaryBranch.ID)
		if err != nil {
			c.HTML(500, "military_branch.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		militaryBranchesView = append(militaryBranchesView, MilitaryBranchView{
			MilitaryBranch: militaryBranch,
			LikesCount:     likesCount,
		})
	}
	c.HTML(200, "military_branch.html", gin.H{
		"militaryBranches": militaryBranchesView,
		"speed":            speedText,
	})
}
func (h *Handler) GetMilitaryBranchesBySpeed(c *gin.Context) {
	speedText := c.Query("speed")

	if speedText == "" {
		h.GetMilitaryBranches(c)
		return
	}

	speed, err := strconv.Atoi(speedText)
	if err != nil {
		c.HTML(400, "military_branch.html", gin.H{
			"error": "Неверное значение скорости",
		})
		return
	}

	militaryBranches, err := h.repository.GetMilitaryBranchesBySpeed(speed)
	if err != nil {
		c.HTML(500, "military_branch.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(200, "military_branch.html", gin.H{
		"militaryBranches": militaryBranches,
	})
}

func (h *Handler) GetMilitaryBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(400, "army.html", gin.H{
			"error": "Неверный идентификатор",
		})
		return
	}

	militaryBranch, err := h.repository.GetMilitaryBranch(id)
	if err != nil {
		c.HTML(404, "army.html", gin.H{
			"error": "Род войск не найден",
		})
		return
	}

	c.HTML(200, "army.html", gin.H{
		"militaryBranch": militaryBranch,
	})
}

func (h *Handler) GetArmy(c *gin.Context) {
	id := int64(1)

	idText := c.Query("id")

	if idText != "" {
		parsedID, err := strconv.ParseInt(idText, 10, 64)

		if err == nil {
			id = parsedID
		}
	}

	var militaryBranch *ds.MilitaryBranch
	var err error

	if c.Query("next") == "true" {
		militaryBranch, err = h.repository.GetNextMilitaryBranch(id)
	} else {
		militaryBranch, err = h.repository.GetMilitaryBranch(id)
	}

	if err != nil {
		c.HTML(404, "army.html", gin.H{
			"error": "Род войск не найден",
		})
		return
	}

	likesCount, err := h.repository.GetLikesCount(militaryBranch.ID)

	if err != nil {
		c.HTML(500, "army.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	militaryBranch.Terrain = "равнине"
	militaryBranch.Speed = militaryBranch.SpeedPlain

	c.HTML(200, "army.html", gin.H{
		"militaryBranch": militaryBranch,
		"likesCount":     likesCount,
	})
}

func (h *Handler) GetNewArmy(c *gin.Context) {
	draft, err := h.repository.GetMilitaryBranchDraft(1)

	if err != nil {
		c.HTML(200, "new_army.html", gin.H{
			"hasDraft": false,
		})
		return
	}

	c.HTML(200, "new_army.html", gin.H{
		"hasDraft": true,
		"draft":    draft,
	})
}
func (h *Handler) CreateMilitaryBranch(c *gin.Context) {
	name := c.PostForm("name")

	_, err := h.repository.CreateMilitaryBranch(name, 1)
	if err != nil {
		c.HTML(500, "new_army.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(302, "/new_army")
}
func (h *Handler) PublishMilitaryBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.PostForm("id"), 10, 64)
	if err != nil {
		c.HTML(400, "new_army.html", gin.H{
			"error": "Неверный идентификатор",
		})
		return
	}

	description := c.PostForm("description")

	speedPlain, err := strconv.Atoi(c.PostForm("speed_plain"))
	if err != nil {
		c.HTML(400, "new_army.html", gin.H{
			"error": "Неверная скорость на равнине",
		})
		return
	}

	speedMountains, err := strconv.Atoi(c.PostForm("speed_mountains"))
	if err != nil {
		c.HTML(400, "new_army.html", gin.H{
			"error": "Неверная скорость в горах",
		})
		return
	}

	speedForest, err := strconv.Atoi(c.PostForm("speed_forest"))
	if err != nil {
		c.HTML(400, "new_army.html", gin.H{
			"error": "Неверная скорость в лесу",
		})
		return
	}

	speedRiver, err := strconv.Atoi(c.PostForm("speed_river"))
	if err != nil {
		c.HTML(400, "new_army.html", gin.H{
			"error": "Неверная скорость у реки",
		})
		return
	}

	food, err := strconv.Atoi(c.PostForm("food"))
	if err != nil {
		c.HTML(400, "new_army.html", gin.H{
			"error": "Неверное значение питания",
		})
		return
	}

	err = h.repository.PublishMilitaryBranch(
		id,
		description,
		speedPlain,
		speedMountains,
		speedForest,
		speedRiver,
		food,
	)
	if err != nil {
		c.HTML(500, "new_army.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(302, "/military_branch")
}
func (h *Handler) DeleteMilitaryBranch(c *gin.Context) {
	id, err := strconv.ParseInt(c.PostForm("id"), 10, 64)
	if err != nil {
		c.HTML(400, "military_branch.html", gin.H{
			"error": "Неверный идентификатор",
		})
		return
	}

	err = h.repository.DeleteMilitaryBranch(id)
	if err != nil {
		c.HTML(500, "military_branch.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(302, "/military_branch")
}
func (h *Handler) RegisterHandler(r *gin.Engine) {
	r.GET("/military_branch", h.GetMilitaryBranches)
	r.GET("/army", h.GetArmy)
	r.GET("/new_army", h.GetNewArmy)

	r.POST("/new_army", h.CreateMilitaryBranch)
	r.POST("/publish_army", h.PublishMilitaryBranch)
	r.POST("/delete_army", h.DeleteMilitaryBranch)
}

func (h *Handler) RegisterStatic(r *gin.Engine) {
	r.Static("/static", "../Web_history_army_RIP_frontend_2026/resources")
	r.LoadHTMLGlob("../Web_history_army_RIP_frontend_2026/templates/*")
}
