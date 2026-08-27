package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	serv *service.StatsService
}

func NewStatsHandler(serv *service.StatsService) *StatsHandler {
	return &StatsHandler{
		serv: serv,
	}
}

// GetKPIStats returns key performance indicators
func (h *StatsHandler) GetKPIStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	stats, err := h.serv.GetKPIStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve KPI stats",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRecentActivity returns recent activities
func (h *StatsHandler) GetRecentActivity(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 50 {
		limit = 50
	}

	activities, err := h.serv.GetRecentActivity(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve recent activity",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// GetPlayerTrends returns player count trends
func (h *StatsHandler) GetPlayerTrends(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days > 365 {
		days = 365
	}
	if days < 7 {
		days = 7
	}

	trends, err := h.serv.GetPlayerTrends(ctx, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve trends",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// GetRoleDistribution returns role distribution
func (h *StatsHandler) GetRoleDistribution(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	distribution, err := h.serv.GetRoleDistribution(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve role distribution",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, distribution)
}

// GetBanStats returns ban statistics
func (h *StatsHandler) GetBanStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	stats, err := h.serv.GetBanStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve ban stats",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetDashboardOverview returns comprehensive dashboard data
func (h *StatsHandler) GetDashboardOverview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 15*time.Second)
	defer cancel()

	kpiStats, _ := h.serv.GetKPIStats(ctx)
	activities, _ := h.serv.GetRecentActivity(ctx, 10)
	roleDistribution, _ := h.serv.GetRoleDistribution(ctx)
	banStats, _ := h.serv.GetBanStats(ctx)
	trends, _ := h.serv.GetPlayerTrends(ctx, 30)

	overview := gin.H{
		"kpi":                kpiStats,
		"recent_activities":  activities,
		"role_distribution":  roleDistribution,
		"ban_stats":          banStats,
		"player_trends":      trends,
		"timestamp":          time.Now(),
	}

	c.JSON(http.StatusOK, overview)
}
