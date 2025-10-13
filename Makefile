.PHONY:
	build 		\
	up 			\
	up_d		\
	exec_api	\
	stop		\
	down_v		\
	test		\
	gen_models  \
	start


build:
	docker compose build --no-cache

up:
	docker compose up

up_d:
	docker compose up -d

exec_api:
	docker exec -it golang_tutorial_api bash

stop:
	docker compose stop

down_v:
	docker compose down -v

test:
	make up_d
	docker compose exec api gotest -v ./... -cover
	docker compose stop

gen_models:
	make up_d
	docker compose exec api go run cmd/gen/gorm_gen.go
	docker compose stop

start:
	make up_d
	cd filmogophery-front
	npm run dev
