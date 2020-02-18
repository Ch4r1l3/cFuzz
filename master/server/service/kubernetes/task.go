package kubernetes

import (
	"encoding/json"
	"strconv"
)

func CreateStorageItem(taskID uint64, exist bool, mtype, path, relPath string) (uint64, error) {
	var result []byte
	var err error
	if exist {
		form := map[string]interface{}{
			"type":          mtype,
			"relPath":       relPath,
			"existsInImage": exist,
			"path":          path,
		}
		result, err = requestProxyPost(uint64(taskID), []string{"storage_item", "exist"}, form)
	} else {
		form := map[string]string{
			"type":    mtype,
			"relPath": relPath,
		}
		result, err = requestProxyPostWithFile(taskID, []string{"storage_item"}, form, path)
	}
	if err != nil {
		return 0, err
	}
	var resp clientStorageItemPostResp
	if err := json.Unmarshal(result, &resp); err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func GetStorageItems(taskID uint64) ([]byte, error) {
	return requestProxyGet(taskID, []string{"storage_item"})
}

func StopTask(taskID uint64) error {
	_, err := requestProxyPost(taskID, []string{"task", "stop"}, []string{"xx"})
	return err
}

func GetTask(taskID uint64) (string, string, error) {
	result, err := requestProxyGet(taskID, []string{"task"})
	if err != nil {
		return "", "", err
	}
	var clientTask clientTaskGetResp
	if err := json.Unmarshal(result, &clientTask); err != nil {
		return "", "", err
	}
	return clientTask.Status, clientTask.ErrorMsg, nil
}

func GetCrashes(taskID uint64) ([]clientCrashGetResp, error) {
	result, err := requestProxyGet(taskID, []string{"task", "crash"})
	if err != nil {
		return nil, err
	}
	var crashes []clientCrashGetResp
	if err = json.Unmarshal(result, &crashes); err != nil {
		return nil, err
	}
	return crashes, nil
}

func DownloadCrash(taskID uint64, crashID uint64, crashesPath string) (string, error) {
	return requestProxySaveFile(taskID, []string{"task", "crash", strconv.Itoa(int(crashID))}, crashesPath)
}

func GetResult(taskID uint64) (*clientResultGetResp, error) {
	result, err := requestProxyGet(taskID, []string{"task", "result"})
	if err != nil {
		return nil, err
	}
	if len(result) > 10 {
		var fuzzResult clientResultGetResp
		if err = json.Unmarshal(result, &fuzzResult); err != nil {
			return nil, err
		}
		return &fuzzResult, nil
	}
	return nil, nil
}

func CreateTask(taskID uint64, postData map[string]interface{}) error {
	_, err := requestProxyPost(taskID, []string{"task"}, postData)
	return err
}

func StartTask(taskID uint64) error {
	_, err := requestProxyPost(taskID, []string{"task", "start"}, struct{}{})
	return err
}
