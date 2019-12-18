package main

import (
	"./fuzzer/common"
	"fmt"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"log"
	"os"
	"os/exec"
)

var handshakeAFLConfig = plugin.HandshakeConfig{
	ProtocolVersion:  3,
	MagicCookieKey:   "afl",
	MagicCookieValue: "afl",
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

	for i := 1; i < 10; i += 1 {
		err = afl.Fuzz()
		if err != nil {
			fmt.Println("Error:", err.Error())
		}
	}
	logger.Debug(err.Error())
}
