SHELL := /bin/bash
GO := go
NAME := cmg
OS := $(shell uname)
MAIN_GO := main.go
ROOT_PACKAGE := $(GIT_PROVIDER)/$(ORG)/$(NAME)
GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
PKGS := $(shell go list ./... | grep -v /vendor | grep -v generated)
PKGS := $(subst  :,_,$(PKGS))
BUILDFLAGS := ''
CGO_ENABLED = 0
VENDOR_DIR=vendor

all: build

check: fmt build test

build:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

test: 
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test $(PACKAGE_DIRS) -test.v -coverprofile cp.out

coverage:
    gopherbadger-md="README.md"

full: $(PKGS)

install:
	GOBIN=${GOPATH}/bin $(GO) install -ldflags $(BUILDFLAGS) $(MAIN_GO)

fmt:
	@gofmt -s -w -l **/*.go

clean:
	rm -rf build release

linux:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

multiarch:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=${ARCH} $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(ARCH)/$(NAME) $(MAIN_GO)

.PHONY: release clean

FGT := $(GOPATH)/bin/fgt
$(FGT):
	go get github.com/GeertJohan/fgt

GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	go get github.com/golang/lint/golint

$(PKGS): $(GOLINT) $(FGT)
	@echo "LINTING"
	@$(FGT) $(GOLINT) $(GOPATH)/src/$@/*.go
	@echo "VETTING"
	@go vet -v $@
	@echo "TESTING"
	@go test -v $@

.PHONY: lint
lint: vendor | $(PKGS) $(GOLINT) # ❷
	@cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
	    test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret

watch:
	reflex -r "\.go$" -R "vendor.*" make skaffold-run

skaffold-run: build
	skaffold run -p dev

dbuild: fmt
	DOCKER_BUILDKIT=1 docker build --tag caladreas/cmg:latest .

dclean:
	docker rm cmg

drun:
	docker run --name cmg -p 8080:8080 caladreas/cmg:latest serve

dpush: dbuild
	docker push caladreas/cmg:latest

gbuild: fmt
	gcloud builds submit --tag gcr.io/$(PROJECT_ID)/cmg

gpush: dbuild
	docker tag caladreas/cmg:latest gcr.io/$(PROJECT_ID)/cmg:latest
	docker push gcr.io/$(PROJECT_ID)/cmg:latest

gdeploy: gpush
	gcloud run deploy cmg --image=gcr.io/$(PROJECT_ID)/cmg:latest --memory=128Mi --max-instances=2 --timeout=30 --project=$(PROJECT_ID) --platform managed --allow-unauthenticated --region=europe-west4