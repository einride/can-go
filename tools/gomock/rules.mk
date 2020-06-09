gomock_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
mockgen := $(gomock_cwd)/bin/mockgen

$(mockgen): $(gomock_cwd)/../../go.mod
	cd $(gomock_cwd)/../.. && go build -o $@ github.com/golang/mock/mockgen
