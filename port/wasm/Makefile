MAKEFLAGS ?= -j4

.PHONY: all init artifact-list init-artifact-list mission-list init-mission-list rockets-tracker init-rockets-tracker past-contracts init-past-contracts

all: artifact-list mission-list past-contracts rockets-tracker

init: init-artifact-list init-mission-list init-past-contracts init-rockets-tracker

artifact-list:
	$(MAKE) -C artifact-list

init-artifact-list:
	$(MAKE) -C artifact-list init

mission-list:
	$(MAKE) -C mission-list

init-mission-list:
	$(MAKE) -C mission-list init

rockets-tracker:
	$(MAKE) -C rockets-tracker

init-rockets-tracker:
	$(MAKE) -C rockets-tracker init

past-contracts:
	$(MAKE) -C past-contracts

init-past-contracts:
	$(MAKE) -C past-contracts init
