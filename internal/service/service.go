package service

import (
	"context"
	"prtask2/internal/repo"
	"prtask2/pkg/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Слой бизнес-логики. Тут основная логика сервиса

// Service - интерфейс для бизнес-логики
type Service interface {
	CreateTask(c *fiber.Ctx) error
	GetAllTasks(c *fiber.Ctx) error
	GetTaskByID(c *fiber.Ctx) error
	UpdateTask(c *fiber.Ctx) error
	DeleteTask(c *fiber.Ctx) error
	GetTasksByUserID(c *fiber.Ctx) error
	GetTasksByUsername(c *fiber.Ctx) error
}

type service struct {
	repo *repo.Repository
}

// NewService - конструктор сервиса
func NewService(repo *repo.Repository) *service {
	return &service{repo: repo}
}

// Создать задачу
func (s *service) CreateTask(c *fiber.Ctx) error {
	var req repo.Task

	// Парсинг JSON
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid request body",
		})
	}

	// Валидация
	if err := validator.Validate(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	// Проверка существования пользователя
	if _, err := s.repo.GetUserByID(context.Background(), req.UserID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": "User not found"})
	}

	if err := s.repo.CreateTask(context.Background(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   req,
	})
}

// Получить все задачи
func (s *service) GetAllTasks(c *fiber.Ctx) error {
	tasks, err := s.repo.GetAllTasks(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   tasks,
	})
}

// Получить задачу по ID
func (s *service) GetTaskByID(c *fiber.Ctx) error {
	// Извлечение ID из URL
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid task ID",
		})
	}

	// Получение задачи из storage
	task, err := s.repo.GetTaskByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  "Task not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   task,
	})
}

// Обновить задачу
func (s *service) UpdateTask(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": "Invalid task ID"})
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": "Invalid request body"})
	}

	if req.Title == "" || req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": "Title and status are required"})
	}

	if err := s.repo.UpdateTask(context.Background(), id, req.Title, req.Description, req.Status); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "error": "Task not found"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Task updated"})
}

// Удалить задачу
func (s *service) DeleteTask(c *fiber.Ctx) error {
	// Извлечение ID из URL
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid task ID",
		})
	}

	// Удаление через storage
	if err := s.repo.DeleteTask(context.Background(), id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  "Task not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Task deleted successfully",
	})
}

// Получить задачи по userID
func (s *service) GetTasksByUserID(c *fiber.Ctx) error {
	// Извлечение userID из URL
	userIDStr := c.Params("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": "Invalid user ID"})
	}

	// Получение задач из storage
	tasks, err := s.repo.GetTasksByUserID(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": tasks})
}

// Получить задачи по username
func (s *service) GetTasksByUsername(c *fiber.Ctx) error {
	// Извлечение username из URL
	username := c.Params("username")

	// Получение задач из storage
	tasks, err := s.repo.GetTasksByUsername(context.Background(), username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": tasks})
}
