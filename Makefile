include Makefile.vars

.SILENT:
.DEFAULT_GOAL := help

.PHONY: help
help:
	$(info Available Commands:)
	$(info -> install                 installs project dependencies)
	$(info -> test                    run all tests)
	$(info -> test-unit               run unit tests)
	$(info -> lint                    check coding style)
	$(info -> run                     run app)
	$(info -> db-start                starts development database)
	$(info -> db-stop                 stops development database)

.PHONY: install
install:
	go get github.com/google/wire/cmd/wire
	go mod tidy -v

.PHONY: wire
wire:
	wire $(SERVERDIR)

.PHONY: test
test:
	go test ./... -v

.PHONY: test-unit
test-unit:
	go test ./... -v -short

.PHONY: lint
lint:
	go get -u golang.org/x/lint/golint
	golint ./...

.PHONY: run
run: wire
	go run $(SERVERDIR)

.PHONY: db-start
db-start:
	docker run -d \
	--rm \
	-p 3306:3306 \
	--name mysql_maskada \
	-e MYSQL_ROOT_PASSWORD=$(DATABASE_ROOT_PASSWORD) \
	-e MYSQL_PASSWORD=$(DATABASE_PASSWORD) \
	-e MYSQL_USER=$(DATABASE_USERNAME) \
	-e MYSQL_DATABASE=$(DATABASE_NAME) \
	-v $(CURDIR)/details/db/migrations:/docker-entrypoint-initdb.d/ \
	mysql:8

.PHONY: db-stop
db-stop:
	docker stop -t 2 mysql_maskada

# ignore unknown commands
%:
    @:
