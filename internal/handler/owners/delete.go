package owners

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/ankorstore/yokai/config"
	"github.com/labstack/echo/v4"
)

type DeleteOwnerHandler struct {
	config  *config.Config
	queries *sqlc.Queries
}

func NewDeleteOwnerHandler(config *config.Config, queries *sqlc.Queries) *DeleteOwnerHandler {
	return &DeleteOwnerHandler{
		config:  config,
		queries: queries,
	}
}

func (h *DeleteOwnerHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid owner id: %v", err))
		}

		ctx := c.Request().Context()

		err = h.queries.DeleteOwner(ctx, int64(id))
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
