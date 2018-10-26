PACKAGES=$(shell go list ./... | grep -v /vendor/)
RACE := $(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))
VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/')

help:
	@echo 'Available commands:'
	@echo
	@echo 'Usage:'
	@echo '    make build    		Compile the project.'
	@echo '    make test    		Run project tests.'
	@echo '    make deps     		Install go deps.'
	@echo '    make add_deps		Add dependencies with govendor.'
	@echo '    make clean    		Clean the project.'
	@echo

test:
	@go test ${RACE} ${PACKAGES}

deps:
	go get -u github.com/kardianos/govendor
	go get -u github.com/tkanos/gonfig
	go get -u github.com/go-chi/chi
	go get -u github.com/go-chi/chi/middleware
	go get -u github.com/go-chi/hostrouter

add_deps:
	govendor fetch github.com/tkanos/gonfig
	govendor fetch github.com/go-chi/chi
	govendor fetch github.com/go-chi/chi/middleware
	govendor fetch github.com/go-chi/hostrouter

build:
	@echo "Compiling Brypt Server"
	@mkdir -p ./bin
	@go build -o ./bin/bserv ./cmd/brypt-server
	@echo "Completed!"
