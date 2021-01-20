MAKEFLAGS += -j4

.PHONY: all init dist html postcss static clean serve

all: dist

init:
	yarn install

dist: html static
	../_tools/build.py dist --additional static

html: postcss
	go run .
	../_tools/build.py html

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

static:
	rm -rf dist/static
	mkdir -p dist/static
	ln ../../../static/egginc/afx_*.png dist/static

clean:
	@$(RM) -r src/index.html dist/index.html dist/*.*.css dist/manifest.*.json dist/static/

serve:
	cd dist && npx serve
