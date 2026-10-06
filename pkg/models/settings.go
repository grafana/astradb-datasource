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
	// jsonData fields.
	URI          string   `json:"uri"`
	GRPCEndpoint string   `json:"grpcEndpoint"`
	AuthEndpoint string   `json:"authEndpoint"`
	UserName     string   `json:"user"`
	Secure       bool     `json:"secure"`
	AuthKind     AuthType `json:"authKind"`
	Database     string   `json:"database"`

	// secureJsonData fields. These are not part of jsonData (json:"-"); they are
	// loaded from the decrypted secrets via mapstructure in LoadSettings.
	Token    string `json:"-"`
	Password string `json:"-"`
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
