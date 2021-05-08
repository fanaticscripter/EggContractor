.PHONY: all init data build dist clean dev serve

all: clean
	$(MAKE) dist

init:
	yarn install

build:
	yarn build

dist: build
	../_tools/build.py dist --additional assets --additional static

clean:
	@$(RM) -r dist

dev:
	yarn dev

serve:
	yarn serve
