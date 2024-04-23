package pets

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/labstack/echo/v4"
)

type DeleteOwnerPetHandler struct {
	querier sqlc.Querier
}

func NewDeleteOwnerPetHandler(querier sqlc.Querier) *DeleteOwnerPetHandler {
	return &DeleteOwnerPetHandler{
		querier: querier,
	}
}

func (h *DeleteOwnerPetHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		ownerId, err := strconv.Atoi(c.Param("owner_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid owner id: %v", err))
		}

		petId, err := strconv.Atoi(c.Param("pet_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid pet id: %v", err))
		}

		err = h.querier.DeleteOwnerPet(ctx, sqlc.DeleteOwnerPetParams{
			OwnerID: sql.NullInt64{Int64: int64(ownerId), Valid: true},
			ID:      int64(petId),
		})
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
