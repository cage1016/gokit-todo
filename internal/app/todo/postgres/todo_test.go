// +build !integration

package postgres_test

import (
	"context"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/kit/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cage1016/todo/internal/app/todo/model"
	psql "github.com/cage1016/todo/internal/app/todo/postgres"
)

func TestTodoRepository_Add(t *testing.T) {
	var (
		mTodo = &model.Todo{
			ID:        "iKe0KxpurIn0E_6vzUDAr",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Text:      "aa",
			Complete:  false,
		}
	)

	type fields struct {
		mock sqlmock.Sqlmock
	}

	type args struct {
		todo *model.Todo
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "Add Todo",
			prepare: func(f *fields) {
				f.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "todos" ("id","created_at","updated_at","text","complete") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(mTodo.ID, mTodo.CreatedAt, mTodo.UpdatedAt, mTodo.Text, mTodo.Complete).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
			args:    args{todo: mTodo},
			wantErr: false,
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
			}
		})
	}
}
