.PHONY: test lint fmt help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

test: ## Run tests
	go test -v -race -cover ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

---

### **3. Verify Structure**

Your directory should look like this:
```
thirdparty-sdk/
├── .gitignore
├── Makefile
├── README.md
├── go.mod
├── client/
├── config/
├── db/
│   └── queries/
├── errors/
├── examples/
├── internal/
│   ├── http/
│   └── retry/
├── logger/
├── models/
└── ui/
    ├── customer_media/
    ├── rebate/
    ├── shift/
    └── ticket_class/