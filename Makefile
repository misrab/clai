.PHONY: test build run clean install

test:
	go test ./...

build:
	@mkdir -p bin
	go build -o bin/clai .

run: build
	./bin/clai $(ARGS)

install:
	go install .

clean:
	rm -f bin/clai
