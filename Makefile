all: \
	go-stringer \
	go-mock-gen \
	testdata \
	go-lint \
	go-test \
	go-mod-tidy \
	git-verify-nodiff

include build/rules.mk

.PHONY: clean
clean:
	rm -rf $(FILES_DIR)

.PHONY: go-lint
go-lint: $(GOLANGCI_LINT)
	# funlen: too strict
	# dupl: allow duplication in tests
	# interfacer: deprecated
	# godox: allow TODOs
	# lll: long go:generate directives
	$(GOLANGCI_LINT) run --enable-all --disable funlen,dupl,interfacer,godox,lll

.PHONY: go-mock-gen
go-mock-gen: \
	internal/mocks/mockcanrunner/mocks.go \
	internal/mocks/mockclock/mocks.go \
	internal/mocks/mocksocketcan/mocks.go

internal/mocks/mockcanrunner/mocks.go: pkg/canrunner/run.go $(GOBIN)
	$(GOBIN) -m -run github.com/golang/mock/mockgen -destination $@ \
		-package mockcanrunner go.einride.tech/can/pkg/canrunner \
		Node,TransmittedMessage,ReceivedMessage,FrameTransmitter,FrameReceiver

internal/mocks/mockclock/mocks.go: internal/clock/clock.go $(GOBIN)
	$(GOBIN) -m -run github.com/golang/mock/mockgen -destination $@ \
		-package mockclock go.einride.tech/can/internal/clock \
		Clock,Ticker

internal/mocks/mocksocketcan/mocks.go: pkg/socketcan/fileconn.go $(GOBIN)
	$(GOBIN) -m -run github.com/golang/mock/mockgen -destination $@ \
		-package mocksocketcan -source $<

.PHONY: go-stringer
go-stringer: \
	pkg/descriptor/sendtype_string.go \
	pkg/socketcan/errorclass_string.go \
	pkg/socketcan/protocolviolationerrorlocation_string.go \
	pkg/socketcan/protocolviolationerror_string.go \
	pkg/socketcan/controllererror_string.go \
	pkg/socketcan/transceivererror_string.go

%_string.go: %.go
	go generate $<

.PHONY: testdata
testdata:
	go run cmd/cantool/main.go generate testdata/dbc testdata/gen/go

.PHONY: go-test
go-test:
	go test -race -cover ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy -v
