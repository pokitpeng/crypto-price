package config

import (
	"crypto-price/exchanges"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

var (
	DefaultConfigPath = path.Join(os.Getenv("HOME"), ".crypto-price")
	DefaultConfigFile = "config.json"
)

type Token struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Exchange  string  `json:"exchange"`
	Displayed bool    `json:"displayed"`
}

type Config struct {
	Tokens      []Token `json:"tokens"`
	RefreshTime int     `json:"refresh_time"` // 刷新间隔（秒）
	ProxyURL    string  `json:"proxy_url"`
	ActiveToken string  `json:"active_token"` // 当前显示的代币
}

func LoadConfig() (*Config, error) {
	configPath := filepath.Join(DefaultConfigPath, DefaultConfigFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefaultConfig()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	if err := os.MkdirAll(DefaultConfigPath, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(DefaultConfigPath, DefaultConfigFile), data, 0644)
}

func createDefaultConfig() (*Config, error) {
	config := &Config{
		Tokens: []Token{
			{Symbol: "BTC", Name: "BTC-Binance", Exchange: exchanges.ExchangeBinance, Displayed: true},
			{Symbol: "ETH", Name: "ETH-Binance", Exchange: exchanges.ExchangeBinance, Displayed: false},
			{Symbol: "DOGE", Name: "DOGE-Binance", Exchange: exchanges.ExchangeBinance, Displayed: false},
			{Symbol: "GRASS", Name: "GRASS-Bitget", Exchange: exchanges.ExchangeBitget, Displayed: false},
			{Symbol: "TRUMP", Name: "TRUMP-XT", Exchange: exchanges.ExchangeXT, Displayed: false},
			{Symbol: "EIGEN", Name: "EIGEN-XT", Exchange: exchanges.ExchangeXT, Displayed: false},
		},
		RefreshTime: 30,
		ProxyURL:    "",
		ActiveToken: "BTC",
	}

	if err := SaveConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}
