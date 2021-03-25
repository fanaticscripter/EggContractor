MAKEFLAGS += -j4

.PHONY: all init dist html webpack postcss data update-events-data clean serve

all: clean
	$(MAKE) dist

init: update-events-data
	yarn install

dist: html
	../_tools/build.py dist

html: webpack postcss
	../_tools/build.py html

webpack: data
	yarn webpack

postcss:
	yarn postcss --env=production src/app.css -o dist/app.css

data:
	go run .

update-events-data:
	curl -o src/events.json https://egginc-events-api.zw.workers.dev/

clean:
	@$(RM) -r dist/index.html dist/*.*.css dist/*.*.js dist/manifest.*.json

serve:
	cd dist && npx serve
