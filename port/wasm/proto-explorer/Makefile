MAKEFLAGS += -j4

.PHONY: all init dist html webpack proto postcss doc clean serve

all: clean
	$(MAKE) dist

init:
	yarn install

dist: html doc
	../_tools/build.py dist --additional doc.html

html: webpack postcss
	../_tools/build.py html

webpack: proto
	yarn webpack

proto:
	yarn run pbjs -p ../../../misc/protobuf -t static-module -w es6 -o src/lib/proto.js ../../../misc/protobuf/*.proto

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

doc:
	mkdir -p dist
	protoc --doc_out=html,doc.html:./dist --proto_path ../../../misc/protobuf ../../../misc/protobuf/*.proto

clean:
	@$(RM) -r dist/*.html dist/*.*.css dist/*.*.js dist/manifest.*.json

serve:
	cd dist && npx serve
