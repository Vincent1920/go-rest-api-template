package dto

import "time"

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=3,max=200"`

	Description string `json:"description" binding:"omitempty,max=500"`

	Status string `json:"status" binding:"omitempty,oneof=Pending 'In Progress' Done"`
}

type UpdateTodoRequest struct {
	Title *string `json:"title" binding:"omitempty,min=3,max=200"`

	Description *string `json:"description" binding:"omitempty,max=500"`

	Status *string `json:"status" binding:"omitempty,oneof=Pending 'In Progress' Done"`
}

type TodoResponse struct {
	ID uint `json:"id"`

	Title string `json:"title"`

	Description string `json:"description"`

	Status string `json:"status"`

	UserID uint `json:"user_id"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

type TodoListResponse struct {
	Data []TodoResponse `json:"data"`

	Total int64 `json:"total"`

	Page int `json:"page"`

	Limit int `json:"limit"`
}
