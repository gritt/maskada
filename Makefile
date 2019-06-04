include Makefile.vars

.SILENT:
.DEFAULT_GOAL := help

.PHONY: help
help:
	$(info Available Commands:)
	$(info -> be-install                 Install backend dependencies)
	$(info -> be-test                    Run backend tests)
	$(info -> be-run                     Run backend locally)
	$(info -> be-build                   Generate backend build)
	$(info -> be-clean                   Cleanup backend build files)

.PHONY: be-install
#

.PHONY: be-test
#

.PHONY: be-run
#

.PHONY: be-build
#

.PHONY: be-clean

# ignore unknown commands
%:
    @:
