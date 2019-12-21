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
	ProtocolVersion:  2,
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
	//plugins[3] = plugin.PluginSet{
	//	"afl": &fuzzer.FuzzerGRPCPlugin{},
	//}
	plugins[2] = plugin.PluginSet{
		"afl": &fuzzer.FuzzerPlugin{},
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
		MaxTime: 60,
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

	v, err := afl.Fuzz(tFuzzArg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)
	fmt.Println("Stats !")
	for k, v := range v.Stats {
		fmt.Println(k + ": " + v)
	}
	fmt.Println("Crashes !")

	if len(v.Crashes) > 0 {
		for _, v := range v.Crashes {
			tReproduceArg := fuzzer.ReproduceArg{
				InputPath: v.InputPath,
				MaxTime:   60,
			}
			result, err := afl.Reproduce(tReproduceArg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Return Code:")
			fmt.Println(result.ReturnCode)
			fmt.Println("TimeExecuted:")

			fmt.Println(result.TimeExecuted)
		}
	}

	err = afl.Clean()
	if err != nil {
		log.Fatal(err)
	}
}
