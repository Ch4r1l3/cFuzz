package main

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"io/ioutil"
	"path"
	"strings"
)

func ParseFuzzerStats(output []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, v := range output {
		if STATS_REGEX.MatchString(v) {
			stat := v[6:]
			kv := strings.Split(stat, ":")
			if len(kv) < 2 {
				continue
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			result[key] = value
		}
	}
	return result, nil
}

func (a *LibFuzzer) GetAllCrashes() ([]fuzzer.Crash, error) {
	a.logger.Debug("libfuzzer in get all crash")
	files, err := ioutil.ReadDir(a.outputDir)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, nil
	}
	crashes := []fuzzer.Crash{}
	for _, v := range files {
		crashes = append(crashes, fuzzer.Crash{
			InputPath: path.Join(a.outputDir, v.Name()),
			FileName:  v.Name(),
		})
	}
	return crashes, nil
}
