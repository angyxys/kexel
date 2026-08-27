package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/angyxys/kexel/internal/config"
	"github.com/angyxys/kexel/internal/database"
	"github.com/angyxys/kexel/internal/handler"
	"github.com/angyxys/kexel/internal/repository"
	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("error on database: %v", err)
	}

	playerRepo := repository.NewPlayerRepository(db)
	playerServ := service.NewPlayerService(playerRepo)
	playerHandl := handler.NewPlayerHandler(playerServ)

	r := gin.Default()
	r.Use(cors.Default())

	webRoute := r.Group("/web")
	webRoute.Use(func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		key := strings.TrimPrefix("Bearer ", auth)
		if key == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "unauthorized",
			})
			return
		}
		if key != "101010" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "unauthorized",
			})
			return
		}
		c.Next()
	})
	webRoute.POST("/player", playerHandl.AddPlayer)

	vrcRoute := r.Group("/vrc")
	vrcRoute.GET("/list/vip", playerHandl.Vip)
	vrcRoute.GET("/list/banned", playerHandl.Banned)
	vrcRoute.GET("/list/moderator", playerHandl.Moderator)
	vrcRoute.GET("/list/owner", playerHandl.Owner)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error %s", err.Error())
	}
}
