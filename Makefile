#TODO DEPENDENCIES make, angular cli, npm, docker, go > 1.17
PROJECTNAME = "pogong"

DB_INITIAL_PATH = "./pg/backup.sql"

.DEFAULT_GOAL := all
IMAGES = "postgres:15.2-alpine" "kjconroy/sqlc" 
GOMODS = "github.com/gzuidhof/tygo@latest" "github.com/gobuffalo/nulls" "github.com/stretchr/testify/require" "github.com/lib/pq"
# TODO run go mod tidy after all downloads
DB_INT_PORT = 5432
CONTAINER_DB = $(PROJECTNAME)db
RUNNING_DB := $(shell docker ps -q -a -f name=$(CONTAINER_DB))

all: init

install: npm && pull && go

go:
	for gomod in $(GOMODS); do \
		go get $$gomod; \
	done

pull:
	for images in $(IMAGES); do \
		docker pull $$images; \
	done

npm: 
	cd ui && npm ci && cd ..

init: $(if $(RUNNING_DB),startdb, initdb)

initdb:
	docker run -v "$(CURDIR)/pg"/:/pg/ --name=$(CONTAINER_DB) -p 5432:$(DB_INT_PORT) -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -e PGDATA="./pg/data" -d postgres:15.2-alpine
 
startdb:
	@if [ -z $$(docker ps -q -f name=$(CONTAINER_DB) -f status=running) ]; then \
        docker start $(CONTAINER_DB); \
    fi

killdb:
	docker kill $(CONTAINER_DB) && docker rm $(CONTAINER_DB)

backupdb:
	docker exec $(CONTAINER_DB) pg_dump -F p --if-exists --clean --create --no-owner --username=root --host=localhost -p $(DB_INT_PORT) pogong > ./pg/backup.sql

restoredb:
	@docker exec $(CONTAINER_DB) psql --quiet -f $(DB_INITIAL_PATH)

psql: init
	docker exec -it $(CONTAINER_DB) psql -U root -d pogong

generate: schemadb schemapy sqlc tygo

schemadb:
	docker exec $(CONTAINER_DB) pg_dump --username=root -F p --encoding=LATIN1 --no-blobs --schema-only --no-privileges --no-comments --no-tablespaces --no-unlogged-table-data --no-owner pogong > ./sql/schema.sql

schemapy:
	python makeschema.py

sqlc:
	docker run --rm -v $(CURDIR):/src -w /src kjconroy/sqlc:1.16.0 generate

tygo:
	tygo generate

test:
	go test -v -cover ./...

.PHONY: sqlc tygo npm pull initdb install go init generate backupdb killdb restoredb schemapy psql test