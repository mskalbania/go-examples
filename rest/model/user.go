package model

type User struct {
	ID    string `json:"id" uri:"id"`
	Email string `json:"email"`
}
