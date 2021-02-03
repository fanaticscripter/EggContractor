MAKEFLAGS += -j4

.PHONY: all init dist html wasm webpack postcss static clean serve

all: dist

init:
	yarn install

dist: html static
	../_tools/build.py dist --additional data.json

html: webpack postcss
	cp -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" dist
	../_tools/build.py html

wasm:
	../_tools/build.py wasm

webpack: wasm
	yarn webpack

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

static:
	mkdir -p dist
	cp -f static/data.json dist/data.json

clean:
	@$(RM) -r dist/index.html dist/*.*.js dist/wasm_exec.js dist/*.*.css dist/*.*.wasm dist/manifest.*.json dist/data.json

serve:
	cd dist && npx serve
