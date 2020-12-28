package postgres

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
)

var _ model.TodoRepository = (*todoRepository)(nil)

type todoRepository struct {
	mu  sync.RWMutex
	log log.Logger
	db  *gorm.DB
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

	result := repo.db.WithContext(ctx).Model(&model.Todo{ID: todo.ID}).UpdateColumn("text", todo.Text)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (repo *todoRepository) CompleteAll(ctx context.Context) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var nonCompleteCount int64
	if err := repo.db.WithContext(ctx).Model(&model.Todo{}).Where("complete = ?", false).Count(&nonCompleteCount).Error; err != nil {
		return err
	}

	if nonCompleteCount > 0 {
		if err := repo.db.WithContext(ctx).Exec("UPDATE todos SET complete = true").Error; err != nil {
			return err
		}
	} else {
		if err := repo.db.WithContext(ctx).Exec("UPDATE todos SET complete = false").Error; err != nil {
			return err
		}
	}
	return nil
}

func (repo *todoRepository) Complete(ctx context.Context, todoID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	result := repo.db.WithContext(ctx).Model(&model.Todo{ID: todoID}).UpdateColumn("complete", gorm.Expr("NOT complete"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (repo *todoRepository) List(ctx context.Context, filter string) (res []*model.Todo, err error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	db := repo.db.WithContext(ctx)
	switch filter {
	case service.ACTIVE:
		db = db.Where("complete", false)
		break
	case service.COMPLETE:
		db = db.Where("complete", true)
		break
	case service.ALL:
		break
	default:

	}

	err = db.Order("created_at").Find(&res).Error
	return
}

func (repo *todoRepository) Clear(ctx context.Context) error {
	if err := repo.db.WithContext(ctx).Where("complete", true).Delete(&model.Todo{}).Error; err != nil {
		return err
	}
	return nil
}

func New(db *gorm.DB, logger log.Logger) model.TodoRepository {
	return &todoRepository{
		mu:  sync.RWMutex{},
		log: logger,
		db:  db,
	}
}
