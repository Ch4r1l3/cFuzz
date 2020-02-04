package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/models"
	"github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestAFL(t *testing.T) {
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
		"fuzzCycleTime": 15,
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

	<-time.After(time.Duration(35) * time.Second)
	e.POST("/task/stop").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(20) * time.Second)
	e.GET("/task").Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskStopped)
	e.GET("/task/result").Expect().Status(http.StatusOK).JSON().Object().Value("timeExecuted").Equal(15)

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
}

func TestLibFuzzer(t *testing.T) {
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
		WithFile("file", "../test_data/libfuzzer").
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
		WithFile("file", "../test_data/libfuzzer_target").
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
		"fuzzCycleTime": 15,
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

	<-time.After(time.Duration(31) * time.Second)
	e.POST("/task/stop").Expect().Status(http.StatusNoContent)
	<-time.After(time.Duration(20) * time.Second)
	e.GET("/task").Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(models.TaskStopped)
	e.GET("/task/crash").Expect().Status(http.StatusOK).JSON().Array().Length().NotEqual(0)
	tobj := e.GET("/task/crash").Expect().Status(http.StatusOK).JSON().Array().First().Object()
	tobj.Value("reproduceAble").Equal(true)
	tobj.Value("fileName").NotEqual("")

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(targetID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/storage_item/" + strconv.Itoa(corpusID)).Expect().Status(http.StatusNoContent)
}
