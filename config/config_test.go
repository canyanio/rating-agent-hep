package config

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	assert.NotNil(t, Defaults)
}

func TestInitError(t *testing.T) {
	err := Init("/path/which/does/not/exist.yaml")
	assert.Error(t, err)
}

func TestInitConfigFile(t *testing.T) {
	const listen = "10.0.0.1:1234"

	f, err := ioutil.TempFile("", "config-*.yaml")
	assert.Nil(t, err)

	defer syscall.Unlink(f.Name())

	configFile := fmt.Sprintf("listen: %s\n", listen)
	ioutil.WriteFile(f.Name(), []byte(configFile), 0644)

	err = Init(f.Name())
	assert.Nil(t, err)

	settingListen := config.Config.GetString(SettingListen)
	assert.Equal(t, listen, settingListen)
}
