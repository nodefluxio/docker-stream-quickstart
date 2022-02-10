package entity

import (
	"time"

	"github.com/lib/pq"
)

type UserRole string

const (
	UserRoleSuperAdmin UserRole = "superadmin"
	UserRoleOperator   UserRole = "operator"
)

type User struct {
	ID        uint64        `json:"id" db:"id"`
	Email     string        `json:"email" db:"email"`
	Username  string        `json:"username" db:"username"`
	Password  string        `json:"-" db:"password"`
	Fullname  string        `json:"fullname" db:"fullname"`
	Avatar    byte          `json:"avatar" db:"avatar"`
	Role      string        `json:"role" db:"role"`
	SiteID    pq.Int64Array `json:"site_id" db:"site_id"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at" db:"deleted_at"`
}

// TableName overrides the table name used by User to `user_access`
func (User) TableName() string {
	return "user_access"
}
