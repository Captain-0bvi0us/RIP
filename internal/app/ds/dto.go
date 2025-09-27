package ds

import "time"

// DTO для API ответов
type FactorDTO struct {
	ID       uint     `json:"id"`
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	Image    *string  `json:"image"`
	Argument *float64 `json:"argument"`
	Status   *bool    `json:"status"`
}

type FactorCreateRequest struct {
	Title    string   `json:"title" binding:"required"`
	Text     string   `json:"text" binding:"required"`
	Argument *float64 `json:"argument"`
}

type FactorUpdateRequest struct {
	Title    *string  `json:"title"`
	Text     *string  `json:"text"`
	Argument *float64 `json:"argument"`
}

type FraxDTO struct {
	ID             uint              `json:"id"`
	Status         int               `json:"status"`
	CreationDate   time.Time         `json:"creation_date"`
	CreatorLogin   string            `json:"creator_login"`
	ModeratorLogin *string           `json:"moderator_login"`
	FormingDate    *time.Time        `json:"forming_date"`
	ComplitionDate *time.Time        `json:"complition_date"`
	Age            *int              `json:"age"`
	Gender         *bool             `json:"gender"`
	Weight         *int              `json:"weight"`
	Height         *int              `json:"height"`
	POF            *float64          `json:"POF"`
	PHF            *float64          `json:"PHF"`
	Factors        []FactorInFraxDTO `json:"factors"`
}

type FactorInFraxDTO struct {
	ID          uint     `json:"id"`
	FactorID    uint     `json:"factor_id"`
	Title       string   `json:"title"`
	Text        string   `json:"text"`
	Image       *string  `json:"image"`
	Argument    *float64 `json:"argument"`
	Description *string  `json:"description"`
}

type FraxUpdateRequest struct {
	Age    *int  `json:"age"`
	Gender *bool `json:"gender"`
	Weight *int  `json:"weight"`
	Height *int  `json:"height"`
}

type FraxResolveRequest struct {
	Action string `json:"action" binding:"required"` // "complete" | "reject"
}

type FactorToFraxUpdateRequest struct {
	Description *string `json:"description"`
}

type CartBadgeDTO struct {
	FraxID *uint `json:"frax_id"`
	Count  int   `json:"count"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDTO struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Moderator bool   `json:"moderator"`
}

type UserUpdateRequest struct {
	Username  *string `json:"username"`
	Password  *string `json:"password"`
	Moderator *bool   `json:"moderator"`
}

type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

// Отображение списка факторов с их атрибутами и количеством
type PaginatedResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}
