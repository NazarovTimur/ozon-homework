ifeq ($(OS),Windows_NT)
    # Windows
    RM = del /f /q
    RMDIR = rmdir /s /q
    MKDIR = if not exist bin mkdir bin
    TOOL_SUFFIX = .exe
    PATH_SEP = \\
    CHECK_EXIST = if not exist
else
    # Unix-like (Linux/macOS)
    RM = rm -f
    RMDIR = rm -rf
    MKDIR = mkdir -p bin
    TOOL_SUFFIX =
    PATH_SEP = /
    CHECK_EXIST = test ! -f
endif

TOOLS_DIR = $(shell go env GOPATH)$(PATH_SEP)bin
GOCYCLO = $(TOOLS_DIR)$(PATH_SEP)gocyclo$(TOOL_SUFFIX)
GOCOGNIT = $(TOOLS_DIR)$(PATH_SEP)gocognit$(TOOL_SUFFIX)

.PHONY: run-all test-coverage cyclo cognit install-tools clean report

run-all:
	go run ./cmd/server/main.go

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

install-tools:
ifeq ($(OS),Windows_NT)
	@if not exist "$(GOCYCLO)" ( \
		echo Installing gocyclo... && \
		go install github.com/fzipp/gocyclo/cmd/gocyclo@latest \
	)
	@if not exist "$(GOCOGNIT)" ( \
		echo Installing gocognit... && \
		go install github.com/uudashr/gocognit/cmd/gocognit@latest \
	)
else
	@if [ ! -f "$(GOCYCLO)" ]; then \
		echo Installing gocyclo... && \
		go install github.com/fzipp/gocyclo/cmd/gocyclo@latest; \
	fi
	@if [ ! -f "$(GOCOGNIT)" ]; then \
		echo Installing gocognit... && \
		go install github.com/uudashr/gocognit/cmd/gocognit@latest; \
	fi
endif

cyclo: install-tools
	@$(MKDIR)
	-@"$(GOCYCLO)" -over 10 . > bin/cyclo.txt
	@echo "Cyclo report saved to bin/cyclo.txt"

cognit: install-tools
	@$(MKDIR)
	-@"$(GOCOGNIT)" -over 10 . > bin/cognit.txt
	@echo "Cognit report saved to bin/cognit.txt"

report: cyclo cognit
	@echo "Reports generated in bin/"

clean:
ifeq ($(OS),Windows_NT)
	@if exist coverage.out $(RM) coverage.out
	@if exist coverage.html $(RM) coverage.html
	@if exist bin $(RMDIR) bin
else
	@if [ -f coverage.out ]; then $(RM) coverage.out; fi
	@if [ -f coverage.html ]; then $(RM) coverage.html; fi
	@if [ -d bin ]; then $(RMDIR) bin; fi
endif