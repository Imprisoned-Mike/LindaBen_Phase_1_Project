package Login

type Update struct {
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required"`
	RoleID uint   `gorm:"not null" json:"role_id"`
}
