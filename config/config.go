package config

import (
	"github.com/mendersoftware/go-lib-micro/config"
)

const (
	// SettingListen is the config key for the listen address
	SettingListen = "listen"
	// SettingListenDefault is the default value for the listen address
	SettingListenDefault = ":9060"
)

var (
	// Defaults are the default configuration settings
	Defaults = []config.Default{
		{Key: SettingListen, Value: SettingListenDefault},
	}
)
