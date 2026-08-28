package repository

import (
	"context"

	"github.com/angyxys/kexel/internal/database/models"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// Create creates a new ticket
func (r *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

// GetByID retrieves a ticket by ID
func (r *TicketRepository) GetByID(ctx context.Context, id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.WithContext(ctx).Preload("User").Preload("AssignedUser").First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// ListUserTickets retrieves all tickets for a user
func (r *TicketRepository) ListUserTickets(ctx context.Context, userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("AssignedUser").
		Order("created_at DESC").
		Find(&tickets).Error
	return tickets, err
}

// ListAllTickets retrieves all tickets (for admins)
func (r *TicketRepository) ListAllTickets(ctx context.Context, limit int, offset int) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("AssignedUser").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tickets).Error
	return tickets, err
}

// FilterTickets filters tickets by status, priority, category
func (r *TicketRepository) FilterTickets(ctx context.Context, status, priority, category string, limit, offset int) ([]models.Ticket, error) {
	query := r.db.WithContext(ctx)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var tickets []models.Ticket
	err := query.
		Preload("User").
		Preload("AssignedUser").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tickets).Error

	return tickets, err
}

// Update updates a ticket
func (r *TicketRepository) Update(ctx context.Context, ticket *models.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

// Delete deletes a ticket
func (r *TicketRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Ticket{}, id).Error
}

// GetTicketCount returns total count of tickets
func (r *TicketRepository) GetTicketCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Ticket{}).Count(&count).Error
	return count, err
}

// GetTicketCountByStatus returns count of tickets by status
func (r *TicketRepository) GetTicketCountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Ticket{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

type TicketCommentRepository struct {
	db *gorm.DB
}

func NewTicketCommentRepository(db *gorm.DB) *TicketCommentRepository {
	return &TicketCommentRepository{db: db}
}

// Create creates a new ticket comment
func (r *TicketCommentRepository) Create(ctx context.Context, comment *models.TicketComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// GetByID retrieves a comment by ID
func (r *TicketCommentRepository) GetByID(ctx context.Context, id uint) (*models.TicketComment, error) {
	var comment models.TicketComment
	err := r.db.WithContext(ctx).Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListTicketComments retrieves all comments for a ticket
func (r *TicketCommentRepository) ListTicketComments(ctx context.Context, ticketID uint, includeInternal bool) ([]models.TicketComment, error) {
	var comments []models.TicketComment
	query := r.db.WithContext(ctx).
		Where("ticket_id = ?", ticketID).
		Preload("User").
		Order("created_at ASC")

	if !includeInternal {
		query = query.Where("is_internal = false")
	}

	err := query.Find(&comments).Error
	return comments, err
}

// Update updates a comment
func (r *TicketCommentRepository) Update(ctx context.Context, comment *models.TicketComment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// Delete deletes a comment
func (r *TicketCommentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TicketComment{}, id).Error
}

// DeleteTicketComments deletes all comments for a ticket
func (r *TicketCommentRepository) DeleteTicketComments(ctx context.Context, ticketID uint) error {
	return r.db.WithContext(ctx).Where("ticket_id = ?", ticketID).Delete(&models.TicketComment{}).Error
}
