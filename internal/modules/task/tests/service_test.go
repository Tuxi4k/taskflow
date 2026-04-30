package task_test

import (
	"testing"

	"github.com/Tuxi4k/taskflow/internal/modules/task"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupService() *task.Service {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&task.Task{})
	repo := task.NewRepository(db)
	return task.NewService(repo)
}

func ptr(s string) *string                 { return &s }
func statusPtr(s task.Status) *task.Status { return &s }

func TestService_Create_Success(t *testing.T) {
	svc := setupService()

	input := task.Input{
		Title:  ptr("Купить молоко"),
		Status: statusPtr(task.StatusTodo),
	}

	result, err := svc.Create(input)

	assert.NoError(t, err)
	assert.Equal(t, "Купить молоко", result.Title)
	assert.Equal(t, task.StatusTodo, result.Status)
	assert.NotZero(t, result.ID)
	assert.False(t, result.CreatedAt.IsZero())
}

func TestService_Create_DefaultStatusTodo(t *testing.T) {
	svc := setupService()

	input := task.Input{
		Title:  ptr("Без указания статуса"),
		Status: statusPtr(task.StatusTodo),
	}

	result, err := svc.Create(input)

	assert.NoError(t, err)
	assert.Equal(t, task.StatusTodo, result.Status)
}

func TestService_Update_Success(t *testing.T) {
	svc := setupService()

	created, _ := svc.Create(task.Input{
		Title:  ptr("Старый заголовок"),
		Status: statusPtr(task.StatusTodo),
	})

	updated, err := svc.Update(created.ID, task.Input{
		Title:  ptr("Новый заголовок"),
		Status: statusPtr(task.StatusDoing),
	})

	assert.NoError(t, err)
	assert.Equal(t, "Новый заголовок", updated.Title)
	assert.Equal(t, task.StatusDoing, updated.Status)
	assert.Equal(t, created.ID, updated.ID)
}

func TestService_Update_OnlyTitle(t *testing.T) {
	svc := setupService()

	created, _ := svc.Create(task.Input{
		Title:  ptr("Оригинал"),
		Status: statusPtr(task.StatusTodo),
	})

	updated, err := svc.Update(created.ID, task.Input{
		Title:  ptr("Только заголовок изменился"),
		Status: nil,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Только заголовок изменился", updated.Title)
	assert.Equal(t, task.StatusTodo, updated.Status)
}

func TestService_Update_OnlyStatus(t *testing.T) {
	svc := setupService()

	created, _ := svc.Create(task.Input{
		Title:  ptr("Без изменений"),
		Status: statusPtr(task.StatusTodo),
	})

	updated, err := svc.Update(created.ID, task.Input{
		Title:  nil,
		Status: statusPtr(task.StatusDone),
	})

	assert.NoError(t, err)
	assert.Equal(t, "Без изменений", updated.Title)
	assert.Equal(t, task.StatusDone, updated.Status)
}

func TestService_Update_NotFound(t *testing.T) {
	svc := setupService()

	_, err := svc.Update(999, task.Input{
		Title:  ptr("Неважно"),
		Status: statusPtr(task.StatusTodo),
	})

	assert.Error(t, err)
}

func TestValidate_Title_Empty(t *testing.T) {
	svc := setupService()

	_, err := svc.Create(task.Input{
		Title:  ptr(""),
		Status: statusPtr(task.StatusTodo),
	})

	assert.Error(t, err)
}

func TestValidate_Title_TooShort(t *testing.T) {
	svc := setupService()

	_, err := svc.Create(task.Input{
		Title:  ptr("ab"),
		Status: statusPtr(task.StatusTodo),
	})

	assert.Error(t, err)
}

func TestValidate_Title_MinBoundary(t *testing.T) {
	svc := setupService()

	_, err := svc.Create(task.Input{
		Title:  ptr("abc"),
		Status: statusPtr(task.StatusTodo),
	})

	assert.NoError(t, err)
}

func TestValidate_Title_TooLong(t *testing.T) {
	svc := setupService()

	longTitle := ""
	for i := 0; i < 201; i++ {
		longTitle += "a"
	}

	_, err := svc.Create(task.Input{
		Title:  ptr(longTitle),
		Status: statusPtr(task.StatusTodo),
	})

	assert.Error(t, err)
}

func TestValidate_Title_MaxBoundary(t *testing.T) {
	svc := setupService()

	title := ""
	for i := 0; i < 200; i++ {
		title += "a"
	}

	_, err := svc.Create(task.Input{
		Title:  ptr(title),
		Status: statusPtr(task.StatusTodo),
	})

	assert.NoError(t, err)
}

func TestValidate_Status_Invalid(t *testing.T) {
	svc := setupService()

	invalidStatus := task.Status("invalid")
	_, err := svc.Create(task.Input{
		Title:  ptr("Нормальный заголовок"),
		Status: &invalidStatus,
	})

	assert.Error(t, err)
}

func TestValidate_Status_AllValid(t *testing.T) {
	svc := setupService()

	for _, status := range []task.Status{task.StatusTodo, task.StatusDoing, task.StatusDone} {
		_, err := svc.Create(task.Input{
			Title:  ptr("Задача"),
			Status: &status,
		})
		assert.NoError(t, err, "Статус %s должен быть валидным", status)
	}
}
