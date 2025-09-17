package entities

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
	Role     string `json:"role"` // admin or user
}
