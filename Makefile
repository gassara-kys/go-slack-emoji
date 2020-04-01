APPNAME := $(shell basename `pwd`)
SRCS := $(shell ls *.go | grep -v '_test.go')
LDFLAGS := -ldflags="-s -w -extldflags \"-static\""
.PHONY: all
all: run

.PHONY: clean
clean:
	rm -rf bin/* && ls images/* | grep -v "_mv.sh" | xargs rm -f 

.PHONY: fmt
fmt: $(SRCS) 
	go fmt

.PHONY: tidy
tidy: fmt
	go mod tidy

.PHONY: test
test: tidy
	go test -v -cover ./...

.PHONY: build
build: test
	go build $(LDFLAGS) -o bin/$(APPNAME) .
	go install

.PHONY: run
run: build
	source env.sh && bin/$(APPNAME)

.PHONY: download
download: build
	source env.sh && bin/$(APPNAME) download -f
