# One file to rule to them all

ifndef PROJECT_NAME
PROJECT_NAME := omg
endif

ifndef PRODUCTION_ENVIRONMENT:
PRODUCTION_ENVIRONMENT := prod
endif

ifndef DOCKER_BIN:
DOCKER_BIN := docker
endif

ifndef DOCKER_COMPOSE_BIN:
DOCKER_COMPOSE_BIN := docker-compose
endif

build-local-go-image:
	${DOCKER_BIN} build -f build/local.go.Dockerfile -t ${PROJECT_NAME}-go-local:latest .
	-${${DOCKER_BIN} images -f "dangling=true" -q}

# ----------------------------
# Project level Methods
# ----------------------------
teardown:
	${COMPOSE} down -v
	${COMPOSE} rm --force --stop -v

setup: api-setup
boilerplate: api-boilerplate
run: api-run
test: api-test

# ----------------------------
# api Methods
# ----------------------------
API_COMPOSE = ${COMPOSE} run --rm --service-ports -w /api api
ifdef CONTAINER_SUFFIX
api-test: api-setup
endif
api-test:
	@${API_COMPOSE} sh -c "go test -mod=vendor -coverprofile=c.out -failfast -timeout 5m ./..."
api-run:
	@${API_COMPOSE} sh -c "go run -mod=vendor cmd/serverd/*.go"
api-build-binaries:
	@${API_COMPOSE} sh -c "\
		go clean -mod=vendor -i -x -cache ./... && \
		go build -mod=vendor -v -a -i -o binaries/serverd ./cmd/serverd && \
		go build -mod=vendor -v -a -i -o binaries/job ./cmd/job && \
		go build -mod=vendor -v -a -i -o binaries/multipartitionconsumer ./cmd/multipartitionconsumer && \
		go build -mod=vendor -v -a -i -o binaries/singlepartitionconsumer ./cmd/singlepartitionconsumer && \
		go build -mod=vendor -v -a -i -o binaries/datagen ./cmd/datagen "
api-update-vendor:
	@${API_COMPOSE} sh -c "go mod tidy -compat=1.17 && go mod vendor"

api-gen-mocks:
	${COMPOSE} run --name ${PROJECT_NAME}-mockery-$${CONTAINER_SUFFIX:-local} --rm -w /api --entrypoint '' mockery /bin/sh -c "\
		mockery --dir internal/authenticate --all --recursive --inpackage && \
		mockery --dir internal/controller --all --recursive --inpackage && \
		mockery --dir internal/repository --all --recursive --inpackage"
api-pg-migrate:
	${COMPOSE} run --rm pg-migrate sh -c './migrate -path /api-migrations -database $$PG_URL up'
api-pg-drop:
	${COMPOSE} run --rm pg-migrate sh -c './migrate -path /api-migrations -database $$PG_URL drop'
api-pg-redo: api-pg-drop api-pg-migrate
api-gen-models:
	@${API_COMPOSE} sh -c 'sqlboiler --wipe psql && GOFLAGS="-mod=vendor" goimports -w internal/repository/orm/*.go'

api-boilerplate: api-setup api-gen-models

ifdef CONTAINER_SUFFIX
api-setup: volumes pg sleep api-pg-migrate
else
api-setup: pg sleep api-pg-migrate
api-setup:
	${DOCKER_BIN} image inspect ${PROJECT_NAME}-go-local:latest >/dev/null 2>&1 || make build-local-go-image
endif

# ----------------------------
# Base Methods
# ----------------------------
volumes:
	${COMPOSE} up -d alpine
	${DOCKER_BIN} cp ${shell pwd}/api/. ${PROJECT_NAME}-alpine-$${CONTAINER_SUFFIX:-local}:/api
	${DOCKER_BIN} cp ${shell pwd}/api/data/migrations/. ${PROJECT_NAME}-alpine-$${CONTAINER_SUFFIX:-local}:/api-migrations
	${DOCKER_BIN} cp ${shell pwd}/web/. ${PROJECT_NAME}-alpine-$${CONTAINER_SUFFIX:-local}:/web

COMPOSE := PROJECT_NAME=${PROJECT_NAME} ${DOCKER_COMPOSE_BIN} -f build/docker-compose.base.yaml
ifdef CONTAINER_SUFFIX
COMPOSE := ${COMPOSE} -f build/docker-compose.ci.yaml -p ${CONTAINER_SUFFIX}
else
COMPOSE := ${COMPOSE} -f build/docker-compose.local.yaml
endif

pg:
	${COMPOSE} up -d pg

sleep:
	sleep 5
