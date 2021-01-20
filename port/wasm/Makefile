MAKEFLAGS += -j4

.PHONY: all init artifact-list init-artifact-list past-contracts init-past-contracts

all: artifact-list past-contracts

init: init-artifact-list init-past-contracts

artifact-list:
	$(MAKE) -C artifact-list

init-artifact-list:
	$(MAKE) -C artifact-list init

past-contracts:
	$(MAKE) -C past-contracts

init-past-contracts:
	$(MAKE) -C past-contracts init
