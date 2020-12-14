package postgres

import (
	"context"
	"sync"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
	"github.com/cage1016/todo/internal/pkg/errors"
)

var _ model.TodoRepository = (*todoRepository)(nil)

var (
	ErrDatabase = errors.New("database internal")
)

type todoRepository struct {
	mu  sync.RWMutex
	log log.Logger
	db  *gorm.DB
}

func New(db *gorm.DB, logger log.Logger) model.TodoRepository {
	return &todoRepository{
		mu:  sync.RWMutex{},
		log: logger,
		db:  db,
	}
}

func (repo *todoRepository) Save(ctx context.Context, todo model.Todo) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if err := repo.db.WithContext(ctx).Save(ModelToDB(todo)).Error; err != nil {
		return err
	}
	return nil
}

func (repo *todoRepository) Complete(ctx context.Context, todoID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	result := repo.db.WithContext(ctx).Model(&Todo{ID: todoID}).UpdateColumn("complete", gorm.Expr("NOT complete"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (repo *todoRepository) List(ctx context.Context, filter string) ([]model.Todo, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	db := repo.db.WithContext(ctx)
	switch filter {
	case "active":
		db = db.Where("complete", false)
		break
	case "complete":
		db = db.Where("complete", true)
		break
	case "all":
		break
	default:

	}

	dTodos := []Todo{}
	if err := db.Find(&dTodos).Order("createdAt").Error; err != nil {
		return []model.Todo{}, nil
	}

	todos := []model.Todo{}
	for _, t := range dTodos {
		todos = append(todos, DBtoModel(t))
	}
	return todos, nil
}

func (repo *todoRepository) Get(ctx context.Context, todoID string) (model.Todo, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	todo := Todo{ID: todoID}
	if err := repo.db.First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Todo{}, errors.Wrap(service.ErrNotFound, err)
		}
		return model.Todo{}, err
	}

	return DBtoModel(todo), nil
}

func (repo *todoRepository) Clear(ctx context.Context) error {
	if err := repo.db.WithContext(ctx).Where("complete", true).Delete(&Todo{}).Error; err != nil {
		return err
	}
	return nil
}
