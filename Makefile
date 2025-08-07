.PHONY: all clean fmt fmt-check lint changelog test-summary help
.DEFAULT_GOAL := all

# Go files for dependency tracking
GO_FILES := $(shell find . -name '*.go' -not -path './bin/*')

# Default target
all: test lint bin/ynab-discover

# Test results depend on Go source files
test-results.json coverage.out: $(GO_FILES)
	go test -race -coverprofile=coverage.out -covermode=atomic -json -v ./... > test-results.json
	grep -v "cmd/ynab-discover/main.go" coverage.out > coverage.filtered || true
	mv coverage.filtered coverage.out

# Test target just ensures test results exist
test: test-results.json

# Coverage HTML depends on coverage.out
coverage.html: coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Test coverage target
test-coverage: coverage.html

# Binary depends on Go source files
bin/ynab-discover: $(GO_FILES)
	@mkdir -p bin
	go build -o bin/ynab-discover ./cmd/ynab-discover

# Build alias
build: bin/ynab-discover

# Run linter
lint:
	docker run --rm -v "$(PWD):/app" -w /app golangci/golangci-lint:v2.3.1 golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html test-results.json

# Format code
fmt:
	go fmt ./...

# Check if code is formatted
fmt-check:
	test -z "$$(gofmt -l .)"

# Run all checks (what CI runs)
ci: fmt-check test-results.json lint

# Generate changelog since last tag (or all commits if no previous tag)
changelog:
	@PREV_TAG=$$(git describe --tags --abbrev=0 HEAD~1 2>/dev/null || echo ""); \
	if [ -n "$$PREV_TAG" ]; then \
		echo "Changes since $$PREV_TAG:"; \
		echo ""; \
		git log --pretty=format:"- %s" $$PREV_TAG..HEAD; \
	else \
		echo "Changes:"; \
		echo ""; \
		git log --pretty=format:"- %s"; \
	fi

# Generate test summary with results and coverage
test-summary: test-results.json coverage.out
	@echo "## 🧪 Test Results"
	@echo ""
	@PASSED_TESTS=$$(grep '"Action":"pass"' test-results.json | wc -l); \
	FAILED_TESTS=$$(grep '"Action":"fail"' test-results.json | wc -l); \
	TOTAL_TESTS=$$((PASSED_TESTS + FAILED_TESTS)); \
	echo "**$$PASSED_TESTS/$$TOTAL_TESTS tests passed**"
	@echo ""
	@if grep -q '"Action":"fail"' test-results.json; then \
		echo "❌ **Some tests failed**"; \
		echo ""; \
		echo "### Failed Tests"; \
		grep '"Action":"fail"' test-results.json | jq -r '"- \`" + .Test + "\`"' 2>/dev/null || echo "- Check logs for details"; \
	else \
		echo "✅ **All tests passed!**"; \
	fi
	@echo ""
	@echo "## 🎯 Test Coverage"
	@echo ""
	@echo "| File | Coverage |"
	@echo "|------|----------|"
	@go tool cover -func=coverage.out | head -n -1 | awk '{print "| " $$1 " | " $$3 " |"}'
	@TOTAL=$$(go tool cover -func=coverage.out | tail -1 | awk '{print $$3}'); \
	echo "| **Total** | **$$TOTAL** |"

# Show help
help:
	@echo "Available targets:"
	@echo "  all              - Run test, lint, and build"
	@echo "  test             - Run tests (generates test-results.json)"
	@echo "  test-coverage    - Generate coverage report (coverage.html)"
	@echo "  lint             - Run linter"
	@echo "  build            - Build binary (bin/ynab-discover)"
	@echo "  clean            - Remove all generated files"
	@echo "  fmt              - Format code"
	@echo "  fmt-check        - Check if code is formatted"
	@echo "  ci               - Run all CI checks"
	@echo "  test-summary     - Show test summary (requires test results)"
	@echo "  changelog        - Generate changelog since last tag"
	@echo "  help             - Show this help"