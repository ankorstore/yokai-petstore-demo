package internal

import (
	"github.com/ankorstore/yokai-petstore-demo/internal/handler/owners"
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

// Router is used to register the application HTTP middlewares and handlers.
func Router() fx.Option {
	return fx.Options(
		fxhttpserver.AsHandler("GET", "/owners", owners.NewListOwnersHandler),
		fxhttpserver.AsHandler("POST", "/owners", owners.NewCreateOwnerHandler),
		fxhttpserver.AsHandler("GET", "/owners/:id", owners.NewGetOwnerHandler),
		fxhttpserver.AsHandler("DELETE", "/owners/:id", owners.NewDeleteOwnerHandler),
	)
}
