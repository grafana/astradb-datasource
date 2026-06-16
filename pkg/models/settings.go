package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/mitchellh/mapstructure"
)

type AuthType uint8

const (
	AuthTypeToken AuthType = iota
	AuthTypeCredentials
)

type Settings struct {
	URI          string   `json:"uri"`
	Token        string   `json:"token"`
	GRPCEndpoint string   `json:"grpcEndpoint"`
	AuthEndpoint string   `json:"authEndpoint"`
	UserName     string   `json:"user"`
	Password     string   `json:"password"`
	Secure       bool     `json:"secure"`
	AuthKind     AuthType `json:"authKind"`
}

func LoadSettings(config backend.DataSourceInstanceSettings) (Settings, error) {
	settings := Settings{}

	if err := json.Unmarshal(config.JSONData, &settings); err != nil {
		return settings, fmt.Errorf("could not unmarshal DataSourceInfo json: %w", err)
	}

	if config.DecryptedSecureJSONData == nil {
		return settings, nil
	}

	if settings.AuthKind == AuthTypeToken {
		secureSettings := Settings{}
		if err := mapstructure.Decode(config.DecryptedSecureJSONData, &secureSettings); err != nil {
			return settings, fmt.Errorf("could not unmarshal secure settings: %w", err)
		}
		settings.Token = secureSettings.Token
	}

	if settings.AuthKind == AuthTypeCredentials {
		secureSettings := Settings{}
		if err := mapstructure.Decode(config.DecryptedSecureJSONData, &secureSettings); err != nil {
			return settings, fmt.Errorf("could not unmarshal secure settings: %w", err)
		}
		settings.Password = secureSettings.Password
	}

	return settings, nil
}
