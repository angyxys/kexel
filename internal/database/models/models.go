package models

type Role string

const (
	ROLE_USER  Role = "user"
	ROLE_VIP   Role = "vip"
	ROLE_MOD   Role = "mod"
	ROLE_OWNER Role = "owner"
)

type Player struct {
	VRChatID string `gorm:"primaryKey;column:vrchat_id"`
	Role     []Role
	IsBanned bool `gorm:"default:false"`
}
