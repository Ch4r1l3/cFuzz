package service

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/Ch4r1l3/cFuzz/bot/server/logger"
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/hashicorp/go-plugin"
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
	for _, c := range fuzzResult.Crashes {
		if _, ok := crashCheckMap[c.InputPath]; !ok {
			crashCheckMap[c.InputPath] = true
			reproduceAble := false
			if v, ok := reproduceResult[c.InputPath]; ok && v.ReturnCode != 0 {
				reproduceAble = true
			}
			models.CreateCrash(c.InputPath, reproduceAble, c.FileName)
		}
	}
	err := models.CreateFuzzResult(fuzzResult.Command, fuzzResult.Stats, fuzzResult.TimeExecuted)
	if err != nil {
		running = false
	}
}

func Fuzz(pluginPath string, targetPath string, corpusDir string, maxTime int, fuzzMaxTime int, arguments map[string]string, environments []string) {
	mutex.Lock()
	defer mutex.Unlock()
	if !running {
		running = true
		go func() {
			var Err error
			defer func() {
				mutex.Lock()
				running = false
				if Err != nil {
					models.DB.Model(&models.Task{}).Update("Status", models.TaskError)
					logger.Logger.Debug("error is !!!!!:" + Err.Error())
				}
				mutex.Unlock()
			}()

			plugins := map[int]plugin.PluginSet{
				1: {
					"fuzzer": &fuzzer.FuzzerPlugin{},
				},
				2: {
					"fuzzer": &fuzzer.FuzzerGRPCPlugin{},
				},
			}
			fuzzerClient := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig:  handshakeConfig,
				VersionedPlugins: plugins,
				Cmd:              exec.Command(pluginPath),
				Logger:           logger.Logger,
				AllowedProtocols: []plugin.Protocol{
					plugin.ProtocolNetRPC, plugin.ProtocolGRPC,
				},
			})
			defer fuzzerClient.Kill()

			fuzzerRpcClient, err := fuzzerClient.Client()
			if err != nil {
				logger.Logger.Debug(err.Error())
				Err = err
				return
			}

			fuzzerRaw, err := fuzzerRpcClient.Dispense("fuzzer")
			if err != nil {
				logger.Logger.Debug(err.Error())
				Err = err
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
				logger.Logger.Debug(err.Error())
				Err = err
				return
			}

			fuzzArg := fuzzer.FuzzArg{
				MaxTime: fuzzMaxTime,
			}
			fuzzChan := make(chan struct{})
			ticker := time.NewTicker(time.Duration(fuzzMaxTime) * time.Second)
			defer ticker.Stop()
			go func() {
				for {
					fuzzChan <- struct{}{}
					<-ticker.C
				}
			}()

			for {
				select {

				case <-time.After(time.Duration(maxTime) * time.Second):
					return

				case <-controlChan:
					return

				case <-fuzzChan:
					fuzzResult, err := fuzzerPlugin.Fuzz(fuzzArg)
					if err != nil {
						Err = err
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
								logger.Logger.Debug(err.Error())
							} else {
								reproduceResult[c.InputPath] = tresult
							}
							logger.Logger.Debug("find crash: " + c.InputPath)
						}
					}
					handleFuzzResult(fuzzResult, reproduceResult)
					//go handleFuzzResult(fuzzResult, reproduceResult)

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
