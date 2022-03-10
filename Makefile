# package manager check
# check homebrew, pip3, pip
PM := $(shell { command -v brew || command -v pip3 || command -v pip; } 2> /dev/null)

all: pre_commit bin
	# Project initialization succeeded!

bin: # creates the demo protoc plugin for demonstrating uses of PG*
	go build -o ./bin/protoc-gen-connection .

pre_commit:
ifeq ($(CI), true)
	$(info CI detected, skip install pre-commit.)
else ifeq ($(PM),)
	$(error No valid package manager found. Please install brew (macOS) or pip (Linux) and try again.)
else ifeq ($(shell which pre-commit),)
	$(info Package manager found: $(PM))
	$(info Please make sure the directory which your pre-commit located is in your PATH.)
	# Installing pre-commit
	@- $(PM) install pre-commit
	# Pre-commit installed!
else
	$(info Pre-commit installed.)
endif

ifneq ($(CI), true)
	# Installing pre-commit hooks
	@- pre-commit install
	# Pre-commit hooks installed!
endif

lint: pre_commit
	@- pre-commit run --all-files

.PHONY: all bin pre_commit lint