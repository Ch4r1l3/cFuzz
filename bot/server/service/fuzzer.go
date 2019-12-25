package service

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	//"github.com/Ch4r1l3/cFuzz/bot/server/config"
	//"github.com/Ch4r1l3/cFuzz/bot/server/models"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
	"os/exec"
	"time"
)

func IsRunning() bool {
	mutex.Lock()
	defer mutex.Unlock()
	tmpRunning := running
	return tmpRunning
}

func handleFuzzResult(result fuzzer.FuzzResult) {

}

func Fuzz(pluginPath string, targetPath string, corpusDir string, maxTime int, fuzzMaxTime int, arguments map[string]string, environments []string) {
	mutex.Lock()
	defer mutex.Unlock()
	if !running {
		running = true
		go func() {
			defer func() {
				mutex.Lock()
				running = false
				mutex.Unlock()
			}()
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
				return
			}

			fuzzerRaw, err := fuzzerRpcClient.Dispense("fuzzer")
			if err != nil {
				logger.Debug(err.Error())
				return
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

			for {
				select {

				case <-time.After(time.Duration(maxTime) * time.Second):
					break

				case <-controlChan:
					break

				default:
					fuzzResult, err := fuzzerPlugin.Fuzz(fuzzArg)
					if err != nil {
						logger.Debug(err.Error())
						return
					}
					go handleFuzzResult(fuzzResult)

				}
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
