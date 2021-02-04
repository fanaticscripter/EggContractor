MAKEFLAGS += -j4

.PHONY: all init dist html postcss clean serve

all: clean
	$(MAKE) dist

init:
	yarn install

dist: html
	../_tools/build.py dist

html: postcss
	go run .
	../_tools/build.py html

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

clean:
	@$(RM) -r src/index.html dist/index.html dist/*.*.css dist/manifest.*.json

serve:
	cd dist && npx serve
