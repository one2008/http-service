package main

import (
	"bytes"
	"encoding/json"
	"http-service/cmd/log"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func GetApolloConfig(url string) (*Config, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respData := struct {
		Content string `json:"content"`
	}{}

	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, err
	}

	config := DefaultConfig()
	content := strings.Replace(respData.Content, "*{", "{", 1)
	content = strings.Replace(content, "}*", "}", -1)

	viper := viper.New()
	viper.SetConfigType("json")
	if err := viper.ReadConfig(bytes.NewBufferString(content)); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

func ParseConfig(env string) (*Config, error) {
	var err error
	viper := viper.New()
	viper.SetConfigType("yaml")
	viper.SetConfigName(env)
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("config/")
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

type Config struct {
	Log    *LogConfig    `mapstructure:"log"`
	Server *ServerConfig `mapstructure:"server"`
	Eth    *EthConfig    `mapstructure:"eth"`

	Database *DatabaseConfig `mapstructure:"database"`

	HttpAuth    *HTTPAuthConfig    `mapstructure:"http_auth"`
	PrivateSign *PrivateSignConfig `mapstructure:"private_sign"`
}

func DefaultConfig() *Config {
	return &Config{
		Log:         DefaultLogConfig(),
		Server:      DefaultServerConfig(),
		Eth:         DefaultEthConfig(),
		Database:    DefaultDatabaseConfig(),
		HttpAuth:    DefaultHTTPAuth(),
		PrivateSign: DefaultPrivateSignConfig(),
	}
}

type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr: "127.0.0.1:8081",
	}
}

type EthConfig struct {
	RPCUrl  string `mapstructure:"rpc_url"`
	ChainId int    `mapstructure:"chain_id"`
}

func DefaultEthConfig() *EthConfig {
	return &EthConfig{
		RPCUrl:  "https://rinkeby.infura.io/v3/3085573404af46be96ead378fcbe443d",
		ChainId: 4,
	}
}

type DatabaseConfig struct {
	DatabaseUrl     string `mapstructure:"database_url"`
	ConnMaxLifetime int    `mapstructure:"conn_max_life_time"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
}

func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DatabaseUrl:     "root:Aa123456@/db",
		ConnMaxLifetime: 180,
		MaxIdleConns:    10,
		MaxOpenConns:    10,
	}
}

type LogConfig struct {
	Level   string `mapstructure:"level"`
	Path    string `mapstructure:"path"`
	Brokers string `mapstructure:"brokers"`
	Topic   string `mapstructure:"topic"`
	AppName string `mapstructure:"app_name"`
	EnvName string `mapstructure:"env_name"`
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level: log.LogLevelInfo,
	}
}

type HTTPAuthConfig struct {
	Key    string `mapstructure:"key"`
	Secret string `mapstructure:"secret"`
}

func DefaultHTTPAuth() *HTTPAuthConfig {
	return &HTTPAuthConfig{
		Key:    "account",
		Secret: "account",
	}
}

type PrivateSignConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

func DefaultPrivateSignConfig() *PrivateSignConfig {
	return &PrivateSignConfig{
		AccessKeyID:     "MIHW8ZRZQb7UDwepjWre",
		SecretAccessKey: "h9xmGxlZayzlhlh694Rp",
	}
}
