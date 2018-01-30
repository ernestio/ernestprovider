test:
	cd ./providers/ && go test ./...

lint:
	gometalinter --config .linter.conf

deps:
	go get github.com/satori/uuid

dev-deps: deps
	go get github.com/ernestio/ernest-config-client
	go get github.com/nats-io/nats
	go get github.com/gucumber/gucumber/cmd/gucumber
	go get github.com/alecthomas/gometalinter
	go get github.com/r3labs/terraform
	gometalinter --install
