package models

import "time"

type Role string

const (
	ROLE_USER  Role = "user"
	ROLE_VIP   Role = "vip"
	ROLE_MOD   Role = "mod"
	ROLE_OWNER Role = "owner"
)

type Player struct {
	VRChatID    string     `gorm:"primaryKey;column:vrchat_id"`
	Role        []Role
	IsBanned    bool       `gorm:"default:false;index"`
	BanReason   string     `gorm:"type:text"`
	BanExpiresAt *time.Time `gorm:"index"` // NULL = permanent ban, date = temporary ban
	BannedAt    *time.Time // When the ban was applied
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}

// User model for web admin panel
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Role      Role      `gorm:"default:'user'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// RefreshToken model for JWT refresh tokens
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID"`
}

// InvitationCode model for user invitations
type InvitationCode struct {
	ID        uint      `gorm:"primaryKey"`
	Code      string    `gorm:"uniqueIndex;not null;size:32"`
	CreatedBy uint      `gorm:"not null;index"`
	Role      string    `gorm:"default:'user'"` // Default role for invited users
	MaxUses   int       `gorm:"default:1"`      // Number of times code can be used (-1 = unlimited)
	Uses      int       `gorm:"default:0"`      // Number of times code has been used
	ExpiresAt *time.Time `gorm:"index"`          // Expiration date (NULL = no expiry)
	IsActive  bool      `gorm:"default:true;index"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	Creator   User      `gorm:"foreignKey:CreatedBy"`
}

// AuditLog model for tracking all changes
type AuditLog struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index"`
	Action       string    `gorm:"index"` // create, update, delete, login, logout
	ResourceType string    `gorm:"index"` // player, user, invitation, etc
	ResourceID   string    // The ID of the resource that was changed
	OldValue     string    `gorm:"type:text"` // JSON string of old values
	NewValue     string    `gorm:"type:text"` // JSON string of new values
	IPAddress    string
	UserAgent    string    `gorm:"type:text"`
	Description  string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index"`
	User         User      `gorm:"foreignKey:UserID"`
}
