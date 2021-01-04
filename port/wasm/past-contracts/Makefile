MAKEFLAGS += -j4

.PHONY: all dist html wasm webpack postcss clean serve

all: dist

dist: html
	../_tools/build.py dist

html: webpack postcss
	cp -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" dist
	../_tools/build.py html

wasm:
	../_tools/build.py wasm

webpack: wasm
	yarn webpack

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

clean:
	@$(RM) dist/*.*.js dist/wasm_exec.js dist/*.*.css dist/*.*.wasm dist/manifest.*.json

serve:
	cd dist && npx serve
