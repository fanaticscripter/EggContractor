MAKEFLAGS += -j4

.PHONY: all init dist html webpack postcss data static clean serve

all: clean
	$(MAKE) dist

init:
	yarn install

dist: html static
	../_tools/build.py dist --additional static

html: webpack postcss
	../_tools/build.py html

webpack: data
	yarn webpack

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

data:
	go run .

static:
	mkdir -p dist/static
	cp -f static/favicon* dist/static

clean:
	@$(RM) -r dist/index.html dist/*.*.css dist/*.*.js dist/manifest.*.json dist/static

serve:
	cd dist && npx serve
