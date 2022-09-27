include Makefile.common

.PHONY: default manager worker


default: manager worker

install:
	$(GOCMD) install

manager:
	$(GOBUILD) -o bin/manager cmd/manager/main.go

worker:
	$(GOBUILD) -o bin/worker cmd/worker/main.go


build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/manager cmd/manager/main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/worker cmd/worker/main.go

clean:
	rm -rf ./bin/*