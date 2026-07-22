BINARY_NAME := fp-estimator
GO_MODULE := github.com/senthilsweb/ai-dlc-fp-estimation

.PHONY: build build-linux build-darwin run clean docker help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build native binary
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo "==> Built $(BINARY_NAME) ($$(du -h $(BINARY_NAME) | cut -f1))"

build-linux: ## Cross-compile for Linux amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-linux-amd64 .

build-darwin: ## Cross-compile for macOS arm64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-arm64 .

run: build ## Build and run with the default dataset
	./$(BINARY_NAME)

docker: ## Build Docker image
	docker build -t $(BINARY_NAME) .

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux-amd64 $(BINARY_NAME)-darwin-arm64
	@echo "==> Cleaned"
