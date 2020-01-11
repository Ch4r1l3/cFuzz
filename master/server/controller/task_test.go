package controller

import (
	"github.com/gavv/httpexpect"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestTaskList(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/task").Expect().Status(http.StatusOK)
}

func TestTask1(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	deploymentPostData := map[string]interface{}{
		"name":    "test",
		"content": "11123",
	}
	deploymentID := int(e.POST("/deployment").WithJSON(deploymentPostData).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	err := ioutil.WriteFile("./fuzzer_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer_test")
	}()

	fuzzerID := int(e.POST("/fuzzer").WithMultipart().WithFile("file", "fuzzer_test").WithFormField("name", "afl").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData := map[string]interface{}{
		"name":         "test",
		"deploymentid": deploymentID,
		"time":         100,
		"fuzzerid":     fuzzerID,
		"environments": []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Array().First().Object()
	obj.Keys().ContainsOnly("id", "deploymentid", "time", "fuzzerid", "running", "environments", "arguments", "image", "name")
	obj.Value("id").NotEqual(0)
	obj.Value("deploymentid").NotEqual(0)
	obj.Value("time").NotEqual(0)
	obj.Value("fuzzerid").NotEqual(0)
	obj.Value("running").Equal(false)
	obj.Value("environments").Array().Elements("123", "2333")
	obj.Value("arguments").Object().Value("a1").Equal("a2")
	obj.Value("arguments").Object().Value("a2").Equal("a3")
	e.DELETE("/deployment/" + strconv.Itoa(deploymentID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/fuzzer/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}

func TestTask2(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	deploymentPostData := map[string]interface{}{
		"name":    "test",
		"content": "11123",
	}
	deploymentID := int(e.POST("/deployment").WithJSON(deploymentPostData).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	err := ioutil.WriteFile("./fuzzer_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer_test")
	}()

	fuzzerID := int(e.POST("/fuzzer").WithMultipart().WithFile("file", "fuzzer_test").WithFormField("name", "afl").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":         "test",
		"deploymentid": deploymentID,
		"time":         100,
		"fuzzerid":     fuzzerID,
		"environments": []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}
	taskPostData2 := map[string]interface{}{
		"deploymentid": -1,
	}
	taskPostData3 := map[string]interface{}{
		"fuzzerid": -1,
	}
	taskPostData4 := map[string]interface{}{
		"time": -1,
	}
	taskPostData5 := map[string]interface{}{
		"environments": []string{"2", "3"},
	}
	taskPostData6 := map[string]interface{}{
		"arguments": map[string]string{
			"a3": "a4",
			"a4": "a5",
		},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData1).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData2).Expect().Status(http.StatusBadRequest)
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData3).Expect().Status(http.StatusBadRequest)
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData4).Expect().Status(http.StatusBadRequest)
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData5).Expect().Status(http.StatusOK)
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData6).Expect().Status(http.StatusOK)

	e.GET("/task").Expect().Status(http.StatusOK).JSON().Array().First().Object().Value("environments").Array().Elements("2", "3")

	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Array().First().Object().Value("arguments").Object()
	obj.Value("a3").Equal("a4")
	obj.Value("a4").Equal("a5")

	e.DELETE("/deployment/" + strconv.Itoa(deploymentID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/fuzzer/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}

func TestTask3(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	err := ioutil.WriteFile("./fuzzer_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer_test")
	}()

	fuzzerID := int(e.POST("/fuzzer").WithMultipart().WithFile("file", "fuzzer_test").WithFormField("name", "afl").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":         "test",
		"image":        "registry.cn-hangzhou.aliyuncs.com/cfuzz/test:v1",
		"time":         100,
		"fuzzerid":     fuzzerID,
		"environments": []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}
	taskPostData2 := map[string]interface{}{
		"running": true,
	}
	taskPostData3 := map[string]interface{}{
		"running": false,
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData1).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData2).Expect().Status(http.StatusOK)
	<-time.After(time.Duration(15) * time.Second)
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData3).Expect().Status(http.StatusOK)
	<-time.After(time.Duration(6) * time.Second)
	e.DELETE("/fuzzer/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}
