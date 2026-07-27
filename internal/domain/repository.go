package domain

type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

type TodoRepository interface {
	Create(todo *Todo) error
	GetByID(id uint) (*Todo, error)
	GetAll(userID uint, page, limit int, search string) ([]Todo, int64, error)
	Update(todo *Todo) error
	Delete(id uint) error
}
