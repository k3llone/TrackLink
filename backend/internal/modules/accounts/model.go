package accounts

import "time"

const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

type User struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email        string    `gorm:"type:text;not null;uniqueIndex:idx_users_email" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null" json:"-"`
	Role         string    `gorm:"type:text;not null;default:'customer'" json:"role"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()" json:"createdAt"`
}

func (User) TableName() string {
	return "users"
}
