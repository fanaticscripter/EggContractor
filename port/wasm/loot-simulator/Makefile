.PHONY: all init data build dist clean dev serve

all: clean
	$(MAKE) dist

init:
	yarn install

data:
	go run .

build: data
	yarn build

dist: build
	../_tools/build.py dist --additional assets --additional static

clean:
	@$(RM) -r dist src/app-data.json

dev:
	yarn dev

serve:
	yarn serve
