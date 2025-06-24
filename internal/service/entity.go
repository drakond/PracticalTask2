package service

// TaskRequest - структура для тела запроса
type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}
