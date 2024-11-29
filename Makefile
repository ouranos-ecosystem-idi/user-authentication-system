GOCMD			=	go
GOFMT			=	$(GOCMD) fmt
GOVET			=	$(GOCMD) vet
GOBUILD		=	$(GOCMD) build
GORETURNS	= goreturns
LINTER		=	golangci-lint
APP				=	authenticator-backend
ENV				= local
PORT			= 8081
MOCK_SRC_REPOSITORY = $(wildcard domain/repository/*.go)
MOCK_SRC_USECASE = $(wildcard usecase/*usecase.go)
MOCK_SRC_HANDLER = $(wildcard presentation/http/echo/handler/*.go)
MOCK_FILES = $(wildcard test/mock/*.go)

.PHONY: test

all:
	make goreturns
	make genmock
	make vet
	make lint-fix
	make test
	make swaggo
	make build-go
	make build-with-docker
	make scan-image

validate:
	$(GOFMT) -w -s ./...

goreturns:
	$(GORETURNS) -w .

lint:
	$(LINTER) run

lint-fix:
	$(LINTER) run --fix

vet:
	$(GOVET) ./...

swaggo:
	swag init

build-with-docker:
	docker build -t $(APP) .

build-go:
	$(GOBUILD) main.go

genmock: $(MOCK_SRC_REPOSITORY)
	rm $(MOCK_FILES)
	go generate $(MOCK_SRC_REPOSITORY)
	go generate $(MOCK_SRC_USECASE)
	go generate $(MOCK_SRC_HANDLER)

test:
	go test -v -cover -covermode=atomic ./...

test-coverage:
	go test -v -cover -coverprofile=cover.out -covermode=atomic ./presentation/http/echo/handler/... ./usecase/... ./infrastructure/persistence/datastore/... ./infrastructure/firebase/...
	go tool cover -html=cover.out -o cover.html

run:
	docker run -v $(PWD)/config/:/app/config/ -td -i --network docker.internal --env-file config/$(ENV).env -p $(PORT):$(PORT) --name $(APP) $(APP)

scan-image:
	docker run -v /var/run/docker.sock:/var/run/docker.sock --rm aquasec/trivy image --severity HIGH,CRITICAL $(APP)

clean:
	docker stop $(APP); docker rm $(APP)
	docker container prune --force

api-scan:
	./scripts/api-scan.sh

idp-add-local:
	go run cmd/add_local_user/main.go
