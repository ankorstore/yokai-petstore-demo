

up:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose up -d

down:
	docker compose down

fresh:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose down --remove-orphans
	docker compose build --no-cache
	docker compose up -d --build -V

logs:
	docker compose logs -f

test:
	go test -v -race -cover -count=1 -failfast ./...

lint:
	golangci-lint run -v

migrate-create:
	docker run -v ./db/migrations:/migrations migrate/migrate create -ext sql -dir migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	docker compose exec petstore-demo-server go run . migrate up

migrate-down:
	docker compose exec petstore-demo-server go run . migrate down

sqlc:
	docker run --rm -v ./:/src -w /src sqlc/sqlc $(filter-out $@,$(MAKECMDGOALS))

%:
    @: