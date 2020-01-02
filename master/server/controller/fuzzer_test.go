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

func TestFuzzerList(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/fuzzer").Expect().Status(http.StatusOK)
}

func TestFuzzer1(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/fuzzer").Expect().Status(http.StatusOK)
	err := ioutil.WriteFile("./fuzzer_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./fuzzer_test")
	}()
	e.POST("/fuzzer").Expect().Status(http.StatusBadRequest)
	e.POST("/fuzzer").WithMultipart().WithFile("file", "fuzzer_test").Expect().Status(http.StatusBadRequest)
	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer_test").
		WithFormField("name", "afl").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("id").NotEqual(0)
	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer_test").
		WithFormField("name", "afl").
		Expect().
		Status(http.StatusBadRequest)

	e.GET("/fuzzer").Expect().Status(http.StatusOK).JSON().Array().Length().Equal(1)
	obj := e.GET("/fuzzer").Expect().
		Status(http.StatusOK).JSON().Array().First().Object()
	obj.ValueEqual("name", "afl")
	id := int(obj.Value("id").Number().Raw())
	e.DELETE("/fuzzer/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}
