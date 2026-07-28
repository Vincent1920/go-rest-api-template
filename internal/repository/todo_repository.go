package repository

import (
	"todo-api/internal/domain"

	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) domain.TodoRepository {
	return &todoRepository{
		db: db,
	}
}

func (r *todoRepository) Create(todo *domain.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) GetByID(id uint) (*domain.Todo, error) {
	var todo domain.Todo

	err := r.db.Preload("User").First(&todo, id).Error
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *todoRepository) GetAll(userID uint, page, limit int, search string) ([]domain.Todo, int64, error) {
	var todos []domain.Todo
	var total int64

	query := r.db.Model(&domain.Todo{}).Where("user_id = ?", userID)

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err = query.
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&todos).Error

	if err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (r *todoRepository) Update(todo *domain.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Todo{}, id).Error
}
