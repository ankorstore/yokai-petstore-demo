package owners

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

type GetOwnerHandler struct {
	config  *config.Config
	queries *sqlc.Queries
}

func NewGetOwnerHandler(config *config.Config, queries *sqlc.Queries) *GetOwnerHandler {
	return &GetOwnerHandler{
		config:  config,
		queries: queries,
	}
}

func (h *GetOwnerHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid owner id: %v", err))
		}

		ctx := c.Request().Context()

		owner, err := h.queries.GetOwner(ctx, int64(id))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, owner)
	}
}
