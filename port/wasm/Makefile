MAKEFLAGS ?= -j4

targets = artifact-explorer artifact-sandbox artifact-list consumption-sheet mission-list past-contracts proto-explorer rockets-tracker
init-targets = $(addprefix init-,$(targets))
clean-targets = $(addprefix clean-,$(targets))

.PHONY: all init clean $(targets) $(init-targets) $(clean-targets) update-loot-data

all: $(targets)

init: $(init-targets) update-loot-data

update-loot-data:
	curl -o _common/data/mission_reward_count.json https://api.ei.mikit.app/mission_reward_count.json

clean: $(clean-targets)

$(init-targets): init-%:
	$(MAKE) -C $(patsubst init-%,%,$@) init

$(targets): %:
	$(MAKE) -C $@

$(clean-targets): clean-%:
	$(MAKE) -C $(patsubst clean-%,%,$@) clean
