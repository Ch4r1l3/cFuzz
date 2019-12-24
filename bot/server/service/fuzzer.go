package service

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
	"os/exec"
)

func IsRunning() bool {
	mutex.Lock()
	defer mutex.Unlock()
	tmpRunning := running
	return tmpRunning
}

func Fuzz(pluginPath string, targetPath string, corpusDir string, maxTime int, fuzzMaxTime int, arguments map[string]string, environments []string) {
	mutex.Lock()
	defer mutex.Unlock()
	if !running {
		running = true
		go func() {
			logger := hclog.New(&hclog.LoggerOptions{
				Name:   "plugin",
				Output: os.Stdout,
				Level:  hclog.Debug,
			})
			fuzzerClient := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: handshakeConfig,
				Cmd:             exec.Command(pluginPath),
				Logger:          logger,
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC, plugin.ProtocolGRPC,
				},
			})
			defer fuzzerClient.Kill()

			fuzzerRpcClient, err := fuzzerClient.Client()
			if err != nil {
				logger.Debug(err.Error())
				mutex.Lock()
				running = true
				mutex.Unlock()
			}

			fuzzerRaw, err := fuzzerRpcClient.Dispense("fuzzer")
			if err != nil {
				logger.Debug(err.Error())
				mutex.Lock()
				running = true
				mutex.Unlock()

			}

			fuzzerPlugin := fuzzerRaw.(fuzzer.Fuzzer)

			prepareArg := fuzzer.PrepareArg{
				CorpusDir:    corpusDir,
				TargetPath:   targetPath,
				Arguments:    arguments,
				Environments: environments,
			}

			err = fuzzerPlugin.Prepare(prepareArg)

			fuzzArg := fuzzer.FuzzArg{
				MaxTime: fuzzMaxTime,
			}

			_, err = fuzzerPlugin.Fuzz(fuzzArg)
			if err != nil {

			}
		}()
	}
}

func StopFuzz() {
	mutex.Lock()
	defer mutex.Unlock()
	if running {
		controlChan <- struct{}{}
		running = false
	}
}
