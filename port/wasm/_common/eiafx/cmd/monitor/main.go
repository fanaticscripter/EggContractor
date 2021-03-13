package main

import (
	"io/ioutil"

	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/port/wasm/_common/eiafx"
)

const _dataFilePath = "eiafx-config.json"

func main() {
	var errorOnDiff bool
	flag.BoolVar(&errorOnDiff, "error-on-diff", false,
		"exit with code 1 if fetched config differs from cached version")
	flag.Parse()

	log.SetLevel(log.InfoLevel)
	err := eiafx.LoadConfig()
	if err != nil {
		log.Fatalf("error loading cached config: %s", err)
	}
	existingConfig := eiafx.Config

	req := &api.ArtifactsConfigurationRequestPayload{
		ClientVersion: api.ClientVersion,
	}
	config := &api.ArtifactsConfigurationResponse{}
	err = api.RequestAuthenticated("/ei_afx/config", req, config)
	if err != nil {
		log.Fatal(err)
	}

	if diff := cmp.Diff(config, existingConfig, protocmp.Transform()); diff != "" {
		reportFunc := log.Warnf
		if errorOnDiff {
			reportFunc = log.Fatalf
		}
		reportFunc("config has diverged from cached version: %s", diff)
	} else {
		log.Info("config has not changed")
		return
	}

	marshaller := protojson.MarshalOptions{
		Indent: "  ",
	}
	encoded, err := marshaller.Marshal(config)
	if err != nil {
		log.Fatalf("error marshalling %+v: %s", config, err)
	}
	err = ioutil.WriteFile(_dataFilePath, encoded, 0o644)
	if err != nil {
		log.Fatalf("error writing to %s: %s", _dataFilePath, err)
	}
	log.Infof("config written to %s", _dataFilePath)
}
