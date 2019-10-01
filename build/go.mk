GOLANGCI_LINT_VERSION := 1.19.1
GOLANGCI_LINT := $(FILES_DIR)/golangci-lint/$(GOLANGCI_LINT_VERSION)/golangci-lint

GOBIN_VERSION := 0.0.13
GOBIN := $(FILES_DIR)/gobin/$(GOBIN_VERSION)/gobin
export PATH := $(dir $(GOBIN)):$(PATH)

$(GOLANGCI_LINT):
	mkdir -p $(dir $@)
	curl -s -L -o $(dir $(GOLANGCI_LINT))/archive.tar.gz \
		https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-$(UNAME)-amd64.tar.gz
	tar xzf $(dir $(GOLANGCI_LINT))/archive.tar.gz -C $(dir $(GOLANGCI_LINT)) --strip 1
	chmod +x $@
	touch $@

$(GOBIN):
	mkdir -p $(dir $@)
	curl -s -L -o $@ \
		https://github.com/myitcv/gobin/releases/download/v$(GOBIN_VERSION)/$(UNAME)-amd64
	chmod +x $@
	touch $@
