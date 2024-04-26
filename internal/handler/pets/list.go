package pets

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/labstack/echo/v4"
)

type ListOwnerPetsHandler struct {
	querier sqlc.Querier
}

func NewListOwnerPetsHandler(querier sqlc.Querier) *ListOwnerPetsHandler {
	return &ListOwnerPetsHandler{
		querier: querier,
	}
}

func (h *ListOwnerPetsHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ownerId, err := strconv.Atoi(c.Param("owner_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid owner id: %v", err))
		}

		pets, err := h.querier.ListOwnerPets(
			c.Request().Context(),
			sql.NullInt32{Int32: int32(ownerId), Valid: true},
		)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, pets)
	}
}
