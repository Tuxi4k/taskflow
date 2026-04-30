package task_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"taskflow/internal/modules/task"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestApp() *fiber.App {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&task.Task{})

	repo := task.NewRepository(db)
	service := task.NewService(repo)
	handler := task.NewHandler(service)

	app := fiber.New()
	handler.RegisterRoutes(app.Group("/tasks"))
	return app
}

func createTask(t *testing.T, app *fiber.App, title string) *task.Task {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"title": title, "status": "todo"})
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)

	var created task.Task
	json.NewDecoder(resp.Body).Decode(&created)
	return &created
}

func TestIntegration_CreateTask(t *testing.T) {
	app := setupTestApp()

	body := `{"title":"Тестовая задача","status":"todo"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)

	var created task.Task
	json.NewDecoder(resp.Body).Decode(&created)
	assert.Equal(t, "Тестовая задача", created.Title)
	assert.Equal(t, task.StatusTodo, created.Status)
	assert.NotZero(t, created.ID)
}

func TestIntegration_CreateTask_EmptyTitle(t *testing.T) {
	app := setupTestApp()

	body := `{"title":""}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 422, resp.StatusCode)
}

func TestIntegration_CreateTask_TooShortTitle(t *testing.T) {
	app := setupTestApp()

	body := `{"title":"ab"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 422, resp.StatusCode)
}

func TestIntegration_CreateTask_WithStatus(t *testing.T) {
	app := setupTestApp()

	body := `{"title":"Срочная","status":"doing"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)

	var created task.Task
	json.NewDecoder(resp.Body).Decode(&created)
	assert.Equal(t, task.StatusDoing, created.Status)
}

func TestIntegration_CreateTask_InvalidStatus(t *testing.T) {
	app := setupTestApp()

	body := `{"title":"Тест","status":"invalid"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 422, resp.StatusCode)
}

func TestIntegration_GetAll_Empty(t *testing.T) {
	app := setupTestApp()

	req, _ := http.NewRequest("GET", "/tasks", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var tasks []task.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	assert.Empty(t, tasks)
}

func TestIntegration_GetAll_TwoTasks(t *testing.T) {
	app := setupTestApp()

	createTask(t, app, "Задача 1")
	createTask(t, app, "Задача 2")

	req, _ := http.NewRequest("GET", "/tasks", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var tasks []task.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	assert.Len(t, tasks, 2)
}

func TestIntegration_GetAll_FilterByStatus(t *testing.T) {
	app := setupTestApp()

	createTask(t, app, "Задача 1")

	task2 := createTask(t, app, "Задача 2")
	body, _ := json.Marshal(map[string]string{"status": "done"})
	req, _ := http.NewRequest("PATCH", "/tasks/"+uintToString(task2.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req)

	req, _ = http.NewRequest("GET", "/tasks?status=done", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var tasks []task.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	assert.Len(t, tasks, 1)
	assert.Equal(t, task.StatusDone, tasks[0].Status)
}

func TestIntegration_GetAll_InvalidStatusFilter(t *testing.T) {
	app := setupTestApp()

	req, _ := http.NewRequest("GET", "/tasks?status=invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 422, resp.StatusCode)
}

func TestIntegration_GetByID_Success(t *testing.T) {
	app := setupTestApp()

	created := createTask(t, app, "Найти меня")

	req, _ := http.NewRequest("GET", "/tasks/"+uintToString(created.ID), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var found task.Task
	json.NewDecoder(resp.Body).Decode(&found)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "Найти меня", found.Title)
}

func TestIntegration_GetByID_NotFound(t *testing.T) {
	app := setupTestApp()

	req, _ := http.NewRequest("GET", "/tasks/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestIntegration_GetByID_InvalidID(t *testing.T) {
	app := setupTestApp()

	req, _ := http.NewRequest("GET", "/tasks/abc", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestIntegration_UpdateTask_Success(t *testing.T) {
	app := setupTestApp()

	created := createTask(t, app, "Старый заголовок")

	body := `{"title":"Новый заголовок","status":"doing"}`
	req, _ := http.NewRequest("PATCH", "/tasks/"+uintToString(created.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var updated task.Task
	json.NewDecoder(resp.Body).Decode(&updated)
	assert.Equal(t, "Новый заголовок", updated.Title)
	assert.Equal(t, task.StatusDoing, updated.Status)
}

func TestIntegration_UpdateTask_OnlyTitle(t *testing.T) {
	app := setupTestApp()

	created := createTask(t, app, "Оригинал")

	body := `{"title":"Только заголовок"}`
	req, _ := http.NewRequest("PATCH", "/tasks/"+uintToString(created.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var updated task.Task
	json.NewDecoder(resp.Body).Decode(&updated)
	assert.Equal(t, "Только заголовок", updated.Title)
	assert.Equal(t, task.StatusTodo, updated.Status)
}

func TestIntegration_UpdateTask_OnlyStatus(t *testing.T) {
	app := setupTestApp()

	created := createTask(t, app, "Без изменений")

	body := `{"status":"done"}`
	req, _ := http.NewRequest("PATCH", "/tasks/"+uintToString(created.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)

	var updated task.Task
	json.NewDecoder(resp.Body).Decode(&updated)
	assert.Equal(t, "Без изменений", updated.Title)
	assert.Equal(t, task.StatusDone, updated.Status)
}

func TestIntegration_UpdateTask_NotFound(t *testing.T) {
	app := setupTestApp()

	body := `{"title":"Неважно"}`
	req, _ := http.NewRequest("PATCH", "/tasks/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestIntegration_DeleteTask_Success(t *testing.T) {
	app := setupTestApp()

	created := createTask(t, app, "Удаляемая задача")

	req, _ := http.NewRequest("DELETE", "/tasks/"+uintToString(created.ID), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 204, resp.StatusCode)

	req, _ = http.NewRequest("GET", "/tasks/"+uintToString(created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestIntegration_DeleteTask_NotFound(t *testing.T) {
	app := setupTestApp()

	req, _ := http.NewRequest("DELETE", "/tasks/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func uintToString(id uint) string {
	return fmt.Sprintf("%d", id)
}
