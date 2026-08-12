default: build

test:
	go test ./...

build:
	go build -o tflint-ruleset-cyrus

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-cyrus ~/.tflint.d/plugins
