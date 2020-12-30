package cmd

import (
	"github.com/spf13/cobra"

	"github.com/fanaticscripter/EggContractor/web"
)

var (
	_serveDevMode  bool
	_serveBindAddr string
)

var _serveCommand = &cobra.Command{
	Use:     "serve",
	Short:   "Run web server",
	Args:    cobra.NoArgs,
	PreRunE: subcommandPreRunE,
	Run: func(cmd *cobra.Command, args []string) {
		web.Serve(web.ServerOptions{
			BindAddr: _serveBindAddr,
			Dev:      _serveDevMode || _debug,
		})
	},
}

func init() {
	_serveCommand.Flags().StringVarP(&_serveBindAddr, "bind", "b", ":8080",
		"bind address in the form of <host>:<port>, where host is optional and defaults to all interfaces")
	_serveCommand.Flags().BoolVarP(&_serveDevMode, "dev", "d", false, "enable dev mode; also implied by --debug")
	_rootCmd.AddCommand(_serveCommand)
}
