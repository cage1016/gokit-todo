// +build !integration

package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cage1016/gokit-todo/internal/app/todo/model"
	psql "github.com/cage1016/gokit-todo/internal/app/todo/postgres"
)

func TestTodoRepository_Add(t *testing.T) {
	var (
		mTodo = &model.Todo{
			ID:        "iKe0KxpurIn0E_6vzUDAr",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Text:      "aa",
			Completed: false,
		}
	)

	type fields struct {
		mock sqlmock.Sqlmock
	}

	type args struct {
		todo *model.Todo
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		checkFunc func(err error)
		wantErr   bool
	}{
		{
			name: "Add Todo",
			prepare: func(f *fields) {
				f.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "todos" ("id","created_at","updated_at","text","completed") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(mTodo.ID, mTodo.CreatedAt, mTodo.UpdatedAt, mTodo.Text, mTodo.Completed).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
			args:    args{todo: mTodo},
			wantErr: false,
		},
		{
			name: "Add Todo Fail without primary key",
			prepare: func(f *fields) {
				f.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "todos" ("id","created_at","updated_at","text","completed") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(mTodo.ID, mTodo.CreatedAt, mTodo.UpdatedAt, mTodo.Text, mTodo.Completed).
					WillReturnResult(sqlmock.NewResult(0, 0)).
					WillReturnError(sql.ErrNoRows)
			},
			args: args{todo: func() *model.Todo {
				mTodo.ID = ""
				return mTodo
			}()},
			checkFunc: func(err error) {
				assert.Equal(t, err, sql.ErrNoRows, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			f := fields{
				mock: mock,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := psql.New(gdb, log.NewLogfmtLogger(os.Stderr))

			if err := repo.Add(context.Background(), tt.args.todo); (err != nil) != tt.wantErr {
				t.Errorf("Add(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(err)
				}
			}
		})
	}
}

func TestTodoRepository_List(t *testing.T) {
	var (
		mTodos = []*model.Todo{
			{
				ID:        "iKe0KxpurIn0E_6vzUDAr",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Text:      "aa",
				Completed: false,
			},
			{
				ID:        "zIYPEK0zEpUc7CoQWIGB2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Text:      "bb",
				Completed: false,
			},
		}
	)

	type fields struct {
		mock sqlmock.Sqlmock
	}

	type args struct {
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res []*model.Todo, err error)
	}{
		{
			name: "List Todo",
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "text", "completed", "created_at", "updated_at"})
				for _, m := range mTodos {
					rows.AddRow(m.ID, m.Text, m.Completed, m.CreatedAt, m.UpdatedAt)
				}
				f.mock.ExpectQuery(`SELECT (.+) FROM "todos"`).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			wantErr: false,
			checkFunc: func(res []*model.Todo, err error) {
				assert.Equal(t, mTodos, res, fmt.Sprintf("models: expected aa got %v", res))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			f := fields{
				mock: mock,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := psql.New(gdb, log.NewLogfmtLogger(os.Stderr))

			if res, err := repo.List(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("List(ctx context.Context) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(res, err)
				}
			}
		})
	}
}

func TestTodoRepository_Get(t *testing.T) {
	var (
		mTodo = &model.Todo{
			ID:        "iKe0KxpurIn0E_6vzUDAr",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Text:      "aa",
			Completed: false,
		}
	)

	type fields struct {
		mock sqlmock.Sqlmock
	}

	type args struct {
		todoID string
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		checkFunc func(res *model.Todo, err error)
		wantErr   bool
	}{
		{
			name: "Get Todo",
			prepare: func(f *fields) {
				rows := sqlmock.NewRows([]string{"id", "text", "completed", "created_at", "updated_at"})
				rows.AddRow(mTodo.ID, mTodo.Text, mTodo.Completed, mTodo.CreatedAt, mTodo.UpdatedAt)

				f.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE "id" = $1`)).
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(rows).
					WillReturnError(nil)
			},
			args:    args{todoID: mTodo.ID},
			wantErr: false,
			checkFunc: func(res *model.Todo, err error) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			f := fields{
				mock: mock,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := psql.New(gdb, log.NewLogfmtLogger(os.Stderr))

			if res, err := repo.Get(context.Background(), tt.args.todoID); (err != nil) != tt.wantErr {
				t.Errorf("Get(ctx context.Context id string) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(res, err)
				}
			}
		})
	}
}

func TestTodoRepository_Delete(t *testing.T) {
	var (
		mTodo = &model.Todo{
			ID:        "iKe0KxpurIn0E_6vzUDAr",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Text:      "aa",
			Completed: false,
		}
	)

	type fields struct {
		mock sqlmock.Sqlmock
	}

	type args struct {
		todoID string
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		checkFunc func(err error)
		wantErr   bool
	}{
		{
			name: "Delete Todo",
			prepare: func(f *fields) {
				f.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "todos" WHERE "todos"."id" = $1`)).
					WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewErrorResult(nil))
			},
			args:    args{todoID: mTodo.ID},
			wantErr: false,
		},
		{
			name: "Delete Todo fail not found",
			prepare: func(f *fields) {
				f.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "todos" WHERE "todos"."id" = $1`)).
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 0)).
					WillReturnError(sql.ErrNoRows)
			},
			args:    args{todoID: mTodo.ID},
			wantErr: true,
			checkFunc: func(err error) {
				assert.Equal(t, err, sql.ErrNoRows, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			f := fields{
				mock: mock,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := psql.New(gdb, log.NewLogfmtLogger(os.Stderr))

			if err := repo.Delete(context.Background(), tt.args.todoID); (err != nil) != tt.wantErr {
				t.Errorf("Delete(ctx context.Context id string) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(err)
				}
			}
		})
	}
}

func TestTodoRepository_Update(t *testing.T) {
	var (
		mTodo = &model.Todo{
			ID:        "iKe0KxpurIn0E_6vzUDAr",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Text:      "aa",
			Completed: false,
		}
	)

	type fields struct {
		mock sqlmock.Sqlmock
	}

	type args struct {
		todo *model.Todo
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		checkFunc func(err error)
		wantErr   bool
	}{
		{
			name: "Update Todo",
			prepare: func(f *fields) {
				f.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args:    args{todo: mTodo},
			wantErr: false,
		},
		{
			name: "Update Todo fail with wrong type",
			prepare: func(f *fields) {
				f.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 0)).
					WillReturnError(sql.ErrNoRows)
			},
			args:    args{todo: mTodo},
			wantErr: true,
			checkFunc: func(err error) {
				assert.Equal(t, err, sql.ErrNoRows, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			f := fields{
				mock: mock,
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			repo := psql.New(gdb, log.NewLogfmtLogger(os.Stderr))

			if err := repo.Update(context.Background(), tt.args.todo); (err != nil) != tt.wantErr {
				t.Errorf("Update(ctx context.Context todo *model.Todo) error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(err)
				}
			}
		})
	}
}
