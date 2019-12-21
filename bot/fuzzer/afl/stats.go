package main

import (
	"github.com/Ch4r1l3/cFuzz/bot/fuzzer/common"
	"io/ioutil"
	"path"
	"strings"
)

func ParseFuzzerStatsFile(outDir string) (map[string]string, error) {
	result := make(map[string]string)
	content, err := ioutil.ReadFile(path.Join(outDir, FUZZERSTATSFILE))
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

func GetAllCrashes(outDir string) ([]fuzzer.Crash, error) {
	crashPath := path.Join(outDir, CRASHPATH)
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
			crashes = append(crashes, fuzzer.Crash{
				InputPath: path.Join(crashPath, v.Name()),
			})
		}
	}
	return crashes, nil
}
