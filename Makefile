.PHONY: build run clean

build:
	@echo "Building EVM..."
	@mkdir -p bin
	go build -o bin/main .

run: build
	@echo "Running EVM..."
	./bin/main

clean:
	@echo "Cleaning up..."
	@rm -f bin/main