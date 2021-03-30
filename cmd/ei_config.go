package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/fanaticscripter/EggContractor/api"
)

var _eiConfigOutputFile string

var _eiConfigCommand = &cobra.Command{
	Use:     "ei-config",
	Short:   "Dump /ei/get_config response in JSON",
	Args:    cobra.NoArgs,
	PreRunE: subcommandPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &api.ConfigRequest{
			Rinfo: api.NewBasicRequestInfo(_config.Players[0].Id),
		}
		config := &api.ConfigResponse{}
		err := api.RequestAuthenticated("/ei/get_config", req, config)
		if err != nil {
			return err
		}

		encoded, err := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			EmitUnpopulated: true,
		}.Marshal(config)
		if err != nil {
			return errors.Wrap(err, "error marshalling response")
		}

		var rm json.RawMessage = encoded
		encoded, err = json.MarshalIndent(rm, "", "  ")
		if err != nil {
			return errors.Wrap(err, "error re-marshalling response")
		}
		encoded = append(encoded, '\n')

		if _eiConfigOutputFile == "" {
			fmt.Print(string(encoded))
		} else {
			err = ioutil.WriteFile(_eiConfigOutputFile, encoded, 0o644)
			if err != nil {
				return errors.Wrapf(err, "error writing to %s", _eiConfigOutputFile)
			}
		}

		return nil
	},
}

func init() {
	_rootCmd.AddCommand(_eiConfigCommand)
	_eiConfigCommand.Flags().StringVarP(&_eiConfigOutputFile, "output", "o", "", "output file path")
}
