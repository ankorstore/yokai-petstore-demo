package internal

import (
	"github.com/ankorstore/yokai-petstore-demo/internal/handler"
	"github.com/ankorstore/yokai-petstore-demo/internal/handler/owners"
	"github.com/ankorstore/yokai-petstore-demo/internal/handler/pets"
	"github.com/ankorstore/yokai/fxhttpserver"
	"go.uber.org/fx"
)

// Router is used to register the application HTTP middlewares and handlers.
func Router() fx.Option {
	return fx.Options(
		// example
		fxhttpserver.AsHandler("GET", "", handler.NewExampleHandler),
		// owners
		fxhttpserver.AsHandler("GET", "/owners", owners.NewListOwnersHandler),
		fxhttpserver.AsHandler("POST", "/owners", owners.NewCreateOwnerHandler),
		fxhttpserver.AsHandler("GET", "/owners/:owner_id", owners.NewGetOwnerHandler),
		fxhttpserver.AsHandler("DELETE", "/owners/:owner_id", owners.NewDeleteOwnerHandler),
		// pets
		fxhttpserver.AsHandler("GET", "/owners/:owner_id/pets", pets.NewListOwnerPetsHandler),
		fxhttpserver.AsHandler("POST", "/owners/:owner_id/pets", pets.NewCreateOwnerPetsHandler),
		fxhttpserver.AsHandler("GET", "/owners/:owner_id/pets/:pet_id", pets.NewGetOwnerPetHandler),
		fxhttpserver.AsHandler("DELETE", "/owners/:owner_id/pets/:pet_id", pets.NewDeleteOwnerPetHandler),
	)
}
