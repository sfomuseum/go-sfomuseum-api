GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/api cmd/api/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/test-methods cmd/test-methods/main.go
