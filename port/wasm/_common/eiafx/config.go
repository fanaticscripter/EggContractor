package eiafx

import (
	_ "embed"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed eiafx-config.json
var _eiafxConfigJSON []byte

var Config *api.ArtifactsConfigurationResponse

func LoadConfig() error {
	Config = &api.ArtifactsConfigurationResponse{}
	err := protojson.Unmarshal(_eiafxConfigJSON, Config)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling eiafx-config.json")
	}
	return nil
}
