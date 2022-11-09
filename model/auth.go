package model

type LoginModel struct{
	Email string `json:"email" binding:"required,email"`
	Password string
}