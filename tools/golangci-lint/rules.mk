golangci_lint_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
golangci_lint := $(golangci_lint_cwd)/bin/golangci-lint

$(golangci_lint): $(golangci_lint_cwd)/go.mod
	@echo building golangci-lint...
	@cd $(golangci_lint_cwd) && go build -o $@ github.com/golangci/golangci-lint/cmd/golangci-lint
	@cd $(golangci_lint_cwd) && go mod tidy

# go-lint: lint Go code with GolangCI-Lint
.PHONY: go-lint
go-lint: $(golangci_lint)
	@echo linting Go code with golangci-lint...
	@$(golangci_lint) run
