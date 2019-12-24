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

func TestFuzzerList(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/fuzzer").Expect().Status(http.StatusOK)
}

func TestFuzzer(t *testing.T) {
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
	e.GET("/fuzzer").Expect().
		Status(http.StatusOK).
		JSON().
		Array().First().Object().ValueEqual("name", "afl")
	e.DELETE("/fuzzer/afl").Expect().Status(http.StatusNoContent)

}

func TestFuzzerWithZipFile(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	err := createZipfile("fuzzer.zip", config.ServerConf.DefaultFuzzerName, []byte("afl"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("fuzzer.zip")
	}()

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer.zip").
		WithFormField("name", "aflzip").
		Expect().
		Status(http.StatusOK)

	e.GET("/fuzzer").Expect().
		Status(http.StatusOK).
		JSON().
		Array().First().Object().ValueEqual("name", "aflzip")
	e.DELETE("/fuzzer/aflzip").Expect().Status(http.StatusNoContent)

}

func TestFuzzerWithWrongZipFile1(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	err := createZipfile("fuzzer.zip", "abc", []byte("afl"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("fuzzer.zip")
	}()

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", "fuzzer.zip").
		WithFormField("name", "aflzip").
		Expect().
		Status(http.StatusBadRequest)

	e.GET("/fuzzer").Expect().
		Status(http.StatusOK).
		JSON().
		Array().Empty()

}

func TestFuzzerWithWrongZipFile2(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	filename := "abc.zip"
	err := ioutil.WriteFile(filename, []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(filename)
	}()

	e.POST("/fuzzer").
		WithMultipart().
		WithFile("file", filename).
		WithFormField("name", "aflzip").
		Expect().
		Status(http.StatusBadRequest)

	e.GET("/fuzzer").Expect().
		Status(http.StatusOK).
		JSON().
		Array().Empty()

}
