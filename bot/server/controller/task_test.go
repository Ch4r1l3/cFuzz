package controller

import (
	"github.com/Ch4r1l3/cFuzz/bot/server/config"
	"github.com/gavv/httpexpect"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
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
	e.DELETE("/fuzzer/afl").Expect().Status(http.StatusNoContent)
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
	postdata1 := map[string]interface{}{
		"status": config.TASK_RUNNING,
	}
	postdata2 := map[string]interface{}{
		"status": config.TASK_STOPPED,
	}
	e.POST("/task").
		WithJSON(postdata).
		Expect().Status(http.StatusOK)
	e.GET("/task").Expect().Status(http.StatusOK)

	e.PUT("/task").WithJSON(postdata1).Expect().Status(http.StatusBadRequest)

	/*
		e.POST("/task/target").
			WithMultipart().
			WithFile("file", "tmp.zip").
			Expect().
			Status(http.StatusBadRequest)

	*/
	e.POST("/task/target").
		WithMultipart().
		WithFile("file", "tmp.zip").
		WithFormField("targetRelPath", "../../abc").
		Expect().
		Status(http.StatusBadRequest).JSON().Object().Value("error").NotEqual("")

	e.POST("/task/target").
		WithMultipart().
		WithFile("file", "tmp.zip").
		WithFormField("targetRelPath", "abc").
		Expect().
		Status(http.StatusOK)

	e.PUT("/task").WithJSON(postdata1).Expect().Status(http.StatusBadRequest)

	e.POST("/task/corpus").
		WithMultipart().
		WithFile("file", "tmp.zip").
		Expect().
		Status(http.StatusOK)

	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Object()

	obj.Value("targetDir").NotEqual("")
	obj.Value("targetPath").NotEqual("")
	obj.Value("corpusDir").NotEqual("")

	e.PUT("/task").WithJSON(postdata1).Expect().Status(http.StatusOK)

	<-time.After(time.Duration(4) * time.Second)
	e.PUT("/task").WithJSON(postdata2).Expect().Status(http.StatusBadRequest)

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/fuzzer/afl").Expect().Status(http.StatusNoContent)
}

func TestTask5(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	err := ioutil.WriteFile("./target.c", []byte("#include<stdio.h>\n\nint main()\n{\nchar buf[0x10];read(0,buf,0x100);return 0;\n}\n"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	//_, err = exec.Command("gcc", "-o", "target", "./target.c").Output()
	_, err = exec.Command("/afl/afl-2.52b/afl-gcc", "-o", "target", "./target.c").Output()
	os.RemoveAll("./target.c")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("target")
	}()

	err = createZipfile("tmp.zip", "abc", []byte("abc"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("tmp.zip")
	}()

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "afl").
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
	postdata1 := map[string]interface{}{
		"status": config.TASK_RUNNING,
	}
	postdata2 := map[string]interface{}{
		"status": config.TASK_STOPPED,
	}
	e.POST("/task").
		WithJSON(postdata).
		Expect().Status(http.StatusOK)

	e.POST("/task/target").
		WithMultipart().
		WithFile("file", "target").
		Expect().
		Status(http.StatusOK)

	e.POST("/task/corpus").
		WithMultipart().
		WithFile("file", "tmp.zip").
		Expect().
		Status(http.StatusOK)

	<-time.After(time.Duration(10) * time.Second)
	e.PUT("/task").WithJSON(postdata1).Expect().Status(http.StatusOK)
	<-time.After(time.Duration(3) * time.Second)

	obj := e.GET("/task").Expect().Status(http.StatusOK).JSON().Object()
	obj.Value("status").Equal(config.TASK_RUNNING)

	e.PUT("/task").WithJSON(postdata2).Expect().Status(http.StatusOK)
	<-time.After(time.Duration(70) * time.Second)
	e.GET("/task").Expect().Status(http.StatusOK).JSON().Object().Value("status").Equal(config.TASK_STOPPED)

	e.DELETE("/task").Expect().Status(http.StatusNoContent)
	e.DELETE("/fuzzer/afl").Expect().Status(http.StatusNoContent)
}

func TestTaskUpdate(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	postdata1 := map[string]interface{}{
		"status": "abc",
	}
	postdata2 := map[string]interface{}{
		"status": config.TASK_STOPPED,
	}
	e.PUT("/task").Expect().Status(http.StatusBadRequest)
	e.PUT("/task").WithJSON(postdata1).Expect().Status(http.StatusBadRequest)
	e.PUT("/task").WithJSON(postdata2).Expect().Status(http.StatusBadRequest)

}
