package config

import (
	"strings"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/pkg/errors"
)

// Supported values for the State Manager setting
const (
	StateManagerMemory = "memory"
	StateManagerRedis  = "redis"
)

const (
	// SettingListenUDP is the config key for enabling the UDP protocol
	SettingListenUDP = "listen_udp"
	// SettingListenUDPDefault is the default value for the UDP protocol
	SettingListenUDPDefault = ":9060"

	// SettingListenTCP is the config key for enabling the TCP protocol
	SettingListenTCP = "listen_tcp"
	// SettingListenTCPDefault is the default value for the TCP protocol
	SettingListenTCPDefault = ":9060"

	// SettingMessageBusURI is the config key for the message bus URI
	SettingMessageBusURI = "message_bus_uri"
	// SettingMessageBusURIDefault is the default value for the message bus URI
	SettingMessageBusURIDefault = "amqp://user:password@localhost:5672//"

	// SettingStateManager is the config key for the state manager
	SettingStateManager = "state_manager"
	// SettingStateManagerDefault is the default value for the state manager
	SettingStateManagerDefault = StateManagerMemory

	// SettingRedisAddress is the config key for the redis address
	SettingRedisAddress = "redis_address"
	// SettingRedisAddressDefault is the default value for the redis address
	SettingRedisAddressDefault = "localhost:6379"

	// SettingRedisPassword is the config key for the redis password
	SettingRedisPassword = "redis_password"

	// SettingRedisDb is the config key for the redis database
	SettingRedisDb = "redis_db"
	// SettingRedisDbDefault is the default value for the redis database
	SettingRedisDbDefault = 0

	// SettingTenant is the config key for the tenant identifier
	SettingTenant = "tenant"
	// SettingTenantDefault is the default value for the tenant
	SettingTenantDefault = "default"

	// SettingSIPHeaderCaller is the SIP header used to extract the identifier of the caller
	SettingSIPHeaderCaller = "sip_header_caller"

	// SettingSIPHeaderCallee is the SIP header used to extract the identifier of the callee
	SettingSIPHeaderCallee = "sip_header_callee"

	// SettingSIPHeaderHistoryInfo is the SIP header used to extract the identifier of the caller on call forwarding
	SettingSIPHeaderHistoryInfo = "sip_header_history_info"
	// SettingSIPHeaderHistoryInfoDefault is the default value for sip_header_history_info
	SettingSIPHeaderHistoryInfoDefault = "History-Info"

	// SettingSIPHeaderHistoryInfoIndex is the index of the SIP header used to extract the identifier of the caller on call forwarding
	SettingSIPHeaderHistoryInfoIndex        = "sip_header_history_info_index"
	SettingSIPHeaderHistoryInfoIndexDefault = 1

	// SettingSIPLocalDomains is a comma separated list of local domains
	SettingSIPLocalDomains = "sip_local_domains"

	// SettingAccountTagMatchRegexp is a regular expression to extract the sip account
	SettingAccountTagMatchRegexp = "account_tag_match_regexp"

	// SettingProductTag is the product tag
	SettingProductTag = "product_tag"
	// SettingProductTagDefault is the product tag default value
	SettingProductTagDefault = "VOICE"

	// SettingTransactionTags is a comma separated list of transaction tags
	SettingTransactionTags = "transaction_tags"
)

var (
	// Defaults are the default configuration settings
	Defaults = []config.Default{
		{Key: SettingListenUDP, Value: SettingListenUDPDefault},
		{Key: SettingListenTCP, Value: SettingListenTCPDefault},
		{Key: SettingMessageBusURI, Value: SettingMessageBusURIDefault},
		{Key: SettingTenant, Value: SettingTenantDefault},
		{Key: SettingStateManager, Value: SettingStateManagerDefault},
		{Key: SettingRedisAddress, Value: SettingRedisAddressDefault},
		{Key: SettingRedisDb, Value: SettingRedisDbDefault},
		{Key: SettingProductTag, Value: SettingProductTagDefault},
		{Key: SettingSIPHeaderHistoryInfo, Value: SettingSIPHeaderHistoryInfoDefault},
		{Key: SettingSIPHeaderHistoryInfoIndex, Value: SettingSIPHeaderHistoryInfoIndexDefault},
	}
)

// Init initializes the configuration from the given config file
func Init(configPath string) error {
	if configPath != "" {
		err := config.FromConfigFile(configPath, Defaults)
		if err != nil {
			return errors.Wrap(err, "error loading configuration file")
		}
	}

	// Enable setting config values by environment variables
	config.Config.SetEnvPrefix("RATING_AGENT_HEP")
	config.Config.AutomaticEnv()
	config.Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	return nil
}
