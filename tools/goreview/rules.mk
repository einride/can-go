goreview_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
goreview := $(goreview_cwd)/bin/goreview

$(goreview): $(goreview_cwd)/go.mod
	@echo building goreview...
	@cd $(goreview_cwd) && go build -o $@ github.com/einride/goreview/cmd/goreview
	@cd $(goreview_cwd) && go mod tidy

.PHONY: go-review
go-review: $(goreview)
	$(info [$@] reviewing Go code for Einride-specific conventions...)
	@$(goreview) -c 1 ./...
