package owners

import (
	"database/sql"
	"net/http"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

type CreateOwnerParams struct {
	Name string `json:"name" form:"name" query:"name"`
	Bio  string `json:"bio" form:"bio" query:"bio"`
}

type CreateOwnerHandler struct {
	config  *config.Config
	queries *sqlc.Queries
}

func NewCreateOwnerHandler(config *config.Config, queries *sqlc.Queries) *CreateOwnerHandler {
	return &CreateOwnerHandler{
		config:  config,
		queries: queries,
	}
}

func (h *CreateOwnerHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		params := new(CreateOwnerParams)
		if err := c.Bind(params); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		result, err := h.queries.CreateOwner(
			ctx,
			sqlc.CreateOwnerParams{
				Name: params.Name,
				Bio:  sql.NullString{String: params.Bio, Valid: true},
			})
		if err != nil {
			return err
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		owner, err := h.queries.GetOwner(ctx, lastInsertId)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, owner)
	}
}
