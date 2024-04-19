package owners

import (
	"net/http"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

type ListOwnersHandler struct {
	config  *config.Config
	queries *sqlc.Queries
}

func NewListOwnersHandler(config *config.Config, queries *sqlc.Queries) *ListOwnersHandler {
	return &ListOwnersHandler{
		config:  config,
		queries: queries,
	}
}

func (h *ListOwnersHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		owners, err := h.queries.ListOwners(c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, owners)
	}
}
