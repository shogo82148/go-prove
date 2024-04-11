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
.PHONY: credits
credits: go.mod go.sum
	go mod download
	gocredits . > CREDITS
