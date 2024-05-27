package plugin_loader

import (
	"errors"
	"fmt"
	"plugin"
	"strings"

	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
)

func LoadPlugins[t any](paths string, typeName string) (map[string]t, error) {
	pluginInstances := make(map[string]t, 0)
	for _, soFilePath := range strings.Split(paths, ",") {
		if soFilePath == "" {
			continue
		}
		plug, err := plugin.Open(soFilePath)
		if err != nil {
			return nil, errors.New("error loading plugin: " + err.Error())
		}
		plugSymbol, err := plug.Lookup(typeName)
		if err != nil {
			return nil, fmt.Errorf("error finding %s in plugin %s: %v", typeName, soFilePath, err)
		}

		pluginCheck, ok := plugSymbol.(plugin_models.Plugin)
		if !ok {
			return nil, errors.New("plugin must implement Plugin interface")
		}

		pluginId := pluginCheck.GetPluginId()
		if _, ok := pluginInstances[pluginId]; ok {
			return nil, errors.New("duplicate plugin name: " + pluginId)
		}

		aiApi, ok := plugSymbol.(t)
		if !ok {
			return nil, errors.New("unexpected type from module symbol")
		}

		pluginInstances[pluginId] = aiApi
	}
	return pluginInstances, nil
}
