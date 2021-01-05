package postgres

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"github.com/cage1016/gokit-todo/internal/app/todo/model"
)

var _ model.TodoRepository = (*todoRepository)(nil)

type todoRepository struct {
	mu  sync.RWMutex
	log log.Logger
	db  *gorm.DB
}

func (repo *todoRepository) Get(ctx context.Context, todoID string) (res *model.Todo, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	res = new(model.Todo)
	err = repo.db.WithContext(ctx).Where("id", todoID).Find(res).Error
	return
}

func (repo *todoRepository) Add(ctx context.Context, todo *model.Todo) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if err := repo.db.WithContext(ctx).Create(todo).Error; err != nil {
		return err
	}
	return nil
}

func (repo *todoRepository) Delete(ctx context.Context, todoID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if err := repo.db.WithContext(ctx).Delete(&model.Todo{ID: todoID}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *todoRepository) Update(ctx context.Context, todo *model.Todo) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	result := repo.db.WithContext(ctx).Model(&model.Todo{ID: todo.ID}).UpdateColumns(
		map[string]interface{}{
			"text":       todo.Text,
			"completed":  todo.Completed,
			"updated_at": todo.UpdatedAt,
		},
	)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *todoRepository) List(ctx context.Context) (res []*model.Todo, err error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	err = repo.db.WithContext(ctx).Order("created_at desc").Find(&res).Error
	return
}

func New(db *gorm.DB, logger log.Logger) model.TodoRepository {
	return &todoRepository{
		mu:  sync.RWMutex{},
		log: logger,
		db:  db,
	}
}
