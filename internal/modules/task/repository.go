package task

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(status *Status) ([]Task, error) {
	var tasks []Task
	query := r.db.Order("created_at DESC")

	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	err := query.Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) GetByID(id uint) (*Task, error) {
	var task Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) Create(task *Task) error {
	return r.db.Create(task).Error
}

func (r *Repository) Update(task *Task) error {
	return r.db.Select("title", "status").Updates(task).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Task{}, id).Error
}
