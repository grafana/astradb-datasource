package main

import (
	"os"

	"github.com/grafana/astradb-datasource/pkg/plugin"
	"github.com/grafana/grafana-enterprise-sdk/plugins"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {

	pluginID := "grafana-astradb-datasource"

	if err := plugins.CheckEnterprisePluginLicense(pluginID); err != nil {
		log.DefaultLogger.Error(err.Error())
		return // Should never get here
	}

	if err := datasource.Manage(pluginID, plugin.NewDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
