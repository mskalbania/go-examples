package model

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type PostUser struct {
	Email string `json:"email" binding:"required"`
}
