package task

import "gorm.io/gorm"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(input Input) (*Task, error) {
	err := input.Validate(true)
	if err != nil {
		return nil, err
	}

	task := &Task{Title: *input.Title, Status: *input.Status}
	err = s.repo.Create(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) GetAll(status *Status) ([]Task, error) {
	return s.repo.GetAll(status)
}

func (s *Service) GetByID(id uint) (*Task, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id uint, input Input) (*Task, error) {
	err := input.Validate(false)
	if err != nil {
		return nil, err
	}

	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		task.Title = *input.Title
	}

	if input.Status != nil {
		task.Status = *input.Status
	}

	err = s.repo.Update(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) Delete(id uint) error {
	result := s.repo.db.Delete(&Task{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}
