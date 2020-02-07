package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/gavv/httpexpect"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestTaskRetrieve(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/task").Expect().Status(http.StatusBadRequest)
}

func TestTask1(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	postdata := map[string]interface{}{
		"fuzzerID":      1,
		"corpusID":      2,
		"targetID":      3,
		"maxTime":       100,
		"fuzzCycleTime": 60,
		"arguments": map[string]string{
			"A": "1",
			"B": "2",
		},
		"environments": []string{
			"A=1",
			"B=2",
		},
	}
	e.POST("/task").
		WithJSON(postdata).
		Expect().Status(http.StatusBadRequest)
}

func TestTask2(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	err := ioutil.WriteFile("./fuzzer", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer")
	}()

	fuzzerID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("type", "fuzzer").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())
	corpusID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("type", "corpus").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())
	targetID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("type", "target").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())

	postdata := map[string]interface{}{
		"fuzzerID":      fuzzerID,
		"corpusID":      corpusID,
		"targetID":      targetID,
		"maxTime":       100,
		"fuzzCycleTime": 60,
		"arguments": map[string]string{
			"A": "1",
			"B": "2",
		},
		"environments": []string{
			"A=1",
			"B=2",
		},
	}
	e.POST("/task").
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK)
	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Object()
	obj.Keys().ContainsOnly("corpusID", "targetID", "fuzzerID", "maxTime", "status", "arguments", "environments", "fuzzCycleTime", "errorMsg")
	obj.Value("fuzzerID").Equal(fuzzerID)
	obj.Value("targetID").Equal(targetID)
	obj.Value("corpusID").Equal(corpusID)
	obj.Value("maxTime").Equal(100)
	obj.Value("fuzzCycleTime").Equal(60)
	obj.Value("status").Equal(models.TaskCreated)

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)

}

func TestTask3(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	err := ioutil.WriteFile("./fuzzer", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer")
	}()

	err = ioutil.WriteFile("./target", []byte("target"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./target")
	}()

	fuzzerID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("type", "fuzzer").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())
	corpusID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("type", "corpus").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())
	targetID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("type", "target").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())

	postdata := map[string]interface{}{
		"fuzzerID":      fuzzerID,
		"corpusID":      corpusID,
		"targetID":      targetID,
		"maxTime":       100,
		"fuzzCycleTime": 60,
		"arguments": map[string]string{
			"A": "1",
			"B": "2",
		},
		"environments": []string{
			"A=1",
			"B=2",
		},
	}
	e.POST("/task").
		WithJSON(postdata).
		Expect().Status(http.StatusOK)
	e.GET("/task").Expect().Status(http.StatusOK).JSON().Object()

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
}

func TestTask4(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	err := ioutil.WriteFile("./fuzzer", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer")
	}()

	err = createZipfile("tmp.zip", "abc", []byte("abc"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("tmp.zip")
	}()

	fuzzerID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("type", "fuzzer").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())

	targetID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "tmp.zip").
		WithFormField("type", "target").
		WithFormField("relPath", "abc").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())

	corpusID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "tmp.zip").
		WithFormField("type", "corpus").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())

	postdata := map[string]interface{}{
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"maxTime":       100,
		"fuzzCycleTime": 60,
		"arguments": map[string]string{
			"A": "1",
			"B": "2",
		},
		"environments": []string{
			"A=1",
			"B=2",
		},
	}
	e.POST("/task").
		WithJSON(postdata).
		Expect().Status(http.StatusOK)

	e.GET("/task").Expect().Status(http.StatusOK)

	e.POST("/task/start").Expect().Status(http.StatusNoContent)

	<-time.After(time.Duration(10) * time.Second)
	e.GET("/task").Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskError)
	e.GET("/task").Expect().Status(http.StatusOK).JSON().Object().Value("errorMsg").NotEqual("")
	e.POST("/task/stop").Expect().Status(http.StatusBadRequest)

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
}

func TestTask5(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	err := createZipfile("tmp.zip", "abc", []byte("abc"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("tmp.zip")
	}()

	fuzzerID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "../test_data/afl").
		WithFormField("type", "fuzzer").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())

	targetID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "../test_data/test").
		WithFormField("type", "target").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())

	corpusID := int(e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "tmp.zip").
		WithFormField("type", "corpus").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").
		Number().
		Raw())

	postdata := map[string]interface{}{
		"fuzzerID":      fuzzerID,
		"targetID":      targetID,
		"corpusID":      corpusID,
		"maxTime":       100,
		"fuzzCycleTime": 60,
		"arguments": map[string]string{
			"A": "1",
			"B": "2",
		},
		"environments": []string{
			"A=1",
			"B=2",
		},
	}
	e.POST("/task").
		WithJSON(postdata).
		Expect().Status(http.StatusOK)

	e.POST("/task/start").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(3) * time.Second)

	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Object()
	obj.Value("status").Equal(models.TaskRunning)

	<-time.After(time.Duration(70) * time.Second)
	e.POST("/task/stop").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(63) * time.Second)
	e.GET("/task").Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskStopped)
	e.GET("/task/result").Expect().Status(http.StatusOK).JSON().Object().Value("timeExecuted").Equal(60)

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
}

func TestTaskGET(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	e.GET("/task").Expect().Status(http.StatusBadRequest)

}
