package config

import (
	"strings"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/pkg/errors"
)

const (
	// SettingListen is the config key for the listen address
	SettingListen = "listen"
	// SettingListenDefault is the default value for the listen address
	SettingListenDefault = ":9060"

	// SettingMessageBusURI is the config key for the message bus URI
	SettingMessageBusURI = "message_bus_uri"
	// SettingMessageBusURIDefault is the default value for the message bus URI
	SettingMessageBusURIDefault = "amqp://user:password@localhost:5672//"

	// SettingTenant is the config key for the tenant identifier
	SettingTenant = "tenant"
	// SettingTenantDefault is the default value for the tenant
	SettingTenantDefault = "default"
)

var (
	// Defaults are the default configuration settings
	Defaults = []config.Default{
		{Key: SettingListen, Value: SettingListenDefault},
		{Key: SettingMessageBusURI, Value: SettingMessageBusURIDefault},
		{Key: SettingTenant, Value: SettingTenantDefault},
	}
)

// Init initializes the configuration from the given config file
func Init(configPath string) error {
	err := config.FromConfigFile(configPath, Defaults)
	if err != nil {
		return errors.Wrap(err, "error loading configuration file")
	}

	// Enable setting config values by environment variables
	config.Config.SetEnvPrefix("RATING_AGENT_HEP")
	config.Config.AutomaticEnv()
	config.Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	return nil
}
