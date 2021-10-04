package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/fanaticscripter/EggContractor/api"
)

const (
	_apiEndpoint      = "/ei_afx/config"
	_rawDataFilePath  = "eiafx-config.txt"
	_jsonDataFilePath = "eiafx-config.json"
)

func main() {
	var errorOnDiff bool
	flag.BoolVar(&errorOnDiff, "error-on-diff", false,
		"exit with code 1 if fetched config differs from cached version")
	flag.Parse()

	log.SetLevel(log.InfoLevel)

	cachedRaw, err := ioutil.ReadFile(_rawDataFilePath)
	if err != nil {
		log.Fatalf("error reading cached config: %s", err)
	}
	existingConfig := &api.ArtifactsConfigurationResponse{}
	err = api.DecodeAPIResponse(_apiEndpoint, cachedRaw, existingConfig, true)
	if err != nil {
		log.Fatalf("error decoding cached config: %s", err)
	}

	cachedJson, err := ioutil.ReadFile(_jsonDataFilePath)
	if err != nil {
		log.Fatalf("error reading cached JSON config: %s", err)
	}

	req := &api.ArtifactsConfigurationRequestPayload{
		ClientVersion: api.ClientVersion,
	}
	raw, err := api.RequestRawPayload(_apiEndpoint, req)
	if err != nil {
		log.Fatal(err)
	}
	config := &api.ArtifactsConfigurationResponse{}
	err = api.DecodeAPIResponse(_apiEndpoint, raw, config, true)
	if err != nil {
		log.Fatal(err)
	}

	if diff := cmp.Diff(existingConfig, config, protocmp.Transform()); diff != "" {
		reportFunc := log.Warnf
		if errorOnDiff {
			reportFunc = log.Fatalf
		}
		reportFunc("config has diverged from cached version: %s", diff)
	} else {
		log.Info("config has not changed")
	}

	// Marshal with protojson first, then marshal with json.MarshalIndent again
	// to remove protojson-introduced indeterminism (double spaces in some
	// places).
	encoded, err := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}.Marshal(config)
	if err != nil {
		log.Fatalf("error marshalling %+v: %s", config, err)
	}
	var rawJson json.RawMessage = encoded
	encoded, err = json.MarshalIndent(rawJson, "", "  ")
	if err != nil {
		log.Fatalf("error re-marshalling %+v: %s", config, err)
	}
	encoded = append(encoded, '\n')

	if diff := cmp.Diff(cachedJson, encoded); diff != "" {
		if errorOnDiff {
			log.Fatalf("JSON-encoded config has diverged from cached version: %s", diff)
		}
	} else {
		log.Info("JSON-encoded config has not changed")
		return
	}

	err = ioutil.WriteFile(_rawDataFilePath, raw, 0o644)
	if err != nil {
		log.Fatalf("error writing to %s: %s", _rawDataFilePath, err)
	}
	log.Infof("raw payload written to %s", _rawDataFilePath)
	err = ioutil.WriteFile(_jsonDataFilePath, encoded, 0o644)
	if err != nil {
		log.Fatalf("error writing to %s: %s", _jsonDataFilePath, err)
	}
	log.Infof("config written to %s", _jsonDataFilePath)
}
