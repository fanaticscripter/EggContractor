MAKEFLAGS += -j4

.PHONY: all init dist html webpack postcss data static clean serve

all: clean
	$(MAKE) dist

init:
	yarn install

dist: html static
	../_tools/build.py dist --additional 'Always Together.otf'

html: webpack postcss
	../_tools/build.py html

webpack: data src/lib/schema_pb.js
	yarn webpack

src/lib/schema_pb.js: src/lib/schema.proto
	protoc --proto_path=src/lib --js_out=import_style=commonjs,binary:src/lib src/lib/schema.proto

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

data:
	go run .

static:
	mkdir -p dist
	cp -f 'static/Always Together.otf' 'dist/Always Together.otf'

clean:
	@$(RM) -r dist/index.html dist/*.*.css dist/*.*.js dist/manifest.*.json

serve:
	cd dist && npx serve
