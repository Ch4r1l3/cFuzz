package main

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"log"
	"os"
	"os/exec"
)

var handshakeAFLConfig = plugin.HandshakeConfig{
	ProtocolVersion:  3,
	MagicCookieKey:   "fuzz",
	MagicCookieValue: "fuzz",
}

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	plugins := map[int]plugin.PluginSet{}
	plugins[3] = plugin.PluginSet{
		"afl": &fuzzer.FuzzerGRPCPlugin{},
	}

	aflclient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeAFLConfig,
		VersionedPlugins: plugins,
		Cmd:              exec.Command("./fuzzer/afl/afl"),
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer aflclient.Kill()

	aflrpcClient, err := aflclient.Client()
	if err != nil {
		log.Fatal(err)
	}

	aflraw, err := aflrpcClient.Dispense("afl")
	fmt.Print("test!")
	if err != nil {
		log.Fatal(err)
	}
	afl := aflraw.(fuzzer.Fuzzer)

	tFuzzArg := fuzzer.FuzzArg{
		TargetPath: "/tmp/test/test",
		MaxTime:    60,
	}
	tPrepareArg := fuzzer.PrepareArg{
		CorpusDir:   "/tmp/test/corpus",
		TargetPath:  "/tmp/test/test",
		Arguments:   map[string]string{},
		Enviroments: []string{},
	}
	err = afl.Prepare(tPrepareArg)
	if err != nil {
		log.Fatal(err)
	}

	err = afl.Fuzz(tFuzzArg)
	if err != nil {
		log.Fatal(err)
	}
}
