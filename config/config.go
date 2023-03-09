package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

var (
	DefaultConfigPath = "./config.toml"
)

// Config defines all necessary configuration parameters.
type INDEPConfig struct {
	Telegram_Token                string `toml:"telegram_token"`
	Telegram_Chat_Id              string `toml:"telegram_chat_id"`
	Pypd_Api_Key                  string `toml:"pypd_api_key"`
	Pypd_Service_Key              string `toml:"pypd_service_key"`
	Pypd_Href                     string `toml:"pypd_herf"`
	Height_Increasing_Time_Period int64  `toml:"height_increasing_time_period"`
	Missing_Block_Trigger         int64  `toml:"missing_block_trigger"`
	Free_Disk_Trigger             int64  `toml:"free_disk_trigger"`
	Node_Home_Dir                 string `toml:"node_home_dir"`
}

type NODEAPPConfig struct {
	MIN_GAS_PRICES string      `toml:"minimum-gas-prices"`
	PRUNING        string      `toml:"pruning"`
	GRPC           *GRPCConfig `toml:"grpc"`
	API            *APIConfig  `toml:"api"`
}

type NODEConfig struct {
	NODE_NAME       string         `toml:"moniker"`
	PRIV_VAL_PATH   string         `toml:"priv_validator_key_file"`
	PRIV_STATE_PATH string         `toml:"priv_validator_state_file"`
	TX_INDEX        *TXINDEXConfig `toml:"tx_index"`
	RPC             *RPCConfig     `toml:"rpc"`
}

type GRPCConfig struct {
	Enable   bool   `toml:"enable"`
	Address  string `toml:"address"`
	Node_rpc string `toml:"node_grpc"`
}
type APIConfig struct {
	Enable  bool   `toml:"enable"`
	Address string `toml:"address"`
}
type RPCConfig struct {
	Address string `toml:"laddr"`
}
type TXINDEXConfig struct {
	Indexer string `toml:"indexer"`
}

// SetupConfig takes the path to a configuration file and returns the properly parsed configuration.
func Read(configPath string) (*INDEPConfig, error) {
	if configPath == "" {
		return nil, fmt.Errorf("empty configuration path")
	}

	log.Debug().Msg("reading config file")

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}

	return ParseString(configData)
}

// ParseString attempts to read and parse  config from the given string bytes.
// An error reading or parsing the config results in a panic.
func ParseString(configData []byte) (*INDEPConfig, error) {
	var cfg INDEPConfig

	log.Debug().Msg("parsing config data")

	err := toml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %s", err)
	}

	return &cfg, nil
}

func ReadNodeConfig(Node_configPath string, App_configPath string) (*NODEConfig, *NODEAPPConfig, error) {
	if Node_configPath == "" {
		return nil, nil, fmt.Errorf("empty Node_config configuration path")
	}
	if App_configPath == "" {
		return nil, nil, fmt.Errorf("empty App_config configuration path")
	}
	log.Debug().Msg("reading config file")

	node_configData, err := ioutil.ReadFile(Node_configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read config: %s", err)
	}
	app_configData, err := ioutil.ReadFile(App_configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read config: %s", err)
	}

	return ParseNodeConfigString(node_configData, app_configData)
}

func ParseNodeConfigString(Node_configData []byte, App_configData []byte) (*NODEConfig, *NODEAPPConfig, error) {

	var n_cfg NODEConfig
	var a_cfg NODEAPPConfig
	log.Debug().Msg("parsing config data")

	err := toml.Unmarshal(Node_configData, &n_cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode config: %s", err)
	}
	err = toml.Unmarshal(App_configData, &a_cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode config: %s", err)
	}
	return &n_cfg, &a_cfg, nil
}
