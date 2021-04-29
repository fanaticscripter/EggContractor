.PHONY: all init data build dist clean dev serve

all: clean
	$(MAKE) dist

init:
	yarn install

build: data src/lib/proto.js
	yarn build

src/lib/proto.js: ../../../misc/protobuf/*.proto
	yarn run pbjs -p ../../../misc/protobuf -t static-module -w es6 -o src/lib/proto.js ../../../misc/protobuf/*.proto
	yarn run pbts -o src/lib/proto.d.ts src/lib/proto.js
	# Fix issue of Long not being imported.
	# https://github.com/protobufjs/protobuf.js/issues/1533
	# https://github.com/protobufjs/protobuf.js/pull/1166
	sed -i "1 i import { Long } from 'protobufjs';" src/lib/proto.d.ts

data:
	go run .

dist: build
	../_tools/build.py dist --additional assets --additional static

clean:
	@$(RM) -r dist src/lib/catalog.json

dev:
	yarn dev

serve:
	yarn serve
