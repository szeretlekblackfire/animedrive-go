package config

import (
	"encoding/json"
	"io/ioutil"
)

type HeadersConfig struct {
	BaseHeaders   map[string]string `json:"baseHeaders"`
	PlayerHeaders map[string]string `json:"playerHeaders"`
}

func LoadHeaders(filePath, headerType string) (map[string]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var headersConfig HeadersConfig
	err = json.Unmarshal(file, &headersConfig)
	if err != nil {
		return nil, err
	}

	switch headerType {
	case "playerHeaders":
		return headersConfig.PlayerHeaders, nil
	default:
		return headersConfig.BaseHeaders, nil
	}
}
