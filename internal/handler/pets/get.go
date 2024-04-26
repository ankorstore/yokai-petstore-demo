package pets

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/labstack/echo/v4"
)

type GetOwnerPetHandler struct {
	querier sqlc.Querier
}

func NewGetOwnerPetHandler(querier sqlc.Querier) *GetOwnerPetHandler {
	return &GetOwnerPetHandler{
		querier: querier,
	}
}

func (h *GetOwnerPetHandler) Handle() echo.HandlerFunc {
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

		pet, err := h.querier.GetOwnerPet(ctx, sqlc.GetOwnerPetParams{
			OwnerID: sql.NullInt32{Int32: int32(ownerId), Valid: true},
			ID:      int32(petId),
		})
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, pet)
	}
}
