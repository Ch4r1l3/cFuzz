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
		TargetPath: "aa",
		MaxTime:    1,
	}
	tPrepareArg := fuzzer.PrepareArg{
		CorpusDir:   "123",
		TargetPath:  "aa",
		Arguments:   map[string]string{},
		Enviroments: map[string]string{},
	}
	for i := 1; i < 10; i += 1 {
		err = afl.Fuzz(tFuzzArg)
		if err != nil {
			fmt.Println("Error:", err.Error())
		}
		err = afl.Prepare(tPrepareArg)
		if err != nil {
			fmt.Println("Error:", err.Error())
		}
	}
}
