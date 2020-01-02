package controller

import (
	"github.com/gavv/httpexpect"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
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
	dockerfilePostData := map[string]interface{}{
		"name":    "test",
		"content": "11123",
	}
	dockerfileID := int(e.POST("/dockerfile").WithJSON(dockerfilePostData).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	err := ioutil.WriteFile("./fuzzer_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer_test")
	}()

	fuzzerID := int(e.POST("/fuzzer").WithMultipart().WithFile("file", "fuzzer_test").WithFormField("name", "afl").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	taskPostData := map[string]interface{}{
		"dockerfileid": dockerfileID,
		"time":         100,
		"fuzzerid":     fuzzerID,
		"environments": []string{"123", "2333"},
	}

	taskID := int(e.POST("/task").WithJSON(taskPostData).Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())

	e.GET("/task").Expect().Status(http.StatusOK).JSON().Array().First().Object().Value("running").Equal(false)
	e.DELETE("/dockerfile/" + strconv.Itoa(dockerfileID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/fuzzer/" + strconv.Itoa(fuzzerID)).Expect().Status(http.StatusNoContent)
	e.DELETE("/task/" + strconv.Itoa(taskID)).Expect().Status(http.StatusNoContent)
}
