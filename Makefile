run:
	go run cmd/ufc_gcal/main.go

build:
	go build -o build/ufc_gcal cmd/ufc_gcal/main.go

test:
	go test ./...

clean:
	rm -rf build/
