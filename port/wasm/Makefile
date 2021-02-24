MAKEFLAGS ?= -j4

targets = artifact-explorer artifact-sandbox artifact-list consumption-sheet mission-list past-contracts rockets-tracker
init-targets = $(addprefix init-,$(targets))
clean-targets = $(addprefix clean-,$(targets))

.PHONY: all init clean $(targets) $(init-targets) $(clean-targets)

all: $(targets)

init: $(init-targets)

clean: $(clean-targets)

$(init-targets): init-%:
	$(MAKE) -C $(patsubst init-%,%,$@) init

$(targets): %:
	$(MAKE) -C $@

$(clean-targets): clean-%:
	$(MAKE) -C $(patsubst clean-%,%,$@) clean
