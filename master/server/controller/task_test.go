package controller

import (
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
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

	fuzzerID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "fuzzer_test").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData := map[string]interface{}{
		"name":          "test",
		"deploymentid":  deploymentID,
		"time":          100,
		"fuzzCycleTime": 60,
		"fuzzerid":      fuzzerID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Array().First().Object()
	obj.Keys().ContainsOnly("id", "deploymentid", "time", "fuzzerid", "corpusid", "targetid", "status", "errorMsg", "environments", "arguments", "image", "name", "fuzzCycleTime", "startedAt")
	obj.Value("id").NotEqual(0)
	obj.Value("deploymentid").NotEqual(0)
	obj.Value("time").NotEqual(0)
	obj.Value("fuzzCycleTime").NotEqual(0)
	obj.Value("fuzzerid").NotEqual(0)
	obj.Value("environments").Array().Elements("123", "2333")
	obj.Value("arguments").Object().Value("a1").Equal("a2")
	obj.Value("arguments").Object().Value("a2").Equal("a3")
	obj.Value("status").NotEqual("")
	obj.Value("startedAt").Equal(0)
	e.DELETE("/deployment/" + strconv.Itoa(deploymentID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
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

	fuzzerID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "fuzzer_test").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"deploymentid":  deploymentID,
		"time":          100,
		"fuzzCycleTime": 60,
		"fuzzerid":      fuzzerID,
		"environments":  []string{"123", "2333"},
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
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData5).Expect().Status(http.StatusNoContent)
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData6).Expect().Status(http.StatusNoContent)

	e.GET("/task").Expect().Status(http.StatusOK).JSON().Array().First().Object().Value("environments").Array().Elements("2", "3")

	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Array().First().Object().Value("arguments").Object()
	obj.Value("a3").Equal("a4")
	obj.Value("a4").Equal("a5")

	e.DELETE("/deployment/" + strconv.Itoa(deploymentID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}

func TestTask3(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	fuzzerID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"image":         "ch4r1l3/cfuzz:test-bot",
		"time":          config.KubernetesConf.CheckTaskTime * 5,
		"fuzzCycleTime": 60,
		"fuzzerid":      fuzzerID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData1).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	e.POST("/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusBadRequest)
	targetID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	e.POST("/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusBadRequest)
	corpusID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	taskPostData4 := map[string]interface{}{
		"targetid": targetID,
		"corpusid": corpusID,
	}
	e.PUT("/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData4).Expect().Status(http.StatusNoContent)

	e.POST("/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*4) * time.Second)
	e.GET("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskRunning)
	e.POST("/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(5) * time.Second)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}

func TestTask4(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	fuzzerID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"image":         "ch4r1l3/cfuzz:test-bot",
		"time":          config.KubernetesConf.CheckTaskTime * 3,
		"fuzzCycleTime": 60,
		"fuzzerid":      fuzzerID,
		"targetid":      targetID,
		"corpusid":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData1).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	e.POST("/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*5) * time.Second)
	e.GET("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskStopped)
	e.POST("/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusBadRequest)
	<-time.After(time.Duration(5) * time.Second)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}

func TestTask5(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	fuzzerID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"image":         "ch4r1l3/cfuzz:test-bot",
		"time":          config.KubernetesConf.CheckTaskTime * 8,
		"fuzzCycleTime": 60,
		"fuzzerid":      fuzzerID,
		"targetid":      targetID,
		"corpusid":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData1).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	e.POST("/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*7) * time.Second)
	e.GET("/task/" + strconv.Itoa(taskID) + "/crash").Expect().Status(http.StatusOK).JSON().Array().Length().NotEqual(0)
	obj := e.GET("/task/" + strconv.Itoa(taskID) + "/result").Expect().Status(http.StatusOK).JSON().Object()
	obj.Keys().ContainsOnly("command", "timeExecuted", "updateAt", "stats")
	obj.Value("command").NotEqual("")
	obj.Value("timeExecuted").NotEqual(0)
	obj.Value("updateAt").NotEqual(0)
	obj.Value("stats").NotNull()
	e.POST("/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(5) * time.Second)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}

func TestTask6(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	fuzzerID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"image":         "deadbeef:v1",
		"time":          config.KubernetesConf.CheckTaskTime * 8,
		"fuzzCycleTime": 60,
		"fuzzerid":      fuzzerID,
		"targetid":      targetID,
		"corpusid":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData1).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	e.POST("/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*3) * time.Second)
	e.GET("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskError)
	e.GET("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("errorMsg").Equal("failed to create deployment")
	<-time.After(time.Duration(5) * time.Second)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}

func TestTask7(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	fuzzerID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"image":         "nginx:1.16-alpine",
		"time":          config.KubernetesConf.CheckTaskTime * 8,
		"fuzzCycleTime": 60,
		"fuzzerid":      fuzzerID,
		"targetid":      targetID,
		"corpusid":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData1).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	e.POST("/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*2) * time.Second)
	e.GET("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskError)
	<-time.After(time.Duration(5) * time.Second)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}
