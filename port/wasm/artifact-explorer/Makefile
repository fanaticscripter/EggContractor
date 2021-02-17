MAKEFLAGS += -j4

.PHONY: all init dist html webpack postcss data clean serve

all: clean
	$(MAKE) dist

init:
	yarn install

dist: html
	../_tools/build.py dist

html: webpack postcss
	../_tools/build.py html

webpack: data
	yarn webpack

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

data:
	go run .

clean:
	@$(RM) -r dist/index.html dist/*.*.css dist/*.*.js dist/manifest.*.json

serve:
	cd dist && npx serve
