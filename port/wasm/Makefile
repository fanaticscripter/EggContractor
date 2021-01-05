MAKEFLAGS += -j4

.PHONY: all init past-contracts init-past-contracts

all: past-contracts

init: init-past-contracts

past-contracts:
	$(MAKE) -C past-contracts

init-past-contracts:
	$(MAKE) -C past-contracts init
