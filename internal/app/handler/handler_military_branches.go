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

func (h *Handler) GetMilitaryBranchesAPI(c *gin.Context) {
	speedText := c.Query("speed")

	var militaryBranches []ds.MilitaryBranch
	var err error

	if speedText == "" {
		militaryBranches, err = h.repository.GetMilitaryBranches()
	} else {
		speed, parseErr := strconv.Atoi(speedText)

		if parseErr != nil {
			c.JSON(400, gin.H{
				"error": "Неверное значение скорости",
			})
			return
		}

		militaryBranches, err = h.repository.GetMilitaryBranchesBySpeed(speed)
	}

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := make([]MilitaryBranchSerializer, 0, len(militaryBranches))

	for _, militaryBranch := range militaryBranches {
		result = append(
			result,
			NewMilitaryBranchSerializer(militaryBranch),
		)
	}

	c.JSON(200, result)
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
	draft, err := h.repository.GetMilitaryBranchDraft(ds.CurrentUserID)

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

	_, err := h.repository.CreateMilitaryBranch(name, ds.CurrentUserID)
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

func (h *Handler) PublishMilitaryBranchAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "Некорректный id",
		})
		return
	}

	description := c.PostForm("description")

	err = h.repository.PublishMilitaryBranch(id, description, 0, 0)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	branch, err := h.repository.GetMilitaryBranch(id)
	if err != nil {
		c.JSON(404, gin.H{
			"status":  "fail",
			"message": "Род войск не найден",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   NewMilitaryBranchSerializer(*branch),
	})
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

func (h *Handler) DeleteMilitaryBranchAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "Некорректный id",
		})
		return
	}

	if err := h.repository.DeleteMilitaryBranch(id); err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Военная ветка удалена",
	})
}

func (h *Handler) CreateMilitaryBranchAPI(c *gin.Context) {
	name := c.PostForm("name")

	if name == "" {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "Название рода войск обязательно",
		})
		return
	}

	branch, err := h.repository.CreateMilitaryBranch(name, ds.CurrentUserID)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Загрузка изображения
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		fileName, err := h.repository.UploadFile(
			file,
			header.Size,
			header.Header.Get("Content-Type"),
			".jpg",
		)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}

		branch.ImageURL = fileName

		if err := h.repository.UpdateMilitaryBranchImage(branch.ID, fileName); err != nil {
			c.JSON(500, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}
	}

	// Загрузка видео
	videoFile, videoHeader, err := c.Request.FormFile("video")
	if err == nil {
		defer videoFile.Close()

		videoFileName, err := h.repository.UploadFile(
			videoFile,
			videoHeader.Size,
			videoHeader.Header.Get("Content-Type"),
			".mp4",
		)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}

		branch.VideoURL = videoFileName

		if err := h.repository.UpdateMilitaryBranchVideo(branch.ID, videoFileName); err != nil {
			c.JSON(500, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}
	}

	c.JSON(201, gin.H{
		"status": "success",
		"data":   NewMilitaryBranchSerializer(*branch),
	})
}

func (h *Handler) GetMilitaryBranchDraftAPI(c *gin.Context) {
	branch, err := h.repository.GetMilitaryBranchDraft(ds.CurrentUserID)

	if err != nil {
		c.JSON(404, gin.H{
			"status":  "fail",
			"message": "Черновик не найден",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   NewMilitaryBranchSerializer(*branch),
	})
}

func (h *Handler) GetMilitaryFeedAPI(c *gin.Context) {
	branches, err := h.repository.GetMilitaryBranches()
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   branches,
	})
}

func (h *Handler) GetMilitaryBranchAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "Некорректный id",
		})
		return
	}

	if c.Query("next") == "true" {
		nextBranch, err := h.repository.GetNextMilitaryBranch(id)

		if err != nil {
			c.JSON(404, gin.H{
				"status":  "fail",
				"message": "Следующий род войск не найден",
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "success",
			"data":   NewMilitaryBranchSerializer(*nextBranch),
		})
		return
	}

	branch, err := h.repository.GetMilitaryBranch(id)
	if err != nil {
		c.JSON(404, gin.H{
			"status":  "fail",
			"message": "Род войск не найден",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   NewMilitaryBranchSerializer(*branch),
	})
}

func (h *Handler) LikeMilitaryBranchAPI(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "Некорректный id",
		})
		return
	}

	likeValue := c.Query("like")

	// like=   0 — отменить лайк
	if likeValue == "0" {
		if err := h.repository.DeleteLike(ds.CurrentUserID, id); err != nil {
			c.JSON(500, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Лайк отменён",
		})
		return
	}

	// like  1 — поставить лайк
	like := &ds.MilitaryBranchLike{
		UserID:           ds.CurrentUserID,
		MilitaryBranchID: id,
	}

	if err := h.repository.CreateLike(like); err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   like,
	})
}

func (h *Handler) RegisterUserAPI(c *gin.Context) {
	var user ds.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": "Некорректные данные",
		})
		return
	}

	if err := h.repository.CreateUser(&user); err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (h *Handler) RegisterHandler(r *gin.Engine) {
	// Старые HTML-маршруты
	r.GET("/military_branch", h.GetMilitaryBranches)
	r.GET("/army", h.GetArmy)
	r.GET("/new_army", h.GetNewArmy)
	// API
	r.GET("/api/military_branches", h.GetMilitaryBranchesAPI)
	r.GET("/api/military_branches/draft", h.GetMilitaryBranchDraftAPI)
	r.GET("/api/military_branches/feed", h.GetMilitaryFeedAPI)
	r.GET("/api/military_branches/:id", h.GetMilitaryBranchAPI)

	r.DELETE("/api/military_branches/:id", h.DeleteMilitaryBranchAPI)

	r.POST("/api/military_branches/:id/like", h.LikeMilitaryBranchAPI)
	r.POST("/api/military_branches", h.CreateMilitaryBranchAPI)
	r.POST("/api/users", h.RegisterUserAPI)

	r.PUT("/api/military_branches/:id", h.PublishMilitaryBranchAPI)
	// Старые HTML-действия
	r.POST("/new_army", h.CreateMilitaryBranch)
	r.POST("/publish_army", h.PublishMilitaryBranch)
	r.POST("/delete_army", h.DeleteMilitaryBranch)
}

func (h *Handler) RegisterStatic(r *gin.Engine) {
	r.Static("/static", "../Web_history_army_RIP_frontend_2026/resources")
	r.LoadHTMLGlob("../Web_history_army_RIP_frontend_2026/templates/*")
}
