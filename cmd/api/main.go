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
		&models.Session{},
		&models.TOTPSecret{},
		&models.APIKey{},
		&models.Webhook{},
		&models.WebhookEvent{},
		&models.DiscordIntegration{},
		&models.PatreonIntegration{},
		&models.RateLimitRule{},
		&models.RateLimitLog{},
		&models.Banner{},
		&models.Ticket{},
		&models.TicketComment{},
	); err != nil {
		log.Fatalf("error migrating database: %v", err)
	}

	// Initialize repositories
	playerRepo := repository.NewPlayerRepository(db)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	totpRepo := repository.NewTOTPRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)
	webhookEventRepo := repository.NewWebhookEventRepository(db)
	discordRepo := repository.NewDiscordIntegrationRepository(db)
	patreonRepo := repository.NewPatreonIntegrationRepository(db)
	rateLimitRuleRepo := repository.NewRateLimitRuleRepository(db)
	rateLimitLogRepo := repository.NewRateLimitLogRepository(db)
	bannerRepo := repository.NewBannerRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	ticketCommentRepo := repository.NewTicketCommentRepository(db)

	// Initialize MinIO client
	minioClient, err := service.NewMinIOClient()
	if err != nil {
		log.Printf("Warning: MinIO initialization failed: %v. Image storage will not work.", err)
		minioClient = nil
	}

	// Initialize services
	playerServ := service.NewPlayerService(playerRepo)
	authServ := service.NewAuthServiceWithInvitations(userRepo, refreshTokenRepo, invitationRepo, cfg.JWT_SECRET)
	auditServ := service.NewAuditService(auditRepo)
	banServ := service.NewBanService(playerRepo)
	statsServ := service.NewStatsService(playerRepo, auditRepo, userRepo)
	invitationServ := service.NewInvitationService(invitationRepo)
	sessionServ := service.NewSessionService(sessionRepo, userRepo)
	totpServ := service.NewTOTPService(totpRepo)
	apiKeyServ := service.NewAPIKeyService(apiKeyRepo)
	webhookServ := service.NewWebhookService(webhookRepo, webhookEventRepo)
	discordServ := service.NewDiscordService(discordRepo)
	patreonServ := service.NewPatreonService(patreonRepo, playerRepo, userRepo)
	rateLimitServ := service.NewRateLimitService(rateLimitRuleRepo, rateLimitLogRepo)
	storageServ := service.NewStorageService(bannerRepo, minioClient, "kexel-banners")
	ticketServ := service.NewTicketService(ticketRepo, ticketCommentRepo, auditServ)

	// Initialize handlers
	playerHandl := handler.NewPlayerHandler(playerServ)
	authHandl := handler.NewAuthHandlerWithSession(authServ, sessionServ)
	auditHandl := handler.NewAuditHandler(auditServ)
	banHandl := handler.NewBanHandler(banServ)
	statsHandl := handler.NewStatsHandler(statsServ)
	invitationHandl := handler.NewInvitationHandler(invitationServ)
	sessionHandl := handler.NewSessionHandler(sessionServ)
	totpHandl := handler.NewTOTPHandler(totpServ)
	apiKeyHandl := handler.NewAPIKeyHandler(apiKeyServ)
	webhookHandl := handler.NewWebhookHandler(webhookServ)
	integrationHandl := handler.NewIntegrationHandler(discordServ, patreonServ, rateLimitServ)
	storageHandl := handler.NewStorageHandler(storageServ)
	ticketHandl := handler.NewTicketHandler(ticketServ)

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

		// Sessions
		webRoute.GET("/sessions", sessionHandl.GetSessions)
		webRoute.GET("/sessions/stats", sessionHandl.GetSessionStats)
		webRoute.DELETE("/sessions/:id", sessionHandl.LogoutSession)
		webRoute.POST("/sessions/logout-all", sessionHandl.LogoutAllSessions)

		// 2FA / TOTP
		webRoute.POST("/2fa/setup", totpHandl.GetTOTPSetup)
		webRoute.POST("/2fa/verify", totpHandl.VerifyTOTP)
		webRoute.GET("/2fa/status", totpHandl.GetTOTPStatus)
		webRoute.DELETE("/2fa", totpHandl.DisableTOTP)
		webRoute.POST("/2fa/verify-code", totpHandl.VerifyTOTPCode)
		webRoute.POST("/2fa/backup-code", totpHandl.VerifyBackupCode)

		// API Keys
		webRoute.POST("/api-keys", apiKeyHandl.CreateAPIKey)
		webRoute.GET("/api-keys", apiKeyHandl.GetAPIKeys)
		webRoute.DELETE("/api-keys/:id", apiKeyHandl.DeleteAPIKey)
		webRoute.POST("/api-keys/:id/revoke", apiKeyHandl.RevokeAPIKey)
		webRoute.PATCH("/api-keys/:id/rate-limit", apiKeyHandl.UpdateRateLimit)
		webRoute.PATCH("/api-keys/:id/scopes", apiKeyHandl.UpdateScopes)
		webRoute.GET("/api-keys/scopes/available", apiKeyHandl.GetAvailableScopes)

		// Webhooks
		webRoute.POST("/webhooks", webhookHandl.CreateWebhook)
		webRoute.GET("/webhooks", webhookHandl.GetWebhooks)
		webRoute.DELETE("/webhooks/:id", webhookHandl.DeleteWebhook)
		webRoute.POST("/webhooks/:id/disable", webhookHandl.DisableWebhook)
		webRoute.GET("/webhooks/:id/events", webhookHandl.GetWebhookEvents)
		webRoute.GET("/webhooks/events/available", webhookHandl.GetAvailableEvents)

		// Banners & Storage
		webRoute.POST("/banners", storageHandl.UploadBanner)
		webRoute.GET("/banners", storageHandl.GetUserBanners)
		webRoute.GET("/banners/:type", storageHandl.GetBannersByType)
		webRoute.PATCH("/banners/:id", storageHandl.UpdateBanner)
		webRoute.DELETE("/banners/:id", storageHandl.DeleteBanner)

		// Tickets & Support
		webRoute.POST("/tickets", ticketHandl.CreateTicket)
		webRoute.GET("/tickets", ticketHandl.ListUserTickets)
		webRoute.GET("/tickets/all", ticketHandl.ListAllTickets)
		webRoute.GET("/tickets/filter", ticketHandl.FilterTickets)
		webRoute.GET("/tickets/stats", ticketHandl.GetTicketStats)
		webRoute.GET("/tickets/:id", ticketHandl.GetTicket)
		webRoute.PATCH("/tickets/:id", ticketHandl.UpdateTicket)
		webRoute.DELETE("/tickets/:id", ticketHandl.DeleteTicket)
		webRoute.POST("/tickets/:id/comments", ticketHandl.AddComment)
		webRoute.GET("/tickets/:id/comments", ticketHandl.GetTicketComments)

		// Integrations (Discord, Patreon, Rate Limiting)
		webRoute.POST("/integrations/discord", integrationHandl.SetupDiscord)
		webRoute.GET("/integrations/discord", integrationHandl.GetDiscordIntegration)
		webRoute.PATCH("/integrations/discord/channels", integrationHandl.ConfigureDiscordChannels)
		webRoute.POST("/integrations/discord/test", integrationHandl.TestDiscordConnection)
		webRoute.DELETE("/integrations/discord", integrationHandl.DisconnectDiscord)

		// Patreon Integration
		webRoute.GET("/integrations/patreon/oauth-url", integrationHandl.GetPatreonOAuthURL)
		webRoute.POST("/integrations/patreon/oauth-callback", integrationHandl.HandlePatreonOAuthCallback)
		webRoute.GET("/integrations/patreon", integrationHandl.GetPatreonIntegration)
		webRoute.PATCH("/integrations/patreon/tier-mapping", integrationHandl.ConfigurePatreonTierMapping)
		webRoute.POST("/integrations/patreon/sync", integrationHandl.SyncPatreonMembers)
		webRoute.DELETE("/integrations/patreon", integrationHandl.DisconnectPatreon)

		// Rate Limiting
		webRoute.GET("/admin/rate-limit/blocks", integrationHandl.GetRateLimitBlocks)

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
