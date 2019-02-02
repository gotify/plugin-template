package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gotify/plugin-api"
)

// GetGotifyPluginInfo returns gotify plugin info.
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:  "github.com/gotify/plugin-template",
		Version:     "1.0.0",
		Author:      "Your Name",
		Website:     "https://gotify.net/docs/plugin",
		Description: "An example plugin with travis-ci building",
		License:     "MIT",
		Name:        "gotify/plugin-template",
	}
}

// MyPlugin is the gotify plugin instance.
type MyPlugin struct {
}

// Enable enables the plugin.
func (c *MyPlugin) Enable() error {
	return nil
}

// Disable disables the plugin.
func (c *MyPlugin) Disable() error {
	return nil
}

// RegisterWebhook implements plugin.Webhooker.
func (c *MyPlugin) RegisterWebhook(basePath string, g *gin.RouterGroup) {
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &MyPlugin{}
}

func main() {
	panic("this should be built as go plugin")
}
