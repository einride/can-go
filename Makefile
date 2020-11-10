SHELL := /bin/bash

all: \
	stringer-generate \
	mockgen-generate \
	testdata \
	go-lint \
	go-review \
	go-test \
	go-mod-tidy \
	git-verify-nodiff

include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/gomock/rules.mk
include tools/goreview/rules.mk
include tools/stringer/rules.mk

.PHONY: clean
clean:
	rm -rf tools/*/*/

.PHONY: mockgen-generate
mockgen-generate: \
	internal/gen/mock/mockcanrunner/mocks.go \
	internal/gen/mock/mockclock/mocks.go \
	internal/gen/mock/mocksocketcan/mocks.go

internal/gen/mock/mockcanrunner/mocks.go: pkg/canrunner/run.go $(mockgen)
	$(mockgen) -destination $@ -package mockcanrunner go.einride.tech/can/pkg/canrunner \
		Node,TransmittedMessage,ReceivedMessage,FrameTransmitter,FrameReceiver

internal/gen/mock/mockclock/mocks.go: internal/clock/clock.go $(mockgen)
	$(mockgen) -destination $@ -package mockclock go.einride.tech/can/internal/clock \
		Clock,Ticker

internal/gen/mock/mocksocketcan/mocks.go: pkg/socketcan/fileconn.go $(mockgen)
	$(mockgen) -destination $@ -package mocksocketcan -source $<

.PHONY: stringer-generate
stringer-generate: \
	pkg/descriptor/sendtype_string.go \
	pkg/socketcan/errorclass_string.go \
	pkg/socketcan/protocolviolationerrorlocation_string.go \
	pkg/socketcan/protocolviolationerror_string.go \
	pkg/socketcan/controllererror_string.go \
	pkg/socketcan/transceivererror_string.go

%_string.go: %.go $(stringer)
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
