# Build, test and install the detoo tools.
#
# Run `make` or `make help` to list the targets.

GO      ?= go
PREFIX  ?= /usr/local
DESTDIR ?=

BUILD_DIR   := build
COMMANDS    := $(notdir $(wildcard cmd/*))
BINARIES    := $(addprefix $(BUILD_DIR)/,$(COMMANDS))
COMPLETIONS := $(addsuffix .bash,$(BINARIES))

# The commands embed their help and their record templates, so those files are
# sources too: a binary is out of date when any of them changes.  The
# documentation site is left out: its node_modules holds tens of thousands of
# files, some of them .go files that are not part of this module.
EMBEDDED    := $(shell find internal -type f ! -name '*.go')
SOURCES     := $(shell find . -path ./docs/site -prune -o -type f -name '*.go' -print) $(EMBEDDED) go.mod go.sum

all: help

#############################################################################
# Build
#############################################################################

.PHONY: build
build: vet test $(BINARIES) $(COMPLETIONS) ## vet, test, and compile every command into build/

# The build directory shares its name with the phony `build` target, so it is
# created by the recipe rather than by a rule of its own.
$(BUILD_DIR)/%: $(SOURCES)
	@mkdir -p $(@D)
	$(GO) build -o $@ ./cmd/$*

# Each command generates its own shell completion, so it is built rather than
# kept in the repository.  This is a static pattern rule, which is explicit, so
# that build/adr.bash is not taken for a command by the pattern rule above.
$(COMPLETIONS): $(BUILD_DIR)/%.bash: $(BUILD_DIR)/%
	$< completion bash > $@

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: test
test:
	$(GO) test ./...

#############################################################################
# Deployment
#############################################################################

.PHONY: deploy
deploy: build ## install the commands and their shell completion under $(PREFIX)
	install -d $(DESTDIR)$(PREFIX)/bin $(DESTDIR)$(PREFIX)/etc/bash_completion.d
	install -m 755 $(BINARIES) $(DESTDIR)$(PREFIX)/bin
	install -m 644 $(COMPLETIONS) $(DESTDIR)$(PREFIX)/etc/bash_completion.d

# `make install` is what most people type, so it is kept, but out of the help.
.PHONY: install
install: deploy

.PHONY: release
release: build ## tag the current commit and publish a release
	git tag `git rev-parse --short HEAD` && git push --tags origin main
	gh release create `git rev-parse --short HEAD` --generate-notes

#############################################################################
# Documentation site
#############################################################################

# The site's build generates its command reference from the tools' help, with
# tools/sitegen, so it needs Go as well as Node.  It is not part of `build`,
# which needs Go alone.
.PHONY: site
site: ## build the documentation site into docs/site/dist, and check its links
	cd docs/site && npm ci && npm run build && npm run check-links

#############################################################################
# Housekeeping
#############################################################################

.PHONY: clean
clean: ## remove the build output
	rm -rf $(BUILD_DIR)

#############################################################################
# Help
#############################################################################

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[92m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DELETE_ON_ERROR:
