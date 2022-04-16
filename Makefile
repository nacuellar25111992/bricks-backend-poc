# Makefile for releasing bricks-backend-poc
#
# The release version is controlled from internal/version

TAG?=latest
NAME:=bricks-backend-poc
DOCKER_REPOSITORY:=nacuellar25111992
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
DOCKER_IMAGE_PLATFORM:=linux/amd64
GIT_COMMIT:=$(shell git describe --dirty --always)
VERSION:=$(shell grep 'VERSION' internal/version/version.go | awk '{ print $$4 }' | tr -d '"')
EXTRA_RUN_ARGS?=

run:
	go run -ldflags "-s -w -X github.com/nacuellar25111992/bricks-backend-poc/internal/version.REVISION=$(GIT_COMMIT)" cmd/bricks-backend-poc/* \
	--log-level=debug \
	--backend-url=https://httpbin.org/status/401 \
	--backend-url=https://httpbin.org/status/500

test:
	go test ./... -coverprofile cover.out

build:
	GIT_COMMIT=$$(git rev-list -1 HEAD) && CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/nacuellar25111992/bricks-backend-poc/internal/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/bricks-backend-poc ./cmd/bricks-backend-poc/*

# TODO: golangci-lint
fmt:
	gofmt -l -s -w .
	goimports -l -w .

build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) . --platform=$(DOCKER_IMAGE_PLATFORM)

test-container: run-container
	@docker ps
	@TOKEN=$$(curl -sd 'test' localhost:9898/token | jq -r .token) && \
	curl -sH "Authorization: Bearer $${TOKEN}" localhost:9898/token/validate | grep test

run-container: remove-container build-container
	docker run -dp 9898:9898 --name=bricks-backend-poc $(DOCKER_IMAGE_NAME):$(VERSION)

remove-container:
	docker rm -f bricks-backend-poc || true

push-container: remove-container build-container
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_IMAGE_NAME):latest

scan-container:
	docker scan $(DOCKER_IMAGE_NAME):$(VERSION)

inspect-container:
	docker inspect $(DOCKER_IMAGE_NAME):$(VERSION)

release:
	git tag $(VERSION)
	git push origin $(VERSION)

swagger:
	go install github.com/swaggo/swag/cmd/swag
	cd internal/api && $$(go env GOPATH)/bin/swag init -g server.go
