CC := go

gcccg: *.go changelog.tmpl
	$(CC) build -o $@

.PHONY: clean
clean:
	@rm -rf gcccg
