CC := go

gcccg: *.go changelog.tmpl
	$(CC) build -o $@
