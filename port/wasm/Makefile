MAKEFLAGS ?= -j4

.PHONY: all init artifact-list init-artifact-list rockets-tracker init-rockets-tracker past-contracts init-past-contracts

all: artifact-list past-contracts

init: init-artifact-list init-past-contracts

artifact-list:
	$(MAKE) -C artifact-list

init-artifact-list:
	$(MAKE) -C artifact-list init

rockets-tracker:
	$(MAKE) -C rockets-tracker

init-rockets-tracker:
	$(MAKE) -C rockets-tracker init

past-contracts:
	$(MAKE) -C past-contracts

init-past-contracts:
	$(MAKE) -C past-contracts init
