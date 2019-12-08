SERVICENAME="device-manager" 
PROJECTNAME="iota-3345"

GCPYAML="cmd/devicemanager/app.yaml"

VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
TIMESTAMP := $(shell date "+%Y%m%d-%H%M%S")
VERSIONGCP := $(shell echo $(VERSION) | tr '.' '-' | tr ':' '-')

LDFLAGS=-ldflags "-X=main.DeviceManagerVersion=$(VERSION) -X=main.DeviceManagerBuild=$(BUILD)"
MAKEFLAGS += --silent
PID := /tmp/.$(SERVICENAME).pid

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(GOPATH)/src/github.com/n7down/iot-core/cmd/devicemanager/*.go
GOFILESLIST=$(shell find . -name '*.go')

.PHONY: install
install:
	echo "installing... \c"
	@go get ./...
	echo "done"

.PHONY: build
build: clean
	echo "building... \c"
	CGO_ENABLED=0 GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(SERVICENAME) $(GOFILES)
	echo "done"
					
.PHONY: generate
generate:
	echo "generating dependency files... \c"
	@GOBIN=$(GOBIN) go generate ./...
	echo "done"

.PHONY: compile
compile: install build

.PHONY: start-server
start-server: stop-server
	echo "starting server... \c"
	@$(GOBIN)/$(SERVICENAME) 2>&1 & echo $$! > $(PID)
	echo "done"
	cat $(PID) | sed "/^/s/^/  \>  PID: /"

.PHONY: stop-server
stop-server:
	echo "stopping server... \c"
	@touch $(PID)
	@kill `cat $(PID)` 2> /dev/null || true
	@rm $(PID)
	echo "done"

.PHONY: start
start: compile start-server

.PHONY: stop
stop: stop-server

.PHONY: test
test:
	@GOCACHE=off go test -short ${GOFILESLIST}

.PHONY: vet
vet:
	@go vet ${GOFILESLIST}

.PHONY: clean
clean:
	echo "cleaning build cache... \c"
	@go clean
	@rm -rf bin/
	echo "done"

.PHONY: docker-build-bin
docker-build-bin: build-amd 
	echo "building $(SERVICENAME) docker container..."
	docker build -t gcr.io/"$(PROJECTNAMEDEV)"/"$(SERVICENAME)":"$(VERSION)" --label "version"="$(VERSION)" --label "build"="$(BUILD)" -f build/dockerfiles/app/Dockerfile.bin .
	echo "done"

.PHONY: deploy-dev-bin
deploy-dev-bin: docker-build-bin docker-push docker-deploy

.PHONY: deploy-dev-bin-promote
deploy-dev-bin-promote: docker-build-bin docker-push docker-deploy-promote

.PHONY: docker-build
docker-build:
	echo "building $(SERVICENAME) docker container..."
	docker build -t gcr.io/"$(PROJECTNAMEDEV)"/"$(SERVICENAME)":"$(VERSION)" --label "version"="$(VERSION)" --label "build"="$(BUILD)" -f build/dockerfiles/app/Dockerfile.dev .
	echo "done"

.PHONY: docker-push
docker-push:
	echo "pushing docker container..."
	docker push gcr.io/"$(PROJECTNAMEDEV)"/"$(SERVICENAME)":"$(VERSION)"
	echo "done"

.PHONY: docker-deploy
docker-deploy:
	echo "deploying version [$(VERSIONGCP)] to gcloud..."
	gcloud app deploy "$(DEVYAML)" -v "$(VERSIONGCP)" --project="$(PROJECTNAMEDEV)" --image-url gcr.io/"$(PROJECTNAMEDEV)"/"$(SERVICENAME)":"$(VERSION)" --no-promote
	echo "done"

.PHONY: deploy-dev
deploy-dev: docker-build docker-push docker-deploy

.PHONY: docker-deploy-promote
docker-deploy-promote:
	echo "deploying version [$(VERSIONGCP)] to gcloud..."
	gcloud app deploy "$(DEVYAML)" -v "$(VERSIONGCP)" --project="$(PROJECTNAMEDEV)" --image-url gcr.io/"$(PROJECTNAMEDEV)"/"$(SERVICENAME)":"$(VERSION)" --stop-previous-version -q
	echo "done"

.PHONY: deploy-dev-promote
deploy-dev-promote: docker-build docker-push docker-deploy-promote

.PHONY: deploy-dev-stanard-promote
deploy-dev-standard-promote:
	echo "deploying to gcloud..."
	gcloud app deploy "$(GCPYAML)" -v "$(VERSIONGCP)" --project="$(PROJECTNAME)" --stop-previous-version -q 
	echo "done"

.PHONY: docker-build-prod
docker-build-prod:
	echo "building $(SERVICENAME) docker container..."
	docker build -t gcr.io/"$(PROJECTNAMEPROD)"/"$(SERVICENAME)":"$(VERSION)" --label "version"="$(VERSION)" --label "build"="$(BUILD)" -f build/dockerfiles/app/Dockerfile.prod .
	echo "done"

.PHONY: docker-push-prod
docker-push-prod:
	echo "pushing docker container..."
	docker push gcr.io/"$(PROJECTNAMEPROD)"/"$(SERVICENAME)":"$(VERSION)"
	echo "done"

.PHONY: docker-deploy-prod
docker-deploy-prod:
	echo "deploying to gcloud..."
	gcloud app deploy "$(PRODYAML)" -v "$(VERSIONGCP)" --project="$(PROJECTNAMEPROD)" --image-url gcr.io/"$(PROJECTNAMEPROD)"/"$(SERVICENAME)":"$(VERSION)" --no-promote
	echo "done"

.PHONY: deploy-prod
deploy-prod: docker-build-prod docker-push-prod docker-deploy-prod

.PHONY: help
help:
	echo "Choose a command run in $(SERVICENAME):"
	echo " install - installs all dependencies for the project"
	echo " build - builds a binary"
	echo " compile - installs all dependencies and builds a binary"
	echo " clean - cleans the cache and cleans up the build files"
						
