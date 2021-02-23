.PHONY: help
help:

.PHONY: test
test:
	go test -v -race ./...

.PHONY: snapshot
snapshot:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: clean
clean:
	-rm -rf dist
