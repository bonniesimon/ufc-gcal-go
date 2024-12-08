.PHONY: build test

run: build ## Run the application after building it
	./build/ufc_gcal

build: ## Build the binary
	go build -o build/ufc_gcal ./cmd/ufc_gcal

test: ## Run tests. No tests as of now.
	go test ./...

clean: ## Delete the build directory
	rm -rf build/

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@echo "Help command"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
