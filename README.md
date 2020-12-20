# gotify/plugin-template [![](https://travis-ci.org/gotify/plugin-template.svg?branch=master)](https://travis-ci.org/gotify/plugin-template)

A plugin template for [gotify/server](https://github.com/gotify/server)
using [gotify/plugin-api](https://github.com/gotify/plugin-api).

## Getting Started

1. Clone, fork or copy this repository.
1. Change `PLUGIN_NAME` in [Makefile](Makefile).
1. Setup building on every release
   * Enable travis-ci in your repository.
   * Add `GH_TOKEN` environment variable see [travis-ci docs](https://docs.travis-ci.com/user/deployment/pages/#setting-the-github-token).
1. Implement your plugin. See [plugin docs](https://gotify.net/docs/plugin).
1. Create a release to automatically build the plugin.

*When you're done, feel free to add your plugin to [gotify/contrib](https://github.com/gotify/contrib).*

## Building

For building the plugin gotify/build docker images are used to ensure compatibility with 
[gotify/server](https://github.com/gotify/server).

`GOTIFY_VERSION` can be a tag, commit or branch from the gotify/server repository.

This command builds the plugin for amd64, arm-7 and arm64. 
The resulting shared object will be compatible with gotify/server version 2.0.20.
```bash
$ make GOTIFY_VERSION="v2.0.20" FILE_SUFFIX="for-gotify-v2.0.20" build
```
