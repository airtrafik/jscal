#############################
# Global vars
#############################
PROJECT_NAME := jscal
PROJECT_MODULE := github.com/airtrafik/jscal
GO := go
GOLINTER := golangci-lint

# Directories
BIN_DIR := ./bin
COVERAGE_DIR := ./coverage
CMD_DIR := ./cmd/jscal

# Build output
BINARY_NAME := jscal
BINARY_PATH := $(BIN_DIR)/$(BINARY_NAME)

#############################
# Main targets
#############################
.DEFAULT_GOAL := all

all: clean lint test cover build

ci-test: test cover

#############################
# Build targets
#############################
build: $(BINARY_PATH)

$(BINARY_PATH): clean
	@echo "=== $(PROJECT_NAME) === [ build ]: Building binary..."
	@mkdir -p $(BIN_DIR)
	@cd $(CMD_DIR) && $(GO) build -o ../../$(BINARY_PATH) .
	@echo "=== $(PROJECT_NAME) === [ build ]: Binary created at $(BINARY_PATH)"

#############################
# Clean targets
#############################
clean:
	@echo "=== $(PROJECT_NAME) === [ clean ]: Removing build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(COVERAGE_DIR)
	@rm -f coverage.out coverage.html
	@find . -type f -name '*.test' -delete
	@find . -type f -name '*.out' -delete
	@echo "=== $(PROJECT_NAME) === [ clean ]: Clean complete"

#############################
# Test targets
#############################
test:
	@echo "=== $(PROJECT_NAME) === [ test ]: Running tests..."
	@mkdir -p $(COVERAGE_DIR)
	@echo "=== $(PROJECT_NAME) === [ test ]: Testing main module..."
	@$(GO) test -v -covermode=atomic -coverprofile=$(COVERAGE_DIR)/main.out ./...
	@for dir in convert/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			module=$$(basename "$$dir"); \
			echo "=== $(PROJECT_NAME) === [ test ]: Testing $$module converter module..."; \
			cd "$$dir" && $(GO) test -v -covermode=atomic -coverprofile=../../$(COVERAGE_DIR)/$$module.out ./... && cd - >/dev/null; \
		fi \
	done
	@echo "=== $(PROJECT_NAME) === [ test ]: Combining coverage reports..."
	@echo 'mode: atomic' > $(COVERAGE_DIR)/coverage.out
	@for file in $(COVERAGE_DIR)/*.out; do \
		if [ -f "$$file" ] && [ "$$file" != "$(COVERAGE_DIR)/coverage.out" ]; then \
			tail -n +2 "$$file" >> $(COVERAGE_DIR)/coverage.out 2>/dev/null || true; \
		fi \
	done
	@echo "=== $(PROJECT_NAME) === [ test ]: Tests complete"

test-verbose:
	@echo "=== $(PROJECT_NAME) === [ test-verbose ]: Running tests with verbose output..."
	@mkdir -p $(COVERAGE_DIR)
	@echo "=== $(PROJECT_NAME) === [ test-verbose ]: Testing main module..."
	@$(GO) test -v -count=1 -covermode=atomic -coverprofile=$(COVERAGE_DIR)/main.out ./...
	@for dir in convert/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			module=$$(basename "$$dir"); \
			echo "=== $(PROJECT_NAME) === [ test-verbose ]: Testing $$module converter module..."; \
			cd "$$dir" && $(GO) test -v -count=1 -covermode=atomic -coverprofile=../../$(COVERAGE_DIR)/$$module.out ./... && cd - >/dev/null; \
		fi \
	done
	@echo "=== $(PROJECT_NAME) === [ test-verbose ]: Combining coverage reports..."
	@echo 'mode: atomic' > $(COVERAGE_DIR)/coverage.out
	@for file in $(COVERAGE_DIR)/*.out; do \
		if [ -f "$$file" ] && [ "$$file" != "$(COVERAGE_DIR)/coverage.out" ]; then \
			tail -n +2 "$$file" >> $(COVERAGE_DIR)/coverage.out 2>/dev/null || true; \
		fi \
	done

#############################
# Coverage targets
#############################
cover: test
	@echo "=== $(PROJECT_NAME) === [ cover ]: Generating coverage report..."
	@$(GO) tool cover -html=$(COVERAGE_DIR)/main.out -o $(COVERAGE_DIR)/main.html 2>/dev/null || true
	@cd convert/ical && $(GO) tool cover -html=../../$(COVERAGE_DIR)/ical.out -o ../../$(COVERAGE_DIR)/ical.html 2>/dev/null || true
	@echo "=== $(PROJECT_NAME) === [ cover ]: Coverage reports generated:"
	@echo "  - Main module: $(COVERAGE_DIR)/main.html"
	@echo "  - iCal converter: $(COVERAGE_DIR)/ical.html"
	@echo "=== $(PROJECT_NAME) === [ cover ]: Coverage summary:"
	@printf "  Main module: "
	@$(GO) tool cover -func=$(COVERAGE_DIR)/main.out 2>/dev/null | tail -1 || echo "N/A"
	@printf "  iCal converter: "
	@cd convert/ical && $(GO) tool cover -func=../../$(COVERAGE_DIR)/ical.out 2>/dev/null | tail -1 || echo "N/A"

cover-view: cover
	@echo "=== $(PROJECT_NAME) === [ cover-view ]: Opening coverage report..."
	@$(GO) tool cover -html=$(COVERAGE_DIR)/main.out

#############################
# Lint targets
#############################
lint:
	@echo "=== $(PROJECT_NAME) === [ lint ]: Running linters..."
	@if command -v $(GOLINTER) >/dev/null 2>&1; then \
		$(GOLINTER) run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Falling back to go vet..."; \
		$(GO) vet ./...; \
	fi
	@echo "=== $(PROJECT_NAME) === [ lint ]: Linting complete"

lint-fix:
	@echo "=== $(PROJECT_NAME) === [ lint-fix ]: Running linters with auto-fix..."
	@if command -v $(GOLINTER) >/dev/null 2>&1; then \
		$(GOLINTER) run --fix ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

#############################
# Development helpers
#############################
fmt:
	@echo "=== $(PROJECT_NAME) === [ fmt ]: Formatting code..."
	@$(GO) fmt ./...
	@echo "=== $(PROJECT_NAME) === [ fmt ]: Formatting complete"

vet:
	@echo "=== $(PROJECT_NAME) === [ vet ]: Running go vet..."
	@$(GO) vet ./...
	@echo "=== $(PROJECT_NAME) === [ vet ]: Vet complete"

mod-tidy:
	@echo "=== $(PROJECT_NAME) === [ mod-tidy ]: Tidying modules..."
	@$(GO) mod tidy
	@for dir in convert/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "=== $(PROJECT_NAME) === [ mod-tidy ]: Tidying $$dir..."; \
			cd "$$dir" && $(GO) mod tidy && cd - >/dev/null; \
		fi \
	done
	@echo "=== $(PROJECT_NAME) === [ mod-tidy ]: Modules tidied"

mod-verify:
	@echo "=== $(PROJECT_NAME) === [ mod-verify ]: Verifying modules..."
	@$(GO) mod verify
	@for dir in convert/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "=== $(PROJECT_NAME) === [ mod-verify ]: Verifying $$dir..."; \
			cd "$$dir" && $(GO) mod verify && cd - >/dev/null; \
		fi \
	done
	@echo "=== $(PROJECT_NAME) === [ mod-verify ]: Modules verified"

#############################
# Install target
#############################
install: build
	@echo "=== $(PROJECT_NAME) === [ install ]: Installing binary..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME) || cp $(BINARY_PATH) ~/go/bin/$(BINARY_NAME)
	@echo "=== $(PROJECT_NAME) === [ install ]: Binary installed to GOPATH/bin"

#############################
# Help target
#############################
help:
	@echo "Available targets:"
	@echo "  all          - Run clean, lint, test, cover, and build"
	@echo "  build        - Build the jscal binary to ./bin/"
	@echo "  clean        - Remove build artifacts and temporary files"
	@echo "  test         - Run all tests"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  cover        - Generate test coverage report"
	@echo "  cover-view   - Generate and open coverage report in browser"
	@echo "  lint         - Run golangci-lint or go vet"
	@echo "  lint-fix     - Run golangci-lint with auto-fix"
	@echo "  fmt          - Format code with go fmt"
	@echo "  vet          - Run go vet"
	@echo "  mod-tidy     - Tidy go modules"
	@echo "  mod-verify   - Verify go modules"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  help         - Show this help message"

.PHONY: all build clean test test-verbose cover cover-view lint lint-fix fmt vet mod-tidy mod-verify install help
