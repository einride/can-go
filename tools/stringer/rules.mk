stringer_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
stringer := $(stringer_cwd)/bin/stringer
PATH := $(PATH):$(dir $(stringer))

$(stringer): $(stringer_cwd)/go.mod
	@echo building stringer...
	@cd $(stringer_cwd) && go build -o $@ golang.org/x/tools/cmd/stringer
	@cd $(stringer_cwd) && go mod tidy
