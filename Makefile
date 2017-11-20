
build:
	go build

publish:
	go get github.com/goreleaser/goreleaser
	./scripts/goreleaser.yaml.sh "$(LD_RELEASE_FLAGS)" >/tmp/gorel.yaml
	goreleaser --config /tmp/gorel.yaml

