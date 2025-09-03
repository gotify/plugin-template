package myplugin

import "github.com/gotify/plugin-api/v2/generated/protobuf"

var (
	PluginInfo = protobuf.Info{
		ModulePath:  "github.com/gotify/plugin-template",
		Version:     "1.0.0",
		Author:      "Your Name",
		Website:     "https://gotify.net/docs/plugin",
		Description: "An example plugin with travis-ci building",
		License:     "MIT",
		Name:        "gotify/plugin-template",
	}

	Capabilities = []protobuf.Capability{
		protobuf.Capability_DISPLAYER,
		protobuf.Capability_CONFIGURER,
		protobuf.Capability_STORAGER,
		protobuf.Capability_WEBHOOKER,
		protobuf.Capability_MESSENGER,
	}
)
