// The `swift` plugin for SHIELD is intended to be a back-end storage
// plugin, wrapping OpenStack Swift.
//
// PLUGIN FEATURES
//
// This plugin implements functionality suitable for use with the following
// SHIELD Job components:
//
//  Target: no
//  Store:  yes
//
// PLUGIN CONFIGURATION
//
// The endpoint configuration passed to this plugin is used to determine
// how to connect to S3, and where to place/retrieve the data once connected.
// your endpoint JSON should look something like this:
//
//    {
//        "auth_url":      "host",
//        "project_name":  "openstack-project",
//        "username":      "your-username",
//        "password":      "secret-access-key",
//        "container":     "bucket-name",
//        "prefix":        "/path/inside/bucket/to/place/backup/data",
//        "debug":         false
//    }
//
// Default Configuration
//
//    {
//        "prefix" : "",
//        "debug"  : false
//    }
//
// STORE DETAILS
//
// When storing data, this plugin connects to the Swift service, and uploads the data
// into the specified container, using a path/filename with the following format:
//
//    <prefix>/<YYYY>/<MM>/<DD>/<HH-mm-SS>-<UUID>
//
// Upon successful storage, the plugin then returns this filename to SHIELD to use
// as the `store_key` when the data needs to be retrieved, or purged.
//
// RETRIEVE DETAILS
//
// When retrieving data, this plugin connects to the Swift service, and retrieves the data
// located in the specified container, identified by the `store_key` provided by SHIELD.
//
// PURGE DETAILS
//
// When purging data, this plugin connects to the Swift service, and deletes the data
// located in the specified container, identified by the `store_key` provided by SHIELD.
//
// DEPENDENCIES
//
// None.
//
package main

// https://github.com/openstack/golang-client/blob/master/examples/objectstorage/objectstorage.go

import (
	"fmt"
	"strings"
	"time"

	"github.com/starkandwayne/shield/plugin"
)

const (
	DefaultDebug  = false
	DefaultPrefix = ""
)

func main() {
	p := SwiftPlugin{
		Name:    "Openstack Swift Backup + Storage Plugin",
		Author:  "Stark & Wayne",
		Version: "0.0.1",
		Features: plugin.PluginFeatures{
			Target: "no",
			Store:  "yes",
		},
		Example: `
{
  "auth_url":      "host",
  "project_name":  "openstack-project",
  "username":      "your-username",
  "password":      "secret-access-key",
  "container":     "bucket-name",
  "prefix":        "/path/inside/bucket/to/place/backup/data",
  "debug":         false
}
`,
		Defaults: `
{
  "prefix":        "",
  "debug":         false
}
`,
	}

	plugin.Run(p)
}

type SwiftPlugin plugin.PluginInfo

type SwiftConnectionInfo struct {
	AuthURL     string
	ProjectName string
	Username    string
	Password    string
	Container   string
	PathPrefix  string
	Debug       bool
}

func (p SwiftPlugin) Meta() plugin.PluginInfo {
	return plugin.PluginInfo(p)
}

func (p SwiftPlugin) Validate(endpoint plugin.ShieldEndpoint) error {
	return plugin.UNIMPLEMENTED
}
func (p SwiftPlugin) Backup(endpoint plugin.ShieldEndpoint) error {
	return plugin.UNIMPLEMENTED
}

func (p SwiftPlugin) Restore(endpoint plugin.ShieldEndpoint) error {
	return plugin.UNIMPLEMENTED
}

func (p SwiftPlugin) Store(endpoint plugin.ShieldEndpoint) (string, error) {
	return "", plugin.UNIMPLEMENTED
}

func (p SwiftPlugin) Retrieve(endpoint plugin.ShieldEndpoint, file string) error {
	return plugin.UNIMPLEMENTED
}

func (p SwiftPlugin) Purge(endpoint plugin.ShieldEndpoint, file string) error {
	return plugin.UNIMPLEMENTED
}

func getConnInfo(e plugin.ShieldEndpoint) (info *SwiftConnectionInfo, err error) {
	info = &SwiftConnectionInfo{}
	info.AuthURL, err = e.StringValue("auth_url")
	if err != nil {
		return
	}

	info.ProjectName, err = e.StringValue("project_name")
	if err != nil {
		return
	}

	info.Username, err = e.StringValue("username")
	if err != nil {
		return
	}

	info.Password, err = e.StringValue("password")
	if err != nil {
		return
	}

	info.Container, err = e.StringValue("container")
	if err != nil {
		return
	}

	info.PathPrefix, err = e.StringValueDefault("prefix", DefaultPrefix)
	if err != nil {
		return
	}
	info.PathPrefix = strings.TrimLeft(info.PathPrefix, "/")

	info.Debug, err = e.BooleanValueDefault("debug", DefaultDebug)
	if err != nil {
		return
	}

	return
}

func (info SwiftConnectionInfo) genBackupPath() string {
	t := time.Now()
	year, mon, day := t.Date()
	hour, min, sec := t.Clock()
	uuid := plugin.GenUUID()
	path := fmt.Sprintf("%s/%04d/%02d/%02d/%04d-%02d-%02d-%02d%02d%02d-%s", info.PathPrefix, year, mon, day, year, mon, day, hour, min, sec, uuid)
	// Remove double slashes
	path = strings.Replace(path, "//", "/", -1)
	return path
}
