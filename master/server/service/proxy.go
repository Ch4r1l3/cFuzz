package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"io"
	"k8s.io/client-go/rest"
	"mime/multipart"
	"os"
	"path/filepath"
)

func requestProxyGet(taskID uint64, url []string) (map[string]interface{}, error, int) {
	urls := append([]string{"proxy"}, url...)
	result := ClientSet.
		CoreV1().
		RESTClient().
		Get().
		Namespace(config.KubernetesConf.Namespace).
		Resource("services").
		Name(fmt.Sprintf(ServiceNameFmt, taskID)).
		Suffix(urls...).Do()

	var statusCode int
	result.StatusCode(&statusCode)
	if statusCode == 0 {
		return nil, result.Error(), 0
	}
	bytesData, _ := result.Raw()
	var temp map[string]interface{}
	if err := json.Unmarshal(bytesData, &temp); err != nil {
		return nil, err, 0
	}
	return temp, nil, statusCode
}

func requestProxyPost(taskID uint64, url []string, data map[string]interface{}) (map[string]interface{}, error, int) {
	return requestProxyPostPut("Post", taskID, url, data)
}

func requestProxyPut(taskID uint64, url []string, data map[string]interface{}) (map[string]interface{}, error, int) {
	return requestProxyPostPut("Put", taskID, url, data)
}

func requestProxyPostPut(method string, taskID uint64, url []string, data map[string]interface{}) (map[string]interface{}, error, int) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err, 0
	}
	client := ClientSet.CoreV1().RESTClient()
	var request *rest.Request
	if method == "Put" {
		request = client.Put()
	} else if method == "Post" {
		request = client.Post()
	} else {
		request = client.Post()
	}
	urls := append([]string{"proxy"}, url...)
	result := request.
		Namespace(config.KubernetesConf.Namespace).
		Resource("services").
		Name(fmt.Sprintf(ServiceNameFmt, taskID)).
		Suffix(urls...).
		Body(bytes).
		SetHeader("Content-Type", "application/json").Do()
	var statusCode int
	result.StatusCode(&statusCode)
	if statusCode == 0 {
		return nil, result.Error(), 0
	}
	bytesData, _ := result.Raw()
	var temp map[string]interface{}
	if err := json.Unmarshal(bytesData, &temp); err != nil {
		return nil, err, 0
	}
	return temp, nil, statusCode
}

func requestProxyPostWithFile(taskID uint64, url []string, form map[string]string, filePath string) (map[string]interface{}, error, int) {
	return requestProxyPostPutWithFile("Post", taskID, url, form, filePath)
}

func requestProxyPutWithFile(taskID uint64, url []string, form map[string]string, filePath string) (map[string]interface{}, error, int) {
	return requestProxyPostPutWithFile("Put", taskID, url, form, filePath)
}

func requestProxyPostPutWithFile(method string, taskID uint64, url []string, form map[string]string, filePath string) (map[string]interface{}, error, int) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err, 0
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err, 0
	}
	io.Copy(part, file)
	for k, v := range form {
		writer.WriteField(k, v)
	}
	writer.Close()
	client := ClientSet.CoreV1().RESTClient()
	var request *rest.Request
	if method == "Put" {
		request = client.Put()
	} else if method == "Post" {
		request = client.Post()
	} else {
		request = client.Post()
	}
	urls := append([]string{"proxy"}, url...)
	result := request.
		Namespace(config.KubernetesConf.Namespace).
		Resource("services").
		Name(fmt.Sprintf(ServiceNameFmt, taskID)).
		Suffix(urls...).
		Body(body.Bytes()).
		SetHeader("Content-Type", writer.FormDataContentType()).
		Do()
	var statusCode int
	result.StatusCode(&statusCode)
	if statusCode == 0 {
		return nil, result.Error(), 0
	}
	var temp map[string]interface{}
	bytesData, _ := result.Raw()
	if err := json.Unmarshal(bytesData, &temp); err != nil {
		return nil, err, 0
	}
	return temp, nil, statusCode
}
