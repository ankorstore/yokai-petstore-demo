package owners

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankorstore/yokai-petstore-demo/db/sqlc"
	"github.com/labstack/echo/v4"
)

type DeleteOwnerHandler struct {
	querier sqlc.Querier
}

func NewDeleteOwnerHandler(querier sqlc.Querier) *DeleteOwnerHandler {
	return &DeleteOwnerHandler{
		querier: querier,
	}
}

func (h *DeleteOwnerHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		ownerId, err := strconv.Atoi(c.Param("owner_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid owner id: %v", err))
		}

		err = h.querier.DeleteOwner(ctx, int32(ownerId))
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
