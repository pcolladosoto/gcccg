CC := go

gcccg: cmd/main.go cmd/changelog.tmpl
	$(CC) build -o $@ ./cmd
