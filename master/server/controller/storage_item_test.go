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

func TestStorageItemList(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/api/storage_item").Expect().Status(http.StatusOK)
}

func TestStorageItem1(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/api/storage_item").Expect().Status(http.StatusOK)
	err := ioutil.WriteFile("./storageItem_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./storageItem_test")
	}()
	e.POST("/api/storage_item").Expect().Status(http.StatusBadRequest)
	e.POST("/api/storage_item").WithMultipart().WithFile("file", "storageItem_test").Expect().Status(http.StatusBadRequest)
	e.POST("/api/storage_item").
		WithMultipart().
		WithFile("file", "storageItem_test").
		WithFormField("name", "afl").
		WithFormField("type", "fuzzer").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("id").NotEqual(0)
	e.POST("/api/storage_item").
		WithMultipart().
		WithFile("file", "storageItem_test").
		WithFormField("name", "afl").
		WithFormField("type", "fuzzer").
		Expect().
		Status(http.StatusBadRequest)

	e.GET("/api/storage_item").Expect().Status(http.StatusOK).JSON().Array().Length().Equal(1)
	e.GET("/api/storage_item").WithQuery("limit", "0").Expect().Status(http.StatusOK).JSON().Array().Length().Equal(0)
	obj := e.GET("/api/storage_item").Expect().
		Status(http.StatusOK).JSON().Array().First().Object()
	obj.ValueEqual("name", "afl")
	obj.ValueEqual("type", "fuzzer")
	id := int(obj.Value("id").Number().Raw())
	e.DELETE("/api/storage_item/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}

func TestStorageItem2(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.POST("/api/storage_item/exist").Expect().Status(http.StatusBadRequest)
	postdata := map[string]interface{}{
		"name": "afl",
		"type": "fuzzer",
		"path": "/tmp/test",
	}
	e.POST("/api/storage_item/exist").
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK).JSON().Object().Value("id").NotEqual(0)

	e.GET("/api/storage_item").Expect().Status(http.StatusOK).JSON().Array().Length().Equal(1)
	obj := e.GET("/api/storage_item").Expect().
		Status(http.StatusOK).JSON().Array().First().Object()
	obj.ValueEqual("name", "afl")
	obj.ValueEqual("type", "fuzzer")
	obj.ValueEqual("existsInImage", true)
	id := int(obj.Value("id").Number().Raw())
	e.DELETE("/api/storage_item/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}
