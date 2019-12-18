package main

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  3,
	MagicCookieKey:   "fuzz",
	MagicCookieValue: "fuzz",
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
