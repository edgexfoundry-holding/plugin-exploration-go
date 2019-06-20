//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package foo

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bogus/my-abstraction/pkg/types"
	"github.com/hashicorp/go-plugin"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "FOO_PLUGIN",
	MagicCookieValue: "foo",
}

// PluginMap is the map of hashicorp we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"foorpc":  &FooRPCPlugin{},
	"foogrpc": &FooGRPCPlugin{},
}

// Required to kill plugin on call to Close()
var pluginClient *plugin.Client

func NewFooPlugin(config types.Config) (Client, error) {

	pluginClient = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             exec.Command(config.Plugin),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	// Can't do this here since this is a factory method. The Close interface replaces this.
	// defer pluginClient.Kill()

	// Connect via RPC
	rpcClient, err := pluginClient.Client()
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(config.PluginName)
	if err != nil {
		fmt.Println("Error dispensing:", err.Error())
		os.Exit(1)
	}

	foo := raw.(Client)

	return foo, nil
}
