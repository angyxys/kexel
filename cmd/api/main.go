package main

import (
	"log"
	"net/http"

	"github.com/angyxys/kexel/internal/config"
	"github.com/angyxys/kexel/internal/database"
	"github.com/angyxys/kexel/internal/database/models"
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

	// Auto migrate tables
	if err := db.AutoMigrate(
		&models.Player{},
		&models.User{},
		&models.RefreshToken{},
		&models.InvitationCode{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("error migrating database: %v", err)
	}

	// Initialize repositories
	playerRepo := repository.NewPlayerRepository(db)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)

	// Initialize services
	playerServ := service.NewPlayerService(playerRepo)
	authServ := service.NewAuthServiceWithInvitations(userRepo, refreshTokenRepo, invitationRepo, cfg.JWT_SECRET)
	auditServ := service.NewAuditService(auditRepo)
	banServ := service.NewBanService(playerRepo)
	statsServ := service.NewStatsService(playerRepo, auditRepo, userRepo)
	invitationServ := service.NewInvitationService(invitationRepo)

	// Initialize handlers
	playerHandl := handler.NewPlayerHandler(playerServ)
	authHandl := handler.NewAuthHandler(authServ)
	auditHandl := handler.NewAuditHandler(auditServ)
	banHandl := handler.NewBanHandler(banServ)
	statsHandl := handler.NewStatsHandler(statsServ)
	invitationHandl := handler.NewInvitationHandler(invitationServ)

	r := gin.Default()
	r.Use(cors.Default())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Public auth routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandl.Register)
		authRoutes.POST("/login", authHandl.Login)
		authRoutes.POST("/refresh", authHandl.Refresh)
		authRoutes.POST("/logout", authHandl.Logout)
		authRoutes.GET("/invitations/validate", invitationHandl.ValidateInvitation)
	}

	// VRChat public routes (for in-game use)
	vrcRoute := r.Group("/vrc")
	vrcRoute.Use(secretKeyMiddleware(cfg.MAP_SECRET))
	{
		vrcRoute.GET("/list/vip", playerHandl.Vip)
		vrcRoute.GET("/list/banned", playerHandl.Banned)
		vrcRoute.GET("/list/moderator", playerHandl.Moderator)
		vrcRoute.GET("/list/owner", playerHandl.Owner)
	}

	// Web protected routes (admin panel)
	webRoute := r.Group("/web")
	webRoute.Use(authHandl.AuthMiddleware)
	{
		// Player management
		webRoute.GET("/player/:id", playerHandl.GetPlayer)
		webRoute.POST("/player", playerHandl.AddPlayer)
		webRoute.PUT("/player/:id", playerHandl.UpdatePlayer)
		webRoute.DELETE("/player/:id", playerHandl.DeletePlayer)
		webRoute.GET("/players", playerHandl.ListPlayers)
		webRoute.GET("/players/search", playerHandl.SearchPlayers)
		webRoute.GET("/players/filter", playerHandl.FilterPlayers)

		// Ban management
		webRoute.POST("/ban/:id", banHandl.BanPlayer)
		webRoute.DELETE("/ban/:id", banHandl.UnbanPlayer)
		webRoute.GET("/ban/:id", banHandl.GetBanInfo)
		webRoute.GET("/bans", banHandl.GetBannedPlayers)
		webRoute.GET("/bans/expiring-soon", banHandl.GetExpiringSoonBans)
		webRoute.POST("/bans/cleanup", banHandl.CleanupExpiredBans)

		// Invitations
		webRoute.POST("/invitations", invitationHandl.CreateInvitation)
		webRoute.GET("/invitations", invitationHandl.GetMyInvitations)
		webRoute.DELETE("/invitations/:id", invitationHandl.RevokeInvitation)
		webRoute.GET("/invitations/stats", invitationHandl.GetInvitationStats)

		// Statistics & Dashboard
		webRoute.GET("/stats/kpi", statsHandl.GetKPIStats)
		webRoute.GET("/stats/activity", statsHandl.GetRecentActivity)
		webRoute.GET("/stats/trends", statsHandl.GetPlayerTrends)
		webRoute.GET("/stats/roles", statsHandl.GetRoleDistribution)
		webRoute.GET("/stats/bans", statsHandl.GetBanStats)
		webRoute.GET("/stats/overview", statsHandl.GetDashboardOverview)

		// Audit logs
		webRoute.GET("/audit-logs", auditHandl.ListAuditLogs)
		webRoute.GET("/audit-logs/stats", auditHandl.GetAuditStats)
		webRoute.GET("/audit-logs/export", auditHandl.ExportAuditLogs)
		webRoute.GET("/audit-logs/resource/:id", auditHandl.GetResourceHistory)

		// User info
		webRoute.GET("/me", authHandl.GetCurrentUser)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error %s", err.Error())
	}
}

// secretKeyMiddleware validates the map secret for VRChat endpoints
func secretKeyMiddleware(expectedSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.Query("secret")
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "missing secret parameter",
			})
			return
		}
		if secret != expectedSecret {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "invalid secret",
			})
			return
		}
		c.Next()
	}
}
