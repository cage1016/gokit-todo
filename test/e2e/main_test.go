// +build e2e

package e2e

import (
	"net/http"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"gorm.io/gorm"

	"github.com/cage1016/todo/internal/app/todo/endpoints"
	"github.com/cage1016/todo/internal/app/todo/postgres"
	"github.com/cage1016/todo/internal/app/todo/service"
	transports "github.com/cage1016/todo/internal/app/todo/transports/http"
	"github.com/cage1016/todo/internal/pkg/errors"
)

const (
	// databaseHost is the host name of the test database.
	// databaseHost = "db"
	databaseHost = "localhost"

	// databasePort is the port that the test database is listening on.
	databasePort = "5432"

	// databaseUser is the user for the test database.
	databaseUser = "postgres"

	// databasePass is the password of the user for the test database.
	databasePass = "password"

	// databaseName is the name of the test database.
	databaseName = "todo"
)

var a *Application

type Application struct {
	DB      *gorm.DB
	handler http.Handler
}

func Truncate(dbc *gorm.DB) error {
	stmt := "TRUNCATE TABLE todos"

	if err := dbc.Exec(stmt).Error; err != nil {
		return errors.Wrap(errors.New("truncate test database tables"), err)
	}

	return nil
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	logger := log.NewLogfmtLogger(os.Stderr)

	db, err := postgres.Connect(postgres.Config{
		Host:        databaseHost,
		Port:        databasePort,
		User:        databaseUser,
		Pass:        databasePass,
		Name:        databaseName,
		SSLMode:     "disable",
		SSLCert:     "",
		SSLKey:      "",
		SSLRootCert: "",
	})
	if err != nil {
		logger.Log("err", err)
		return 1
	}

	zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
	tracer := opentracing.GlobalTracer()

	repo := postgres.New(db, logger)
	svc := service.New(repo, logger)
	eps := endpoints.New(svc, logger, tracer, zkt)

	a = &Application{
		DB:      db,
		handler: transports.NewHTTPHandler(eps, tracer, zkt, logger),
	}

	return m.Run()
}
