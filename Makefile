.PHONY:
	build 		\
	up 			\
	up_d		\
	exec_api	\
	stop		\
	down_v		\
	test		\
	gen_models  \
	mock        \
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
	docker compose exec api gotest -v ./tests/... -cover
	docker compose stop

gen_models:
	make up_d
	docker compose exec api go run cmd/gen/gorm_gen.go
	docker compose stop

mock:
	make up_d
	docker compose exec api mockgen -package=mocks -destination=tests/mocks/mock_repositories.go filmogophery/internal/app/repositories IGenreRepository,IImpressionRepository,IMediaRepository,IMovieRepository,IRecordRepository
	docker compose exec api mockgen -package=mocks -destination=tests/mocks/mock_services.go filmogophery/internal/app/services IMovieService
	docker compose stop

start:
	make up_d
	cd filmogophery-front
	npm run dev
