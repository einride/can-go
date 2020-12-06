commitlint_dir := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
commitlint_version := 9.1.1
commitlint := $(commitlint_dir)/node_modules/.bin/commitlint

$(commitlint):
	$(info [commitlint] installing command version $(commitlint_version)...)
	@npm install --no-save --no-audit --prefix $(commitlint_dir) @commitlint/config-conventional@$(commitlint_version)
	@npm install --no-save --no-audit --prefix $(commitlint_dir) @commitlint/cli@$(commitlint_version)

.PHONY: commitlint
commitlint: $(commitlint)
	$(info [$@] linting commit messages...)
	@git fetch --tags
	@$(commitlint) -x "$(commitlint_dir)/node_modules/@commitlint/config-conventional" --from origin/master --to HEAD
