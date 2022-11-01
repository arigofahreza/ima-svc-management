package model

import "time"

type EnumRole string

const (
	SUPERADMIN EnumRole = "superadmin"
	USER       EnumRole = "user"
	GUEST      EnumRole = "guest"
)

type RoleModel struct {
	Id          string     `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string     `json:"name" bson:"name"`
	Role        EnumRole   `json:"role" bson:"role"`
	Description string     `json:"description" bson:"description"`
	CreatedAt   time.Time  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type PaginateRoleModel struct {
	Order   string `json:"order,omitempty" bson:"order,omitempty"`
	OrderBy string `json:"orderBy,omitempty" bson:"orderBy,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Size    int    `json:"size,omitempty" bson:"size,omitempty"`
}
