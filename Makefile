TARGET_COMPOSE ?= compose.yml

.PHONY:
	generate_docs \
	build 		  \
	up 			  \
	up_d		  \
	exec_api	  \
	stop		  \
	down_v		  \
	test_echo	  \
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

exec_api:
	docker exec -it golang_tutorial_api bash

stop:
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) stop

down_v:
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) down -v

test_echo:
	make up_d
	docker compose -f compose.shared.yml -f compose.echo.yml exec api gotest -v ./... -cover
	docker compose -f compose.shared.yml -f compose.echo.yml stop

gen_models:
	make up_d
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) exec api go run cmd/gen/gorm_gen.go
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) stop

mock:
	make up_d
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) exec api mockgen -package=mocks -destination=tests/mocks/mock_repositories.go filmogophery/internal/app/repositories ITmdbRepository
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) exec api mockgen -package=mocks -destination=tests/mocks/mock_services.go filmogophery/internal/app/services ITmdbService
	docker compose -f compose.shared.yml -f $(TARGET_COMPOSE) stop

start:
	make up_d
	cd web-next-js
	npm run dev

clean:
	docker compose -f compose.shared.yml -f compose.echo.yml down -v
	docker compose -f compose.shared.yml -f compose.hono.yml down -v
	docker rm -f api_echo api_hono 2>/dev/null || true
