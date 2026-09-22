package main

import (
	"log"

	"army/internal/app/config"
	"army/internal/app/dsn"
	"army/internal/app/handler"
	"army/internal/app/repository"
	"army/internal/pkg"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	h := handler.New(repo)

	app := pkg.NewApp(cfg, router, h)
	app.RunApp()
}
