package controller

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/models"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func listenCallback() int {
	taskIDChan := make(chan int)
	srv := &http.Server{Addr: ":38232"}
	defer func() {
		srv.Shutdown(nil)
	}()
	go func() {
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				fmt.Fprintf(w, "error")
				return
			}
			taskID := r.Form.Get("taskID")
			id, err := strconv.Atoi(taskID)
			if err != nil {
				fmt.Fprintf(w, "error")
				return
			}
			taskIDChan <- id
			fmt.Fprintf(w, "ok")
		})
		srv.ListenAndServe()
	}()
	select {
	case <-time.After(60 * time.Second):
		return 0
	case id := <-taskIDChan:
		return id
	}
	return 0
}

func TestTaskList(t *testing.T) {
	e := getExpect(t)
	e.GET("/api/task").Expect().Status(http.StatusOK)
}

func TestTask1(t *testing.T) {
	e := getExpect(t)
	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "11123",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	err := ioutil.WriteFile("./fuzzer_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer_test")
	}()

	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "fuzzer_test").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          100,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"callbackUrl":   "http://127.0.0.1/callback",
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	obj := e.GET("/api/task").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object()
	obj.Keys().ContainsOnly("id", "imageID", "time", "fuzzerID", "corpusID", "targetID", "status", "errorMsg", "environments", "arguments", "name", "fuzzCycleTime", "startedAt", "crashNum", "userID", "callbackUrl")
	obj.Value("id").NotEqual(0)
	obj.Value("imageID").NotEqual(0)
	obj.Value("time").NotEqual(0)
	obj.Value("fuzzCycleTime").NotEqual(0)
	obj.Value("fuzzerID").NotEqual(0)
	obj.Value("environments").Array().Elements("123", "2333")
	obj.Value("arguments").Object().Value("a1").Equal("a2")
	obj.Value("arguments").Object().Value("a2").Equal("a3")
	obj.Value("status").NotEqual("")
	obj.Value("startedAt").Equal(0)
	obj.Value("callbackUrl").Equal("http://127.0.0.1/callback")
	ae := getAdminExpect(t)
	ae.GET("/api/task").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().NotEqual(0)
	e.GET("/api/task").WithQuery("name", "t").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	e.GET("/api/task").WithQuery("name", "a").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(0)
	e.GET("/api/task").WithQuery("name", "t").WithQuery("offset", 0).WithQuery("limit", 0).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(0)
	e.GET("/api/task").WithQuery("name", "t").WithQuery("offset", 1).WithQuery("limit", 1).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(0)
	e.GET("/api/task").WithQuery("name", "t").WithQuery("offset", 0).WithQuery("limit", 1).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	e.GET("/api/task").WithQuery("name", "t").WithQuery("offset", 1).WithQuery("limit", 0).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(0)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
}

func TestTask2(t *testing.T) {
	e := getExpect(t)
	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "11123",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	err := ioutil.WriteFile("./fuzzer_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer_test")
	}()

	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "fuzzer_test").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          100,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"callbackUrl":   "http://127.0.0.1/callback",
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}
	taskPostData2 := map[string]interface{}{
		"imageID": -1,
	}
	taskPostData3 := map[string]interface{}{
		"fuzzerID": -1,
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
	taskPostData7 := map[string]interface{}{
		"callbackUrl": "abc",
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	e.PUT("/api/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData2).Expect().Status(http.StatusBadRequest)
	e.PUT("/api/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData3).Expect().Status(http.StatusBadRequest)
	e.PUT("/api/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData4).Expect().Status(http.StatusBadRequest)
	e.PUT("/api/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData5).Expect().Status(http.StatusCreated)
	e.PUT("/api/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData6).Expect().Status(http.StatusCreated)
	e.PUT("/api/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData7).Expect().Status(http.StatusCreated)

	e.GET("/api/task").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object().Value("environments").Array().Elements("2", "3")
	e.GET("/api/task").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object().Value("callbackUrl").Equal("abc")

	obj := e.GET("/api/task").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object().Value("arguments").Object()
	obj.Value("a3").Equal("a4")
	obj.Value("a4").Equal("a5")

	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
}

func TestTask3(t *testing.T) {
	e := getExpect(t)

	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "ch4r1l3/cfuzz:test-bot",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 3,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	taskPostData4 := map[string]interface{}{
		"targetID": targetID,
		"corpusID": corpusID,
	}
	e.PUT("/api/task/" + strconv.Itoa(taskID)).WithJSON(taskPostData4).Expect().Status(http.StatusCreated)

	e.GET("/api/task").WithQuery("limit", "0").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(0)
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*2) * time.Second)
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskRunning)
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(2) * time.Second)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}

func TestTask4(t *testing.T) {
	e := getExpect(t)

	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "ch4r1l3/cfuzz:test-bot",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 2,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"callbackUrl":   "http://127.0.0.1:38232/callback",
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	//<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*4) * time.Second)
	callbacbID := listenCallback()
	if callbacbID != taskID {
		t.Errorf("callback id wrong")
	}
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskStopped)
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusBadRequest)
	<-time.After(time.Duration(2) * time.Second)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}

func TestTask5(t *testing.T) {
	e := getExpect(t)

	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "ch4r1l3/cfuzz:test-bot",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 5,
		"fuzzCycleTime": 15,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*4) * time.Second)
	e.GET("/api/task/" + strconv.Itoa(taskID) + "/crash").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().NotEqual(0)
	obj := e.GET("/api/task/" + strconv.Itoa(taskID) + "/result").Expect().Status(http.StatusOK).JSON().Object()
	obj.Keys().ContainsOnly("command", "timeExecuted", "updateAt", "stats", "id", "taskid")
	obj.Value("command").NotEqual("")
	obj.Value("timeExecuted").NotEqual(0)
	obj.Value("updateAt").NotEqual(0)
	obj.Value("stats").NotNull()
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(2) * time.Second)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}

func TestTask6(t *testing.T) {
	e := getExpect(t)

	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "deadbeef:v1",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 8,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*4) * time.Second)
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskError)
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("errorMsg").Equal("failed to create deployment")
	<-time.After(time.Duration(2) * time.Second)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}

func TestTask7(t *testing.T) {
	e := getExpect(t)

	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "nginx:1.16-alpine",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpus").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 8,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*2) * time.Second)
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskError)
	<-time.After(time.Duration(2) * time.Second)
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusBadRequest)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}

func TestTask8(t *testing.T) {
	e := getExpect(t)

	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "ch4r1l3/cfuzz:test-exist",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	fuzzerData := map[string]string{
		"name": "f1",
		"type": "fuzzer",
		"path": "/cfuzz/test_data/afl",
	}
	targetData := map[string]string{
		"name": "t1",
		"type": "target",
		"path": "/cfuzz/test_data/test",
	}
	corpusData := map[string]string{
		"name": "c1",
		"type": "corpus",
		"path": "/cfuzz/test_data/corpus",
	}

	fuzzerID := int(e.POST("/api/storage_item/exist").WithJSON(fuzzerData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item/exist").WithJSON(targetData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item/exist").WithJSON(corpusData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 5,
		"fuzzCycleTime": 15,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*4) * time.Second)
	e.GET("/api/task/" + strconv.Itoa(taskID) + "/crash").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().NotEqual(0)
	obj := e.GET("/api/task/" + strconv.Itoa(taskID) + "/result").Expect().Status(http.StatusOK).JSON().Object()
	obj.Keys().ContainsOnly("command", "timeExecuted", "updateAt", "stats", "id", "taskid")
	obj.Value("command").NotEqual("")
	obj.Value("timeExecuted").NotEqual(0)
	obj.Value("updateAt").NotEqual(0)
	obj.Value("stats").NotNull()
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(2) * time.Second)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}

func TestTask9(t *testing.T) {
	e := getExpect(t)

	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "ch4r1l3/cfuzz:test-exist",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	fuzzerData := map[string]string{
		"name": "f1",
		"type": "fuzzer",
		"path": "/cfuzz/test/afl",
	}
	targetData := map[string]string{
		"name": "t1",
		"type": "target",
		"path": "/cfuzz/test_data/test",
	}
	corpusData := map[string]string{
		"name": "c1",
		"type": "corpus",
		"path": "/cfuzz/test_data/corpus",
	}

	fuzzerID := int(e.POST("/api/storage_item/exist").WithJSON(fuzzerData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item/exist").WithJSON(targetData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item/exist").WithJSON(corpusData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 8,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*2) * time.Second)
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskError)
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("errorMsg").NotEqual("")
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}

func TestTask10(t *testing.T) {
	e := getExpect(t)

	imagePostData := map[string]interface{}{
		"name":         "test",
		"isDeployment": false,
		"content":      "ch4r1l3/cfuzz:test-bot",
	}
	imageID := int(e.POST("/api/image").WithJSON(imagePostData).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test.zip").WithFormField("relPath", "test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpuscc").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 3,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*2) * time.Second)
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskRunning)
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(2) * time.Second)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}

func TestTask11(t *testing.T) {
	e := getExpect(t)

	fuzzerID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/afl").WithFormField("name", "afl").WithFormField("type", "fuzzer").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	targetID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/test.zip").WithFormField("relPath", "test").WithFormField("name", "test_target").WithFormField("type", "target").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	corpusID := int(e.POST("/api/storage_item").WithMultipart().WithFile("file", "../test_data/corpus").WithFormField("name", "test_corpuscc").WithFormField("type", "corpus").Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	postdata1 := map[string]interface{}{
		"name":         "test1",
		"isDeployment": true,
		"content": `
apiVersion: apps/v1 
kind: Deployment
spec:
  replicas: 1 
  template:
    spec:
      containers:
      - name: nginx
        image: ch4r1l3/cfuzz:test-bot
        ports:
        - containerPort: 80
`,
	}
	imageID := int(e.POST("/api/image").WithJSON(postdata1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())

	taskPostData1 := map[string]interface{}{
		"name":          "test",
		"imageID":       imageID,
		"time":          config.KubernetesConf.CheckTaskTime * 3,
		"fuzzCycleTime": 60,
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"environments":  []string{"123", "2333"},
		"arguments": map[string]string{
			"a1": "a2",
			"a2": "a3",
		},
	}

	taskID := int(e.POST("/api/task").WithJSON(taskPostData1).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/start").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(config.KubernetesConf.CheckTaskTime*2) * time.Second)
	e.GET("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskRunning)
	e.POST("/api/task/" + strconv.Itoa(taskID) + "/stop").Expect().Status(http.StatusAccepted)
	<-time.After(time.Duration(2) * time.Second)
	e.DELETE("/api/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/api/image/" + strconv.Itoa(imageID)).Expect().Status(http.StatusNoContent)
}
