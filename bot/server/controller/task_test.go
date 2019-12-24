package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/gavv/httpexpect"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
		"fuzzerName": "afl",
		"maxTime":    100,
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

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("name", "afl").
		Expect().
		Status(http.StatusOK)

	postdata := map[string]interface{}{
		"fuzzerName": "afl",
		"maxTime":    100,
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
	obj.Keys().ContainsOnly("corpusDir", "targetDir", "targetPath", "fuzzerName", "maxTime", "status", "arguments", "environments")
	obj.Value("corpusDir").Equal("")
	obj.Value("targetDir").Equal("")
	obj.Value("targetPath").Equal("")
	obj.Value("fuzzerName").Equal("afl")
	obj.Value("maxTime").Equal(100)
	obj.Value("status").Equal(config.TASK_CREATED)

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/fuzzer/afl").Expect().Status(http.StatusNoContent)

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

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer").
		WithFormField("name", "afl").
		Expect().
		Status(http.StatusOK)

	postdata := map[string]interface{}{
		"fuzzerName": "afl",
		"maxTime":    100,
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

	e.POST("/task/target").
		WithMultipart().
		WithFile("file", "target").
		Expect().
		Status(http.StatusOK)

	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Object()

	obj.Value("targetDir").NotEqual("")
	obj.Value("targetPath").NotEqual("")

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
}
