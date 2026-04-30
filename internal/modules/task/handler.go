package task

import (
	"errors"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type request struct {
	Title  *string `json:"title"`
	Status *Status `json:"status"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app fiber.Router) {
	app.Post("/", h.Create)
	app.Get("/", h.GetAll)
	app.Get("/:id", h.GetByID)
	app.Patch("/:id", h.Update)
	app.Delete("/:id", h.Delete)
}

// @Summary  Создать задачу
// @Tags     tasks
// @Accept   json
// @Produce  json
// @Param    task body request true "Данные задачи"
// @Success  201 {object} Task
// @Failure  422  {object}  map[string]string
// @Router   /tasks [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "невалидный JSON"})
	}

	input := Input{
		Title:  req.Title,
		Status: req.Status,
	}

	task, err := h.service.Create(input)
	if err != nil {
		if errs, ok := err.(validation.Errors); ok {
			return c.Status(422).JSON(errs)
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(task)
}

// @Summary      Получить все задачи
// @Description  Возвращает список всех задач, отсортированных по дате создания
// @Tags         tasks
// @Produce      json
// @Success      200  {array}   Task
// @Failure      422  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks [get]
func (h *Handler) GetAll(c *fiber.Ctx) error {
	var status *Status

	statusParam := Status(c.Query("status"))
	if statusParam != "" {
		if !statusParam.IsValid() {

			return c.Status(422).JSON(fiber.Map{"status": "не корректный статус"})
		}

		status = &statusParam
	}

	tasks, err := h.service.GetAll(status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tasks)
}

// @Summary      Получить задачу по ID
// @Description  Возвращает задачу по её уникальному идентификатору
// @Tags         tasks
// @Produce      json
// @Param        id path int true "ID задачи"
// @Success      200  {object}  Task
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /tasks/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "невалидный id"})
	}

	task, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

// @Summary      Обновить задачу
// @Description  Частично обновляет заголовок и/или статус задачи. Поддерживает PATCH.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID задачи"
// @Param        task body task.request false "Обновляемые поля"
// @Success      200  {object}  Task
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      422  {object}  map[string]string
// @Router       /tasks/{id} [patch]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "невалидный id"})
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "невалидный JSON"})
	}

	input := Input{
		Title:  req.Title,
		Status: req.Status,
	}

	task, err := h.service.Update(uint(id), input)
	if err != nil {
		if errs, ok := err.(validation.Errors); ok {
			return c.Status(422).JSON(errs)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

// @Summary      Удалить задачу
// @Description  Удаляет задачу по её идентификатору
// @Tags         tasks
// @Produce      json
// @Param        id path int true "ID задачи"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /tasks/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "невалидный id"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}
