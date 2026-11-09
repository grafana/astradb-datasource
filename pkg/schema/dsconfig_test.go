package schema_test

import (
	_ "embed"
	"testing"

	"github.com/grafana/astradb-datasource/pkg/models"
	"github.com/grafana/dsconfig/schema"
)

//go:embed dsconfig.json
var configSchemaJSON []byte

//go:generate go test -run TestPlugin -generateArtifacts
func TestPlugin(t *testing.T) {
	schema.RunPluginTests(t, schema.PluginUnderTest{
		ID:                "grafana-astradb-datasource",
		ConfigSchemaJSON:  configSchemaJSON,
		SettingsJSONModel: models.Settings{},
		SecureKeys:        []string{"token", "password"},
	})
}