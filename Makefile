.PHONY: help
help:

.PHONY: test
test:
	go test -v -race ./...

.PHONY: snapshot
snapshot:
	goreleaser --snapshot --skip=publish --clean

.PHONY: clean
clean:
	-rm -rf dist

# generate CREDITS file by using https://github.com/Songmu/gocredits
CREDITS: go.mod go.sum
	gocredits . > CREDITS
