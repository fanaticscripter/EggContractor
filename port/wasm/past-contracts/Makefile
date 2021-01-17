MAKEFLAGS += -j4

.PHONY: all init dist html wasm webpack postcss clean serve

all: dist

init:
	yarn install

dist: html dist/contracts.csv
	../_tools/build.py dist --additional contracts.csv

html: webpack postcss
	cp -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" dist
	../_tools/build.py html

wasm:
	../_tools/build.py wasm

webpack: wasm
	yarn webpack

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

dist/contracts.csv: ../../../misc/ContractAggregator/data/contracts.csv
	cp -f ../../../misc/ContractAggregator/data/contracts.csv dist

clean:
	@$(RM) dist/index.html dist/*.*.js dist/wasm_exec.js dist/*.*.css dist/*.*.wasm dist/manifest.*.json dist/contracts.csv

serve:
	cd dist && npx serve
