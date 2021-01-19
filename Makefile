MAKEFLAGS += -j4

.PHONY: all dev go protobuf webpack webpack-dev postcss postcss-dev fmt serve serve-prod docker clean

all: go webpack postcss

dev: go webpack-dev postcss-dev
	rm -f public/egginc public/egginc-extras
	ln -sf ../static/egginc public/egginc
	ln -sf ../static/egginc-extras public/egginc-extras

go: protobuf
	go build

protobuf:
	protoc --proto_path=. --go_out=paths=source_relative:. api/egginc.proto
	protoc --proto_path=. --go_out=paths=source_relative:. solo/pb/solo.proto
	gofumpt -w api/egginc.pb.go solo/pb/solo.pb.go

webpack:
	yarn webpack --config webpack.prod.js

webpack-dev:
	yarn webpack --config webpack.dev.js

postcss:
	yarn postcss --env=production css/app.css -o static/app.css

postcss-dev:
	yarn postcss --env=development css/app.css -o public/app.css

fmt:
	gofumpt -w .
	clang-format -i api/egginc.proto solo/pb/solo.proto

# Hot-reloading server based on entr(1).
serve:
	while true; do { echo EggContractor; find templates; find static; } | entr -dr ./EggContractor serve --dev; [ $$? = 2 ] || break; done

# Same as serve, except serving the production version.
serve-prod:
	while true; do { echo EggContractor; find templates; find static; } | entr -dr ./EggContractor serve; [ $$? = 2 ] || break; done

docker:
	docker build -t fanaticscripter/eggcontractor .

clean:
	@$(RM) EggContractor static/*.*.js static/*.*.css static/manifest.*.json
	@$(RM) -r node_modules/.cache/webpack
