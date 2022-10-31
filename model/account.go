package model

type AccountModel struct {
	Id        string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
	Password  string `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type PaginateAccountModel struct {
	Order   string `json:"order,omitempty" bson:"order,omitempty"`
	OrderBy string `json:"orderBy,omitempty" bson:"orderBy,omitempty"`
	Page    int    `json:"page,omitempty" bson:"page,omitempty"`
	Size    int    `json:"size,omitempty" bson:"size,omitempty"`
}

type RangePage struct {
	End   int `json:"end,omitempty" bson:"end,omitempty"`
	Start int `json:"start,omitempty" bson:"start,omitempty"`
}
