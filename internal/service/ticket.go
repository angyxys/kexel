package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
)

type TicketService struct {
	ticketRepo  *repository.TicketRepository
	commentRepo *repository.TicketCommentRepository
	auditServ   *AuditService
}

func NewTicketService(
	ticketRepo *repository.TicketRepository,
	commentRepo *repository.TicketCommentRepository,
	auditServ *AuditService,
) *TicketService {
	return &TicketService{
		ticketRepo:  ticketRepo,
		commentRepo: commentRepo,
		auditServ:   auditServ,
	}
}

type TicketInfo struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	Priority     string    `json:"priority"`
	Status       string    `json:"status"`
	AssignedTo   *uint     `json:"assigned_to"`
	AssignedName *string   `json:"assigned_name"`
	Resolution   string    `json:"resolution"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TicketCommentInfo struct {
	ID        uint      `json:"id"`
	TicketID  uint      `json:"ticket_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	IsInternal bool     `json:"is_internal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TicketStats struct {
	Total     int64 `json:"total"`
	Open      int64 `json:"open"`
	InProgress int64 `json:"in_progress"`
	Resolved  int64 `json:"resolved"`
	Closed    int64 `json:"closed"`
}

// CreateTicket creates a new ticket
func (s *TicketService) CreateTicket(ctx context.Context, userID uint, title, description, category, priority string) (*TicketInfo, error) {
	if title == "" || description == "" {
		return nil, errors.New("title and description are required")
	}

	if category == "" {
		category = "other"
	}
	if priority == "" {
		priority = "medium"
	}

	ticket := &models.Ticket{
		UserID:      userID,
		Title:       title,
		Description: description,
		Category:    category,
		Priority:    priority,
		Status:      "open",
	}

	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	// Log action
	auditLog := &models.AuditLog{
		UserID:       userID,
		Action:       "ticket.created",
		ResourceType: "ticket",
		ResourceID:   fmt.Sprintf("%d", ticket.ID),
	}
	s.auditServ.LogAction(ctx, auditLog)

	return s.getTicketInfo(ticket, &ticket.User, 0), nil
}

// GetTicket retrieves a ticket by ID
func (s *TicketService) GetTicket(ctx context.Context, ticketID uint, userID uint) (*TicketInfo, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	// Check access: ticket owner or admin/mod
	if ticket.UserID != userID {
		// TODO: Check if user is admin/mod
	}

	comments, _ := s.commentRepo.ListTicketComments(ctx, ticketID, false)
	commentCount := len(comments)

	return s.getTicketInfo(ticket, &ticket.User, commentCount), nil
}

// ListUserTickets retrieves all tickets for a user
func (s *TicketService) ListUserTickets(ctx context.Context, userID uint) ([]TicketInfo, error) {
	tickets, err := s.ticketRepo.ListUserTickets(ctx, userID)
	if err != nil {
		return nil, err
	}

	var infos []TicketInfo
	for _, ticket := range tickets {
		infos = append(infos, *s.getTicketInfo(&ticket, &ticket.User, 0))
	}
	return infos, nil
}

// ListAllTickets retrieves all tickets (admin only)
func (s *TicketService) ListAllTickets(ctx context.Context, page, pageSize int) ([]TicketInfo, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	tickets, err := s.ticketRepo.ListAllTickets(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	var infos []TicketInfo
	for _, ticket := range tickets {
		var assignedName *string
		if ticket.AssignedUser != nil {
			assignedName = &ticket.AssignedUser.Username
		}

		comments, _ := s.commentRepo.ListTicketComments(ctx, ticket.ID, true)
		infos = append(infos, TicketInfo{
			ID:           ticket.ID,
			UserID:       ticket.UserID,
			Username:     ticket.User.Username,
			Title:        ticket.Title,
			Description:  ticket.Description,
			Category:     ticket.Category,
			Priority:     ticket.Priority,
			Status:       ticket.Status,
			AssignedTo:   ticket.AssignedTo,
			AssignedName: assignedName,
			Resolution:   ticket.Resolution,
			ResolvedAt:   ticket.ResolvedAt,
			CommentCount: len(comments),
			CreatedAt:    ticket.CreatedAt,
			UpdatedAt:    ticket.UpdatedAt,
		})
	}

	return infos, nil
}

// FilterTickets filters tickets
func (s *TicketService) FilterTickets(ctx context.Context, status, priority, category string, page, pageSize int) ([]TicketInfo, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	tickets, err := s.ticketRepo.FilterTickets(ctx, status, priority, category, pageSize, offset)
	if err != nil {
		return nil, err
	}

	var infos []TicketInfo
	for _, ticket := range tickets {
		var assignedName *string
		if ticket.AssignedUser != nil {
			assignedName = &ticket.AssignedUser.Username
		}

		comments, _ := s.commentRepo.ListTicketComments(ctx, ticket.ID, true)
		infos = append(infos, TicketInfo{
			ID:           ticket.ID,
			UserID:       ticket.UserID,
			Username:     ticket.User.Username,
			Title:        ticket.Title,
			Description:  ticket.Description,
			Category:     ticket.Category,
			Priority:     ticket.Priority,
			Status:       ticket.Status,
			AssignedTo:   ticket.AssignedTo,
			AssignedName: assignedName,
			Resolution:   ticket.Resolution,
			ResolvedAt:   ticket.ResolvedAt,
			CommentCount: len(comments),
			CreatedAt:    ticket.CreatedAt,
			UpdatedAt:    ticket.UpdatedAt,
		})
	}

	return infos, nil
}

// UpdateTicket updates a ticket
func (s *TicketService) UpdateTicket(ctx context.Context, ticketID uint, userID uint, status, priority, assignedTo, resolution string) error {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := ticket.Status
	oldPriority := ticket.Priority

	if status != "" {
		ticket.Status = status
		if status == "resolved" && ticket.ResolvedAt == nil {
			now := time.Now()
			ticket.ResolvedAt = &now
		}
	}

	if priority != "" {
		ticket.Priority = priority
	}

	if assignedTo != "" {
		// Parse assignedTo as uint
		var assignedToID uint
		fmt.Sscanf(assignedTo, "%d", &assignedToID)
		if assignedToID > 0 {
			ticket.AssignedTo = &assignedToID
		}
	}

	if resolution != "" {
		ticket.Resolution = resolution
	}

	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return err
	}

	// Log changes
	if oldStatus != ticket.Status || oldPriority != ticket.Priority {
		auditLog := &models.AuditLog{
			UserID:       userID,
			Action:       "ticket.updated",
			ResourceType: "ticket",
			ResourceID:   fmt.Sprintf("%d", ticketID),
			OldValue:     oldStatus,
			NewValue:     ticket.Status,
		}
		s.auditServ.LogAction(ctx, auditLog)
	}

	return nil
}

// AddComment adds a comment to a ticket
func (s *TicketService) AddComment(ctx context.Context, ticketID uint, userID uint, content string, isInternal bool) (*TicketCommentInfo, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}

	comment := &models.TicketComment{
		TicketID:   ticketID,
		UserID:     userID,
		Content:    content,
		IsInternal: isInternal,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Log action
	auditLog := &models.AuditLog{
		UserID:       userID,
		Action:       "ticket.comment",
		ResourceType: "ticket_comment",
		ResourceID:   fmt.Sprintf("%d", comment.ID),
	}
	s.auditServ.LogAction(ctx, auditLog)

	return &TicketCommentInfo{
		ID:        comment.ID,
		TicketID:  comment.TicketID,
		UserID:    comment.UserID,
		Username:  comment.User.Username,
		Content:   comment.Content,
		IsInternal: comment.IsInternal,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}, nil
}

// GetTicketComments retrieves comments for a ticket
func (s *TicketService) GetTicketComments(ctx context.Context, ticketID uint, includeInternal bool) ([]TicketCommentInfo, error) {
	comments, err := s.commentRepo.ListTicketComments(ctx, ticketID, includeInternal)
	if err != nil {
		return nil, err
	}

	var infos []TicketCommentInfo
	for _, comment := range comments {
		infos = append(infos, TicketCommentInfo{
			ID:         comment.ID,
			TicketID:   comment.TicketID,
			UserID:     comment.UserID,
			Username:   comment.User.Username,
			Content:    comment.Content,
			IsInternal: comment.IsInternal,
			CreatedAt:  comment.CreatedAt,
			UpdatedAt:  comment.UpdatedAt,
		})
	}
	return infos, nil
}

// GetTicketStats retrieves ticket statistics
func (s *TicketService) GetTicketStats(ctx context.Context) (*TicketStats, error) {
	total, _ := s.ticketRepo.GetTicketCount(ctx)
	open, _ := s.ticketRepo.GetTicketCountByStatus(ctx, "open")
	inProgress, _ := s.ticketRepo.GetTicketCountByStatus(ctx, "in-progress")
	resolved, _ := s.ticketRepo.GetTicketCountByStatus(ctx, "resolved")
	closed, _ := s.ticketRepo.GetTicketCountByStatus(ctx, "closed")

	return &TicketStats{
		Total:      total,
		Open:       open,
		InProgress: inProgress,
		Resolved:   resolved,
		Closed:     closed,
	}, nil
}

// DeleteTicket deletes a ticket (admin only)
func (s *TicketService) DeleteTicket(ctx context.Context, ticketID uint, userID uint) error {
	// Delete comments first
	if err := s.commentRepo.DeleteTicketComments(ctx, ticketID); err != nil {
		return err
	}

	// Delete ticket
	if err := s.ticketRepo.Delete(ctx, ticketID); err != nil {
		return err
	}

	// Log action
	auditLog := &models.AuditLog{
		UserID:       userID,
		Action:       "ticket.deleted",
		ResourceType: "ticket",
		ResourceID:   fmt.Sprintf("%d", ticketID),
	}
	s.auditServ.LogAction(ctx, auditLog)

	return nil
}

func (s *TicketService) getTicketInfo(ticket *models.Ticket, user *models.User, commentCount int) *TicketInfo {
	var assignedName *string
	if ticket.AssignedUser != nil {
		assignedName = &ticket.AssignedUser.Username
	}

	return &TicketInfo{
		ID:           ticket.ID,
		UserID:       ticket.UserID,
		Username:     user.Username,
		Title:        ticket.Title,
		Description:  ticket.Description,
		Category:     ticket.Category,
		Priority:     ticket.Priority,
		Status:       ticket.Status,
		AssignedTo:   ticket.AssignedTo,
		AssignedName: assignedName,
		Resolution:   ticket.Resolution,
		ResolvedAt:   ticket.ResolvedAt,
		CommentCount: commentCount,
		CreatedAt:    ticket.CreatedAt,
		UpdatedAt:    ticket.UpdatedAt,
	}
}
