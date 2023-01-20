ifneq (,$(wildcard ./.env))
    include .env
    export
endif

myup:
	docker-compose -f "docker-compose.mysql.yml" up -d

mydown:
	docker-compose -f "docker-compose.mysql.yml" down

mystop:
	docker-compose -f "docker-compose.mysql.yml" stop

mydbup:
	dbmate -u ${MYSQL_URL} -d ./migrations/mysql up

mydbdown:
	dbmate -u ${MYSQL_URL} -d ./migrations/mysql down

mydbnew:
	dbmate -d ./migrations/mysql new $(name)

pgup:
	docker-compose -f "docker-compose.pgsql.yml" up -d

pgdown:
	docker-compose -f "docker-compose.pgsql.yml" down

pgstop:
	docker-compose -f "docker-compose.pgsql.yml" stop

pgdbup:
	dbmate -u ${POSTGRES_URL} -d ./migrations/pgsql up

pgdbdown:
	dbmate -u ${POSTGRES_URL} -d ./migrations/pgsql down

pgdbnew:
	dbmate -d ./migrations/pgsql new $(name)

sqlc:
	sqlc generate