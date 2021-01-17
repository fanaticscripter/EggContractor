module github.com/fanaticscripter/EggContractor/port/wasm/past-contracts

go 1.15

replace github.com/fanaticscripter/EggContractor => ../../..

require (
	github.com/fanaticscripter/EggContractor v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/net v0.0.0-20201029221708-28c70e62bb1d
	google.golang.org/protobuf v1.25.0
)
