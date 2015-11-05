package main

/*

This is a generic and not terribly helpful plugin. However, it shows the basics
of what is needed in a backup plugin, and how they execute.

*/

import (
	"fmt"
	"plugin"
)

func main() {
	// Create an object representing this plugin, which is a type conforming to the Plugin interface
	dummy := DummyPlugin{
		// give it some authorship info
		meta: plugin.PluginInfo{
			Name:    "Dummy Plugin",
			Author:  "Stark & Wane",
			Version: "1.0.0",
			Features: plugin.PluginFeatures{
				Target: "yes",
				Store:  "yes",
			},
		},
	}

	// Run the plugin - the plugin framework handles all arg parsing, exit handling, error/debug formatting for you
	plugin.Run(dummy)
}

// Define my DummyPlugin type
type DummyPlugin struct {
	meta plugin.PluginInfo // needs a place to store metadata
}

// This function should be used to return the plugin's PluginInfo, however you decide to implement it
func (p DummyPlugin) Meta() plugin.PluginInfo {
	return p.meta
}

// Called when you want to back data up. Examine the ShieldEndpoint passed in, and perform actions accordingly
func (p DummyPlugin) Backup(endpoint plugin.ShieldEndpoint) (int, error) {
	data, err := endpoint.StringValue("data")
	if err != nil {
		return plugin.PLUGIN_FAILURE, err
	}

	return plugin.Exec(plugin.STDOUT, fmt.Sprintf("/bin/echo %s", data))
}

// Called when you want to restore data Examine the ShieldEndpoint passed in, and perform actions accordingly
func (p DummyPlugin) Restore(endpoint plugin.ShieldEndpoint) (int, error) {
	file, err := endpoint.StringValue("file")
	if err != nil {
		return plugin.PLUGIN_FAILURE, err
	}

	return plugin.Exec(plugin.STDIN, fmt.Sprintf("/bin/sh -c \"/bin/cat > %s\"", file))
}

// Called when you want to store backup data. Examine the ShieldEndpoint passed in, and perform actions accordingly
func (p DummyPlugin) Store(endpoint plugin.ShieldEndpoint) (string, int, error) {
	directory, err := endpoint.StringValue("directory")
	if err != nil {
		return "", plugin.PLUGIN_FAILURE, err
	}

	file := plugin.GenUUID()

	success, err := plugin.Exec(plugin.STDIN, fmt.Sprintf("/bin/sh -c \"/bin/cat > %s/%s\"", directory, file))
	return file, success, err
}

// Called when you want to retreive backup data. Examine the ShieldEndpoint passed in, and perform actions accordingly
func (p DummyPlugin) Retrieve(endpoint plugin.ShieldEndpoint, file string) (int, error) {
	directory, err := endpoint.StringValue("directory")
	if err != nil {
		return plugin.PLUGIN_FAILURE, err
	}

	return plugin.Exec(plugin.STDOUT, fmt.Sprintf("/bin/cat %s/%s", directory, file))
}

func (p DummyPlugin) Purge(endpoint plugin.ShieldEndpoint, key string) (int, error) {
	return plugin.UNSUPPORTED_ACTION, fmt.Errorf("I'm just a dummy plugin. I don't know how to purge")
}

//That's all there is to writing a plugin. If your plugin doesn't need to implement Store/Retireve, or Backup/Restore,
// Define the functions, and have them return plugin.UNSUPPORTED_ACTION and a useful error message.
