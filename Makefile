GO ?= go
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go" -type f -not -path './vendor/*')
DOCKER_COMPOSE_FILES = -f docker-compose.yaml -f docker-compose.mongodb.yaml -f docker-compose.rabbitmq.yaml -f docker-compose.redis.yaml -f docker-compose.rating.yaml -f docker-compose.tester.yaml -f docker-compose.carrier.yaml


.PHONY: all
all: fmt lint vet test

.PHONY: build
build:
	$(GO) build -o bin/rating-agent-hep .

.PHONY: test
test:
	$(GO) test -cover -coverprofile=coverage.txt $(PACKAGES) && echo "\n==>\033[32m Ok\033[m\n" || exit 1

.PHONY: test-short
test-short:
	$(GO) test -cover -coverprofile=coverage.txt --short $(PACKAGES) && echo "\n==>\033[32m Ok\033[m\n" || exit 1

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: lint
lint:
	for pkg in ${PACKAGES}; do \
		golint -set_exit_status $$pkg || GOLINT_FAILED=1; \
	done; \
	[ -z "$$GOLINT_FAILED" ]

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: clean
clean:
	$(GO) clean -modcache -x -i ./...
	find . -name coverage.txt -delete
	rm bin/*

.PHONY: docker-build
docker-build:
	docker build . -t canyan/rating-agent-hep:master

.PHONY: docker-build-acceptance
docker-build-acceptance:
	docker build . -t canyan/rating-agent-hep:master -f Dockerfile.acceptance

.PHONY: docker-start
docker-start:
	docker-compose -f docker-compose.rabbitmq.yaml up -d rabbitmq
	for i in `seq 1 60`; do docker-compose -f docker-compose.rabbitmq.yaml exec -T rabbitmq rabbitmqctl await_online_nodes 1 && break || true; sleep 5; done
	sleep 5
	docker-compose $(DOCKER_COMPOSE_FILES) up -d

.PHONY: docker-pull
docker-pull:
	docker-compose $(DOCKER_COMPOSE_FILES) -f docker-compose.kamailio.yaml -f docker-compose.opensips.yaml pull

.PHONY: docker-start-kamailio
docker-start-kamailio: docker-start
	docker-compose -f docker-compose.kamailio.yaml up -d

.PHONY: docker-start-opensips
docker-start-opensips: docker-start
	docker-compose -f docker-compose.opensips.yaml up -d

.PHONY: docker-test-kamailio
docker-test-kamailio:
	docker exec rating-agent-hep_tester_1 pytest -k kamailio /tests/ -vv

.PHONY: docker-test-opensips
docker-test-opensips:
	docker exec rating-agent-hep_tester_1 pytest -k opensips /tests/ -vv

.PHONY: docker-test
docker-test: docker-test-kamailio docker-test-opensips

.PHONY: docker-logs
docker-logs:
	docker-compose $(DOCKER_COMPOSE_FILES) ps -a
	docker-compose $(DOCKER_COMPOSE_FILES) logs

.PHONY: docker-stop
docker-stop:
	docker-compose $(DOCKER_COMPOSE_FILES) -f docker-compose.kamailio.yaml -f docker-compose.opensips.yaml  down
