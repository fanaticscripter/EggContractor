package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	ProgName       = "EggContractor"
	ConfigTemplate = `# Place this config file at ~/.config/EggContractor/config.toml.
#
# Alternatively, use the --config <config-file> option or the environment
# variable EGGCONTRACTOR_CONFIG_FILE=<config-file> to load config placed
# elsewhere.

[player]
# player.id, required.
#
# Your unique player ID. To view your player ID, go to Main Menu -> Settings ->
# Privacy & Data, and the ID should be in the bottom left corner. On iOS (at
# least when signed in via Game Center) this would be of the format G:1234567890.
# Copy the string verbatim for this field.
#id = "G:1234567890"

[database]
# database.path, optional.
#
# The path to the SQLite3 database for all data storage needs of this program.
# Should be absolute (tilde expansion for current user is allowed).
#
# Default is ~/.local/share/EggContractor/data.db.
#path = ""
`
)

const _defaultDatabasePath = "~/.local/share/" + ProgName + "/data.db"

type Config struct {
	Player struct {
		Id string
	}
	Database struct {
		Path string
	}
}

func (c *Config) ResolveAndValidate() error {
	var err error
	if c.Player.Id == "" {
		return errors.New("player.id required")
	}
	if c.Database.Path == "" {
		c.Database.Path = _defaultDatabasePath
	}
	c.Database.Path, err = tildeExpand(c.Database.Path)
	if err != nil {
		return err
	}
	return nil
}

func tildeExpand(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.Wrapf(err, "failed to tilde expand %#v", path)
		}
		if path == "~" {
			return homedir, nil
		}
		return filepath.Join(homedir, path[2:]), nil
	} else {
		return path, nil
	}
}
