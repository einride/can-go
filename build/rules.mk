SHELL := /bin/bash

BUILD_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
FILES_DIR := $(BUILD_DIR)/files

UNAME := $(shell uname -s)
UNAME_LOWERCASE := $(shell uname -s | tr '[:upper:]' '[:lower:]')

ifeq ($(UNAME),Linux)
else ifeq ($(UNAME),Darwin)
else
$(error This Makefile only supports Linux and OSX build agents.)
endif

ifneq ($(shell uname -m),x86_64)
$(error This Makefile only supports x86_64 build agents.)
endif

ifneq ($(shell which curl >/dev/null; echo $$?),0)
$(error cURL not installed. This Makefile requires cURL.)
endif

ifneq ($(shell which realpath >/dev/null; echo $$?),0)
$(error Coreutils not installed. OSX users run: brew install coreutils)
endif
