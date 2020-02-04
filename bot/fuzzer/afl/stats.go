package main

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func ParseFuzzerStatsFile(outDir string) (map[string]string, error) {
	result := make(map[string]string)
	content, err := ioutil.ReadFile(path.Join(outDir, FUZZER_STATS_FILE))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	for _, v := range lines {
		values := strings.Split(v, ":")
		if len(values) >= 2 {
			key := strings.TrimSpace(values[0])
			value := strings.TrimSpace(values[1])
			result[key] = value
		}
	}
	return result, nil
}

func (a *AFL) GetAllCrashes(outDir string) ([]fuzzer.Crash, error) {
	a.logger.Debug("afl in get all crash")
	crashPath := path.Join(outDir, CRASH_PATH)
	crashStorePath := path.Join(outDir, CRASH_STORE_PATH)
	os.MkdirAll(crashStorePath, os.ModePerm)
	files, err := ioutil.ReadDir(crashPath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, nil
	}
	crashes := []fuzzer.Crash{}
	for _, v := range files {
		if v.Name() != "README.txt" {
			crashName := fmt.Sprintf("crash%d", a.crashNum)
			a.crashNum += 1
			os.Link(path.Join(crashPath, v.Name()), path.Join(crashStorePath, crashName))
			crashes = append(crashes, fuzzer.Crash{
				InputPath: path.Join(crashStorePath, crashName),
				FileName:  v.Name(),
			})
		}
	}
	return crashes, nil
}
