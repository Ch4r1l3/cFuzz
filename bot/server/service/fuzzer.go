package service

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
	"os/exec"
	"time"
)

func IsRunning() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return running
}

func handleFuzzResult(fuzzResult fuzzer.FuzzResult, reproduceResult map[string]fuzzer.ReproduceResult) {
	mutex.Lock()
	defer mutex.Unlock()
	if running {
		for _, c := range fuzzResult.Crashes {
			if _, ok := crashCheckMap[c.InputPath]; !ok {
				crashCheckMap[c.InputPath] = true
				reproduceAble := false
				if v, ok := reproduceResult[c.InputPath]; ok && v.ReturnCode != 0 {
					reproduceAble = true

				}
				models.CreateCrash(c.InputPath, reproduceAble)

			}
		}
		err := models.CreateFuzzResult(fuzzResult.Command, fuzzResult.Stats, fuzzResult.TimeExecuted)
		if err != nil {
			running = false
		}
	}
}

func Fuzz(pluginPath string, targetPath string, corpusDir string, maxTime int, fuzzMaxTime int, arguments map[string]string, environments []string) {
	mutex.Lock()
	defer mutex.Unlock()
	if !running {
		running = true
		go func() {
			var err error
			logger := hclog.New(&hclog.LoggerOptions{
				Name:   "plugin",
				Output: os.Stdout,
				Level:  hclog.Debug,
			})
			defer func() {
				mutex.Lock()
				running = false
				if err != nil {
					models.DB.Model(&models.Task{}).Update("Status", config.TASK_ERROR)
				}
				logger.Debug(err.Error())
				mutex.Unlock()
			}()

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
			if err != nil {
				logger.Debug(err.Error())
				return
			}

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
					reproduceResult := map[string]fuzzer.ReproduceResult{}
					for _, c := range fuzzResult.Crashes {
						if _, ok := crashCheckMap[c.InputPath]; !ok {
							targ := fuzzer.ReproduceArg{
								InputPath: c.InputPath,
								MaxTime:   config.ServerConf.DefaultReproduceTime,
							}
							tresult, err := fuzzerPlugin.Reproduce(targ)
							if err != nil {
								logger.Debug(err.Error())
							} else {
								reproduceResult[c.InputPath] = tresult
							}
						}
					}
					go handleFuzzResult(fuzzResult, reproduceResult)

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
