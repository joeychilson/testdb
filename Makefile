ifneq (,$(wildcard ./.env))
    include .env
    export
endif

up:
	docker-compose up -d

down:
	docker-compose down

stop:
	docker-compose stop

dbup:
	dbmate -u ${POSTGRES_URL} -d ./migrations up

dbdown:
	dbmate -u ${POSTGRES_URL} -d ./migrations down

dbnew:
	dbmate -d ./migrations new $(name)

sqlc:
	sqlc generate