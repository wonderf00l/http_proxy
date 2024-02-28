PROJECT_BIN = $(CURDIR)/bin

build:
	go build -o $(PROJECT_BIN)/app cmd/*.go

.PHONY:
run: build
	bin/app