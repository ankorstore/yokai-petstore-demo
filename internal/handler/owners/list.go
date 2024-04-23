package owners

import (
	"net/http"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/labstack/echo/v4"
)

type ListOwnersHandler struct {
	querier sqlc.Querier
}

func NewListOwnersHandler(querier sqlc.Querier) *ListOwnersHandler {
	return &ListOwnersHandler{
		querier: querier,
	}
}

func (h *ListOwnersHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		owners, err := h.querier.ListOwners(c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, owners)
	}
}
