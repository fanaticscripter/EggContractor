MAKEFLAGS += -j4

.PHONY: all init dist html wasm webpack postcss static test clean serve

all: clean
	$(MAKE) dist

init:
	yarn install

dist: html
	../_tools/build.py dist --additional static

html: webpack postcss static
	cp -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" dist
	../_tools/build.py html

wasm:
	../_tools/build.py wasm

webpack: wasm
	yarn webpack

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

static:
	mkdir -p dist/static
	cp -f static/favicon* dist/static

# Uses go_js_wasm_exec from https://github.com/agnivade/wasmbrowsertest.
test:
	GOOS=js GOARCH=wasm go test

clean:
	@$(RM) -r dist/index.html dist/*.*.js dist/wasm_exec.js dist/*.*.css dist/*.*.wasm dist/manifest.*.json

serve:
	cd dist && npx serve
