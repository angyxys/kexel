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

// Session model for tracking user sessions
type Session struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	SessionToken string    `gorm:"uniqueIndex;not null;size:64"`
	IPAddress    string    `gorm:"index"`
	UserAgent    string    `gorm:"type:text"`
	DeviceName   string    // e.g., "Chrome on Windows", "Safari on iPhone"
	IsActive     bool      `gorm:"default:true;index"`
	LastActivity time.Time `gorm:"index"` // Last action timestamp
	LoginAt      time.Time `gorm:"autoCreateTime;index"`
	LogoutAt     *time.Time // NULL if session is active
	ExpiresAt    time.Time `gorm:"index"` // Token expiration
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	User         User      `gorm:"foreignKey:UserID"`
}

// TOTPSecret model for storing 2FA TOTP secrets
type TOTPSecret struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null;uniqueIndex"`
	Secret        string    `gorm:"not null"` // Base32-encoded TOTP secret
	BackupCodes   string    `gorm:"type:text"` // JSON array of backup codes
	IsEnabled     bool      `gorm:"default:false;index"`
	EnabledAt     *time.Time
	LastUsedAt    *time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	User          User      `gorm:"foreignKey:UserID"`
}

// APIKey model for API authentication
type APIKey struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null;index"`
	Name          string    `gorm:"not null"` // User-friendly name
	Key           string    `gorm:"uniqueIndex;not null;size:64"` // Hashed API key
	KeyPrefix     string    `gorm:"index;size:8"` // First 8 chars for preview
	Scopes        string    `gorm:"type:text"` // JSON array of scopes
	IPWhitelist   string    `gorm:"type:text"` // JSON array of IPs (optional)
	IsActive      bool      `gorm:"default:true;index"`
	LastUsedAt    *time.Time
	ExpiresAt     *time.Time `gorm:"index"` // NULL = no expiry
	RateLimit     int       `gorm:"default:1000"` // Requests per hour
	CreatedAt     time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	User          User      `gorm:"foreignKey:UserID"`
}

// Webhook model for event subscriptions
type Webhook struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null;index"`
	Name          string    `gorm:"not null"`
	URL           string    `gorm:"not null;type:text"` // Webhook endpoint
	Events        string    `gorm:"type:text"` // JSON array of event types
	Secret        string    `gorm:"not null"` // For HMAC verification
	Headers       string    `gorm:"type:text"` // JSON custom headers
	IsActive      bool      `gorm:"default:true;index"`
	FailureCount  int       `gorm:"default:0"`
	LastTriedAt   *time.Time
	LastSuccessAt *time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	User          User      `gorm:"foreignKey:UserID"`
}

// WebhookEvent model for event delivery history
type WebhookEvent struct {
	ID          uint      `gorm:"primaryKey"`
	WebhookID   uint      `gorm:"not null;index"`
	EventType   string    `gorm:"index"` // player.created, player.updated, ban.created, etc
	Payload     string    `gorm:"type:text"` // JSON payload
	StatusCode  int       // HTTP response code
	Response    string    `gorm:"type:text"` // Response body
	Attempts    int       `gorm:"default:1"`
	NextRetry   *time.Time
	IsDelivered bool      `gorm:"default:false;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Webhook     Webhook   `gorm:"foreignKey:WebhookID"`
}

// DiscordIntegration model for Discord bot settings
type DiscordIntegration struct {
	ID               uint      `gorm:"primaryKey"`
	UserID           uint      `gorm:"not null;uniqueIndex"`
	BotToken         string    `gorm:"not null"` // Discord bot token
	GuildID          string    `gorm:"not null"` // Discord server ID
	ModLogChannelID  string    // Channel for moderation logs
	AnnouncementChannelID string // Channel for announcements
	RoleMapping      string    `gorm:"type:text"` // JSON: {kexel_role: discord_role_id}
	IsConnected      bool      `gorm:"default:false;index"`
	LastSyncAt       *time.Time
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
	User             User      `gorm:"foreignKey:UserID"`
}

// PatreonIntegration model for Patreon sync
type PatreonIntegration struct {
	ID                uint      `gorm:"primaryKey"`
	UserID            uint      `gorm:"not null;uniqueIndex"`
	CampaignID        string    `gorm:"not null"` // Patreon campaign ID
	AccessToken       string    `gorm:"not null"` // OAuth token
	RefreshToken      string    // OAuth refresh token
	PatreonWebhookID  string    // Webhook ID for events
	SyncInterval      int       `gorm:"default:3600"` // Sync every 1 hour
	RoleMapping       string    `gorm:"type:text"` // JSON: {tier_id: kexel_role}
	IsEnabled         bool      `gorm:"default:false;index"`
	LastSyncAt        *time.Time
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
	User              User      `gorm:"foreignKey:UserID"`
}

// RateLimitRule model for IP-based rate limiting
type RateLimitRule struct {
	ID           uint      `gorm:"primaryKey"`
	IPAddress    string    `gorm:"uniqueIndex;not null"`
	Endpoint     string    `gorm:"index"` // Empty = all endpoints
	RequestLimit int       `gorm:"default:100"` // Requests per window
	WindowSize   int       `gorm:"default:3600"` // Time window in seconds
	IsActive     bool      `gorm:"default:true;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// RateLimitLog model for tracking rate limit hits
type RateLimitLog struct {
	ID          uint      `gorm:"primaryKey"`
	RuleID      uint      `gorm:"index"`
	IPAddress   string    `gorm:"index"`
	Endpoint    string
	RequestCount int
	BlockedAt   time.Time `gorm:"autoCreateTime;index"`
}

// Banner model for configurable banners and posters
type Banner struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	Name         string    `gorm:"not null"` // Banner name
	Type         string    `gorm:"index"` // banner, poster, avatar, etc
	Title        string    // Display title
	Description  string    `gorm:"type:text"`
	ImageURL     string    `gorm:"not null"` // S3/MinIO URL
	S3Key        string    `gorm:"not null"` // MinIO object key
	Width        int       // Image width
	Height       int       // Image height
	FileSize     int64     // File size in bytes
	MimeType     string
	IsActive     bool      `gorm:"default:true;index"`
	DisplayOrder int       `gorm:"default:0"` // Order for display
	CreatedAt    time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	User         User      `gorm:"foreignKey:UserID"`
}

// Ticket model for bug reports and feature requests
type Ticket struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	Title        string    `gorm:"not null;index"` // Ticket title
	Description  string    `gorm:"type:text"` // Detailed description
	Category     string    `gorm:"index"` // bug, feature-request, other
	Priority     string    `gorm:"index;default:'medium'"` // low, medium, high, critical
	Status       string    `gorm:"index;default:'open'"` // open, in-progress, resolved, closed
	AssignedTo   *uint     `gorm:"index"` // Admin user ID
	Resolution   string    `gorm:"type:text"` // Resolution details
	ResolvedAt   *time.Time
	CreatedAt    time.Time `gorm:"autoCreateTime;index"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	User         User      `gorm:"foreignKey:UserID"`
	AssignedUser *User     `gorm:"foreignKey:AssignedTo"`
}

// TicketComment model for ticket discussions
type TicketComment struct {
	ID        uint      `gorm:"primaryKey"`
	TicketID  uint      `gorm:"not null;index"`
	UserID    uint      `gorm:"not null;index"`
	Content   string    `gorm:"type:text"` // Comment content
	IsInternal bool     `gorm:"default:false"` // Only visible to admins
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Ticket    Ticket    `gorm:"foreignKey:TicketID"`
	User      User      `gorm:"foreignKey:UserID"`
}
