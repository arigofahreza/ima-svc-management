package model

type LoginModel struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty" binding:"required"`
	Password string `json:"password,omitempty" bson:"password,omitempty" binding:"required"`
}
