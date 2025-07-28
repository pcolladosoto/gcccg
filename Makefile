CC := go

all: gcccg-darwin-arm64 gcccg-darwin-amd64 gcccg-linux-arm64 gcccg-linux-amd64

gcccg-darwin-arm64: *.go changelog.tmpl
	GOOS=darwin GOARCH=arm64 $(CC) build -o $@

gcccg-darwin-amd64: *.go changelog.tmpl
	GOOS=darwin GOARCH=amd64 $(CC) build -o $@

gcccg-linux-arm64: *.go changelog.tmpl
	GOOS=linux GOARCH=arm64 $(CC) build -o $@

gcccg-linux-amd64: *.go changelog.tmpl
	GOOS=linux GOARCH=amd64 $(CC) build -o $@

.PHONY: clean
clean:
	@rm -rf gcccg*
