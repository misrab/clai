.PHONY: test build run clean

test:
	go test ./...

build:
	@mkdir -p bin
	go build -o bin/clai .

run: build
	./bin/clai $(ARGS)

clean:
	rm -f bin/clai
