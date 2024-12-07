.PHONY: build test

run: build
	./build/ufc_gcal

build:
	go build -o build/ufc_gcal ./cmd/ufc_gcal

test:
	go test ./...

clean:
	rm -rf build/
