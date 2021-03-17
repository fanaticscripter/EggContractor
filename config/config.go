package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ProgName       = "EggContractor"
	ConfigTemplate = `# Place this config file at ~/.config/EggContractor/config.toml.
#
# Alternatively, use the --config <config-file> option or the environment
# variable EGGCONTRACTOR_CONFIG_FILE=<config-file> to load config placed
# elsewhere.

[[players]]

# Player identifiers are configured here. Multiple accounts are supported. To
# add another account, use another [[players]] section, as demonstrated below.

# player.id, required.
#
# A unique account ID. To view your account ID, go to Main Menu -> Settings ->
# Privacy & Data, and the ID should be in the bottom left corner.
#id = "EI1234567890123456"

# player.device_id, optional.
#divice_id = "880684B9-756B-451D-944D-48D3453822B7"

# Example of another account:
#[[players]]
#id = "EI9876543210987654"

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
	Players      []PlayerConfig
	Database     DatabaseConfig
	Notification NotificationConfig
}

type PlayerConfig struct {
	Id       string
	DeviceId string `mapstructure:"device_id"`
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

type Deprecations struct {
	HasLegacyPlayerField bool
}

func (c *Config) ResolveAndValidate() error {
	var err error
	if c.Player.Id != "" {
		// Prepend Player to Players, if the player isn't already in Players.
		player := c.Player
		duplicate := false
		for _, p := range c.Players {
			if p.Id == player.Id {
				duplicate = true
				break
			}
		}
		if !duplicate {
			c.Players = append([]PlayerConfig{player}, c.Players...)
		}
		log.Warn("DeprecationWarning: player deprecated, use players instead; see `EggContractor config-template`")
	}
	if len(c.Players) == 0 {
		return errors.New("at least one player required")
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

func (c Config) LockPath(lockName string) string {
	return filepath.Join(filepath.Dir(c.Database.Path), fmt.Sprintf("%s.lock", lockName))
}

func (c Config) MultiAccountMode() bool {
	return len(c.Players) > 1
}

func (c Config) Deprecations() Deprecations {
	return Deprecations{
		HasLegacyPlayerField: c.HasLegacyPlayerField(),
	}
}

func (c Config) HasLegacyPlayerField() bool {
	return c.Player.Id != ""
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
