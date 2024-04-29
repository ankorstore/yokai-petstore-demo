package owners

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/labstack/echo/v4"
)

type GetOwnerHandler struct {
	querier sqlc.Querier
}

func NewGetOwnerHandler(querier sqlc.Querier) *GetOwnerHandler {
	return &GetOwnerHandler{
		querier: querier,
	}
}

func (h *GetOwnerHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		ownerId, err := strconv.Atoi(c.Param("owner_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid owner id: %v", err))
		}

		owner, err := h.querier.GetOwner(ctx, int32(ownerId))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("owner with id %d not found", ownerId))
			}

			return err
		}

		return c.JSON(http.StatusOK, owner)
	}
}
