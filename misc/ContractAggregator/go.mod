module github.com/fanaticscripter/EggContractor/misc/ContractAggregator

go 1.15

replace github.com/fanaticscripter/EggContractor => ../..

require (
	github.com/fanaticscripter/EggContractor v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	google.golang.org/protobuf v1.25.0
)
