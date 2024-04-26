package owners_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/ankorstore/yokai-petstore-demo/internal"
	"github.com/ankorstore/yokai-petstore-demo/internal/module/fxdatabase"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestListOwnersHandlerSuccess(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var querier sqlc.Querier

	internal.RunTest(
		t,
		fxdatabase.RunFxDatabaseMigration(fxdatabase.Up, false),
		fx.Populate(&httpServer, &logBuffer, &traceExporter, &querier),
	)

	// populate database
	_, err := querier.CreateOwner(context.Background(), sqlc.CreateOwnerParams{
		Name: "test name",
		Bio:  sql.NullString{String: "test bio", Valid: true},
	})
	assert.NoError(t, err)

	// [GET] /owners response assertion
	req := httptest.NewRequest(http.MethodGet, "/owners", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"name": "test name"`)
}

func TestListOwnersHandlerSuccessAgain(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var querier sqlc.Querier

	internal.RunTest(
		t,
		fxdatabase.RunFxDatabaseMigration(fxdatabase.Up, false),
		fx.Populate(&httpServer, &logBuffer, &traceExporter, &querier),
	)

	// populate database
	_, err := querier.CreateOwner(context.Background(), sqlc.CreateOwnerParams{
		Name: "test name2",
		Bio:  sql.NullString{String: "test bio2", Valid: true},
	})
	assert.NoError(t, err)

	// [GET] /owners response assertion
	req := httptest.NewRequest(http.MethodGet, "/owners", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"name": "test name2"`)
}
