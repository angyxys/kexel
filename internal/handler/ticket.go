package handler

import (
	"net/http"
	"strconv"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketServ *service.TicketService
}

func NewTicketHandler(ticketServ *service.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketServ: ticketServ,
	}
}

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
}

type UpdateTicketRequest struct {
	Status     string `json:"status"`
	Priority   string `json:"priority"`
	AssignedTo string `json:"assigned_to"`
	Resolution string `json:"resolution"`
}

type AddCommentRequest struct {
	Content    string `json:"content" binding:"required"`
	IsInternal bool   `json:"is_internal"`
}

// CreateTicket creates a new ticket
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	ticket, err := h.ticketServ.CreateTicket(c, userID.(uint), req.Title, req.Description, req.Category, req.Priority)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// GetTicket retrieves a ticket
func (h *TicketHandler) GetTicket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ticket id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	ticket, err := h.ticketServ.GetTicket(c, uint(ticketID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "ticket not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// ListUserTickets lists user's tickets
func (h *TicketHandler) ListUserTickets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	tickets, err := h.ticketServ.ListUserTickets(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch tickets",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   tickets,
		"status": http.StatusOK,
	})
}

// ListAllTickets lists all tickets (admin)
func (h *TicketHandler) ListAllTickets(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			page = val
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 && val <= 100 {
			pageSize = val
		}
	}

	tickets, err := h.ticketServ.ListAllTickets(c, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch tickets",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   tickets,
		"status": http.StatusOK,
	})
}

// FilterTickets filters tickets
func (h *TicketHandler) FilterTickets(c *gin.Context) {
	status := c.Query("status")
	priority := c.Query("priority")
	category := c.Query("category")
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil {
			page = val
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 && val <= 100 {
			pageSize = val
		}
	}

	tickets, err := h.ticketServ.FilterTickets(c, status, priority, category, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch tickets",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   tickets,
		"status": http.StatusOK,
	})
}

// UpdateTicket updates a ticket
func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ticket id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.ticketServ.UpdateTicket(c, uint(ticketID), userID.(uint), req.Status, req.Priority, req.AssignedTo, req.Resolution); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ticket updated",
		"status":  http.StatusOK,
	})
}

// AddComment adds a comment to a ticket
func (h *TicketHandler) AddComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ticket id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	comment, err := h.ticketServ.AddComment(c, uint(ticketID), userID.(uint), req.Content, req.IsInternal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetTicketComments gets comments for a ticket
func (h *TicketHandler) GetTicketComments(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ticket id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	comments, err := h.ticketServ.GetTicketComments(c, uint(ticketID), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch comments",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   comments,
		"status": http.StatusOK,
	})
}

// GetTicketStats gets ticket statistics
func (h *TicketHandler) GetTicketStats(c *gin.Context) {
	stats, err := h.ticketServ.GetTicketStats(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch stats",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DeleteTicket deletes a ticket (admin only)
func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid ticket id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.ticketServ.DeleteTicket(c, uint(ticketID), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ticket deleted",
		"status":  http.StatusOK,
	})
}
