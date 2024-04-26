package pets

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/labstack/echo/v4"
)

type CreateOwnerPetParams struct {
	Name    string `json:"name" form:"name" query:"name"`
	Type    string `json:"type" form:"type" query:"type"`
	OwnerId int    `json:"owner_id" form:"owner_id" query:"owner_id"`
}

type CreateOwnerPetsHandler struct {
	querier sqlc.Querier
}

func NewCreateOwnerPetsHandler(querier sqlc.Querier) *CreateOwnerPetsHandler {
	return &CreateOwnerPetsHandler{
		querier: querier,
	}
}

func (h *CreateOwnerPetsHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		ownerId, err := strconv.Atoi(c.Param("owner_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid owner id: %v", err))
		}

		params := new(CreateOwnerPetParams)
		if err := c.Bind(params); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid parameters: %v", err))
		}

		result, err := h.querier.CreateOwnerPet(
			ctx,
			sqlc.CreateOwnerPetParams{
				Name:    params.Name,
				Type:    params.Type,
				OwnerID: sql.NullInt32{Int32: int32(ownerId), Valid: true},
			})
		if err != nil {
			return err
		}

		petId, err := result.LastInsertId()
		if err != nil {
			return err
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
