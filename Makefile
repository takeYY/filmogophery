TARGET_COMPOSE ?= compose.yml

.PHONY:
	generate_docs \
	build 		  \
	up 			  \
	up_d		  \
	stop		  \
	down_v		  \
	test_echo	  \
	test_hono     \
	gen_models    \
	mock          \
	start         \
	clean


generate_docs:
	cd docs
	npx @redocly/cli build-docs docs/openapi.yaml --output docs/index.html
	cd ..

build:
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) build --no-cache

up:
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) up

up_d:
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) up -d

stop:
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) stop

down_v:
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) down -v

test_echo:
	make up_d TARGET_COMPOSE=compose.echo.yml
	docker compose -f compose.shared.yml -f compose.echo.yml exec api_echo gotest -v ./... -cover
	make stop TARGET_COMPOSE=compose.echo.yml

test_hono:
	make up_d TARGET_COMPOSE=compose.hono.yml
	docker compose -f compose.shared.yml -f compose.hono.yml exec api_hono bun test
	make stop TARGET_COMPOSE=compose.hono.yml

gen_models:
	make up_d TARGET_COMPOSE=compose.echo.yml
	docker compose -f compose.shared.yml -f compose.echo.yml exec api_echo go run cmd/gen/gorm_gen.go
	make stop TARGET_COMPOSE=compose.echo.yml

mock:
	make up_d TARGET_COMPOSE=compose.echo.yml
	docker compose -f compose.shared.yml -f compose.echo.yml exec api_echo mockgen -package=mocks -destination=tests/mocks/mock_repositories.go filmogophery/internal/app/repositories ITmdbRepository
	docker compose -f compose.shared.yml -f compose.echo.yml exec api_echo mockgen -package=mocks -destination=tests/mocks/mock_services.go filmogophery/internal/app/services ITmdbService
	make stop TARGET_COMPOSE=compose.echo.yml

start:
	make up_d
	cd web-next-js
	npm run dev

clean:
	make down_v TARGET_COMPOSE=compose.echo.yml
	make down_v TARGET_COMPOSE=compose.hono.yml
	docker rm -f api_echo api_hono 2>/dev/null || true
