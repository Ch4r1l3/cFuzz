package main

import (
	"../common"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

type LibFuzzer struct {
	logger hclog.Logger
}

func (a *LibFuzzer) Fuzz() error {
	a.logger.Debug("fuzz in libfuzzer!")
	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "libfuzzer",
	MagicCookieValue: "libfuzzer",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	libfuzzer := &LibFuzzer{
		logger: logger,
	}

	var pluginMap = map[string]plugin.Plugin{
		"libfuzzer": &fuzzer.FuzzerPlugin{Impl: libfuzzer},
	}

	logger.Debug("message from plugin")
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
