package models

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
)

func TestLoadSettings_TokenFromSecureJSONData(t *testing.T) {
	cfg := backend.DataSourceInstanceSettings{
		// authKind omitted -> 0 == AuthTypeToken. Note: no "token" in jsonData.
		JSONData:                []byte(`{"uri":"db-host:443","database":"grafana"}`),
		DecryptedSecureJSONData: map[string]string{"token": "AstraCS:secret"},
	}

	s, err := LoadSettings(cfg)
	require.NoError(t, err)

	require.Equal(t, AuthTypeToken, s.AuthKind)
	require.Equal(t, "AstraCS:secret", s.Token, "token must be loaded from secureJsonData via mapstructure")
	require.Equal(t, "db-host:443", s.URI)
	require.Equal(t, "grafana", s.Database)
}

func TestLoadSettings_PasswordFromSecureJSONData(t *testing.T) {
	cfg := backend.DataSourceInstanceSettings{
		JSONData:                []byte(`{"authKind":1,"user":"alice","authEndpoint":"auth-host"}`),
		DecryptedSecureJSONData: map[string]string{"password": "s3cret"},
	}

	s, err := LoadSettings(cfg)
	require.NoError(t, err)

	require.Equal(t, AuthTypeCredentials, s.AuthKind)
	require.Equal(t, "s3cret", s.Password, "password must be loaded from secureJsonData via mapstructure")
	require.Equal(t, "alice", s.UserName)
}

// TestLoadSettings_SecretsNotReadFromJSONData documents that token/password are
// secureJsonData only: a value placed in jsonData must NOT populate them
// (json:"-"), so secrets can never leak in through unencrypted jsonData.
func TestLoadSettings_SecretsNotReadFromJSONData(t *testing.T) {
	cfg := backend.DataSourceInstanceSettings{
		JSONData:                []byte(`{"token":"leaked","password":"leaked"}`),
		DecryptedSecureJSONData: map[string]string{},
	}

	s, err := LoadSettings(cfg)
	require.NoError(t, err)

	require.Empty(t, s.Token, "token must not be read from jsonData")
	require.Empty(t, s.Password, "password must not be read from jsonData")
}
