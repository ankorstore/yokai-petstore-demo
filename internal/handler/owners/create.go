package owners

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/labstack/echo/v4"
)

type CreateOwnerParams struct {
	Name string `json:"name" form:"name" query:"name"`
	Bio  string `json:"bio" form:"bio" query:"bio"`
}

type CreateOwnerHandler struct {
	querier sqlc.Querier
}

func NewCreateOwnerHandler(querier sqlc.Querier) *CreateOwnerHandler {
	return &CreateOwnerHandler{
		querier: querier,
	}
}

func (h *CreateOwnerHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		params := new(CreateOwnerParams)
		if err := c.Bind(params); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid parameters: %v", err))
		}

		result, err := h.querier.CreateOwner(
			ctx,
			sqlc.CreateOwnerParams{
				Name: params.Name,
				Bio:  sql.NullString{String: params.Bio, Valid: true},
			})
		if err != nil {
			return err
		}

		ownerId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		owner, err := h.querier.GetOwner(ctx, int32(ownerId))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, owner)
	}
}
