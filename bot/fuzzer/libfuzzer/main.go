package main

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "fuzz",
	MagicCookieValue: "fuzz",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	libFuzzer := &LibFuzzer{
		logger: logger,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: {
				"fuzzer": &fuzzer.FuzzerPlugin{Impl: libFuzzer},
			},
			2: {
				"fuzzer": &fuzzer.FuzzerGRPCPlugin{Impl: libFuzzer},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
