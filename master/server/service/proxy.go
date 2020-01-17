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
	"time"
)

func bytes2interface(data []byte) (interface{}, error) {
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}
	return temp, nil
}

func requestProxyGet(taskID uint64, url []string) ([]byte, error) {
	urls := append([]string{"proxy"}, url...)
	result := ClientSet.
		CoreV1().
		RESTClient().
		Get().
		Namespace(config.KubernetesConf.Namespace).
		Resource("services").
		Name(fmt.Sprintf(ServiceNameFmt, taskID)).
		Timeout(time.Duration(config.KubernetesConf.RequestTimeout)).
		Suffix(urls...).Do()

	bytesData, _ := result.Raw()
	err := result.Error()
	return bytesData, err
}

func requestProxyPost(taskID uint64, url []string, data interface{}) ([]byte, error) {
	return requestProxyPostPut("Post", taskID, url, data)
}

func requestProxyPut(taskID uint64, url []string, data interface{}) ([]byte, error) {
	return requestProxyPostPut("Put", taskID, url, data)
}

func requestProxyPostRaw(taskID uint64, url []string, data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	client := ClientSet.CoreV1().RESTClient()
	var request *rest.Request
	request = client.Post()
	urls := append([]string{"proxy"}, url...)
	return request.
		Namespace(config.KubernetesConf.Namespace).
		Resource("services").
		Name(fmt.Sprintf(ServiceNameFmt, taskID)).
		Suffix(urls...).
		Body(bytes).
		Timeout(time.Duration(config.KubernetesConf.RequestTimeout)).
		SetHeader("Content-Type", "application/json").DoRaw()
}

func requestProxyPostPut(method string, taskID uint64, url []string, data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
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
		Timeout(time.Duration(config.KubernetesConf.RequestTimeout)).
		SetHeader("Content-Type", "application/json").Do()
	bytesData, _ := result.Raw()
	err = result.Error()
	return bytesData, err
}

func requestProxyPostWithFile(taskID uint64, url []string, form map[string]string, filePath string) ([]byte, error) {
	return requestProxyPostPutWithFile("Post", taskID, url, form, filePath)
}

func requestProxyPutWithFile(taskID uint64, url []string, form map[string]string, filePath string) ([]byte, error) {
	return requestProxyPostPutWithFile("Put", taskID, url, form, filePath)
}

func requestProxyPostPutWithFile(method string, taskID uint64, url []string, form map[string]string, filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
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
		Timeout(time.Duration(config.KubernetesConf.RequestTimeout)).
		Do()
	bytesData, _ := result.Raw()
	err = result.Error()
	return bytesData, err
}
