package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/mitchellh/mapstructure"
)

type Settings struct {
	URI   string `json:"uri"`
	Token string `json:"token"`
}

func LoadSettings(config backend.DataSourceInstanceSettings) (Settings, error) {
	settings := Settings{}

	if err := json.Unmarshal(config.JSONData, &settings); err != nil {
		return settings, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}

	if config.DecryptedSecureJSONData == nil {
		return settings, nil
	}

	uri := settings.URI
	if err := mapstructure.Decode(config.DecryptedSecureJSONData, &settings); err != nil {
		return settings, fmt.Errorf("could not unmarshal secure settings: %w", err)
	}
	settings.URI = uri

	return settings, nil
}
