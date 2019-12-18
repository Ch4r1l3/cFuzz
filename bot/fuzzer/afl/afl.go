package main

import (
	"../common"
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

type AFL struct {
	logger hclog.Logger
}

func (a *AFL) Fuzz() error {
	a.logger.Debug("fuzz in afl!")
	return errors.New("afl!")
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  3,
	MagicCookieKey:   "afl",
	MagicCookieValue: "afl",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	afl := &AFL{
		logger: logger,
	}

	/*
		var pluginMap = map[string]plugin.Plugin{
			"afl": &fuzzer.FuzzerPlugin{Impl: afl},
		}
	*/

	logger.Debug("message from plugin")
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		VersionedPlugins: map[int]plugin.PluginSet{
			2: {
				"afl": &fuzzer.FuzzerPlugin{Impl: afl},
			},
			3: {
				"afl": &fuzzer.FuzzerGRPCPlugin{Impl: afl},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
