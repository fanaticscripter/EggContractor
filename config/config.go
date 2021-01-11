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

[notification]
  [notification.pushover]
  # notification.pushover.on, optional.
  #
  # If true, turn on notifications through Pushover; in that case, api_key and
  # user_key below are required.
  #
  # Default is false.
  #on = true

  # See https://pushover.net/api#registration.
  #api_key = "azGDORePK8gMaC0QOYAMyEEuzJnyUi"

  # See https://pushover.net/api#identifiers. A group key, or a comma-delimited
  # list of user keys may be used here instead.
  #user_key = "uQiRzpo4DXghDmr9QzzfQu27cmVRsG"
`
)

const _defaultDatabasePath = "~/.local/share/" + ProgName + "/data.db"

type Config struct {
	Player       PlayerConfig
	Database     DatabaseConfig
	Notification NotificationConfig
}

type PlayerConfig struct {
	Id string
}

type DatabaseConfig struct {
	Path string
}

type NotificationConfig struct {
	Pushover struct {
		On      bool
		APIKey  string `mapstructure:"api_key"`
		UserKey string `mapstructure:"user_key"`
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
	if c.Notification.Pushover.On {
		if c.Notification.Pushover.APIKey == "" {
			return errors.New("notification.pushover.api_key required when notification.pushover.on is true")
		}
		if c.Notification.Pushover.UserKey == "" {
			return errors.New("notification.pushover.user_key required when notification.pushover.on is true")
		}
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
