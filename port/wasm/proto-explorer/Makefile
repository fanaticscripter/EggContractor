MAKEFLAGS += -j4

.PHONY: all init dist html webpack proto postcss doc clean serve

all: clean
	$(MAKE) dist

init:
	yarn install

dist: html doc
	../_tools/build.py dist --additional doc

html: webpack postcss
	../_tools/build.py html

webpack: proto
	yarn webpack

proto:
	yarn run pbjs -p ../../../misc/protobuf -t static-module -w es6 -o src/lib/proto.js ../../../misc/protobuf/*.proto

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

doc:
	mkdir -p dist/doc
	protoc --doc_out=html,index.html:./dist/doc --proto_path ../../../misc/protobuf ../../../misc/protobuf/*.proto

clean:
	@$(RM) -r dist/*.html dist/*.*.css dist/*.*.js dist/manifest.*.json dist/doc

serve:
	# Symlink proto-explorer to . so that absolute links of assets still work
	# even though we're serving from a wrong root.
	#
	# Note that /doc isn't correctly served by `npx serve --single`. I attempted
	# to use a serve.json instead but ran into
	# https://github.com/vercel/serve/issues/636. Maybe I should switch to
	# superstatic.
	cd dist && \
	rm -f proto-explorer && \
	ln -s . proto-explorer && \
	npx serve --single
