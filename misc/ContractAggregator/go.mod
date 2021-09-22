module github.com/fanaticscripter/EggContractor/misc/ContractAggregator

go 1.16

replace github.com/fanaticscripter/EggContractor => ../..

require (
	github.com/fanaticscripter/EggContractor v0.0.0-20210922143229-808c77d03657
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	google.golang.org/protobuf v1.27.1
)
