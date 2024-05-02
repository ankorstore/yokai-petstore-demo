package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ExampleHandler struct {
	db *sql.DB
}

func NewExampleHandler(db *sql.DB) *ExampleHandler {
	return &ExampleHandler{
		db: db,
	}
}

func (h *ExampleHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := h.db.BeginTx(ctx, nil)
		defer tx.Rollback()

		stmt, err := tx.PrepareContext(c.Request().Context(), "UPDATE pets SET name = $1")
		if err != nil {
			return err
		}
		defer stmt.Close()

		res, err := stmt.Exec("generic")
		if err != nil {
			return err
		}

		lastInsertedId, err := res.LastInsertId()
		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		return c.String(
			http.StatusOK,
			fmt.Sprintf("result: lastInsertedId: %v, rowsAffected: %v", lastInsertedId, rowsAffected))
	}
}
