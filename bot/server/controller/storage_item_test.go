package controller

import (
	"github.com/gavv/httpexpect"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestStorageItemList(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/storage_item").Expect().Status(http.StatusOK)
}

func TestStorageItem1(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/storage_item").Expect().Status(http.StatusOK)
	err := ioutil.WriteFile("./storageItem_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./storageItem_test")
	}()
	e.POST("/storage_item").Expect().Status(http.StatusBadRequest)
	e.POST("/storage_item").WithMultipart().WithFile("file", "storageItem_test").Expect().Status(http.StatusBadRequest)
	e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "storageItem_test").
		WithFormField("type", "fuzzer").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("id").NotEqual(0)

	e.GET("/storage_item").Expect().Status(http.StatusOK).JSON().Array().Length().Equal(1)
	obj := e.GET("/storage_item").Expect().
		Status(http.StatusOK).JSON().Array().First().Object()
	obj.ValueEqual("type", "fuzzer")
	id := int(obj.Value("id").Number().Raw())
	e.DELETE("/storage_item/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}

func TestStorageItem2(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.POST("/storage_item/exist").Expect().Status(http.StatusBadRequest)
	err := ioutil.WriteFile("./storageItem_test", []byte("afl"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll("./storageItem_test")
	}()

	absPath, err := filepath.Abs("./storageItem_test")
	if err != nil {
		t.Fatal(err)
	}

	postdata := map[string]interface{}{
		"type": "fuzzer",
		"path": absPath,
	}
	e.POST("/storage_item/exist").
		WithJSON(postdata).
		Expect().
		Status(http.StatusOK).JSON().Object().Value("id").NotEqual(0)

	e.GET("/storage_item").Expect().Status(http.StatusOK).JSON().Array().Length().Equal(1)
	obj := e.GET("/storage_item").Expect().
		Status(http.StatusOK).JSON().Array().First().Object()
	obj.ValueEqual("type", "fuzzer")
	obj.ValueEqual("existsInImage", true)
	id := int(obj.Value("id").Number().Raw())
	e.GET("/storage_item/" + strconv.Itoa(id)).Expect().Status(http.StatusOK)
	e.DELETE("/storage_item/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}

func TestStorageItem3(t *testing.T) {
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

	e.POST("/storage_item").
		WithMultipart().
		WithFile("file", "tmp.zip").
		WithFormField("targetRelPath", "../../abc").
		Expect().
		Status(http.StatusBadRequest).JSON().Object().Value("error").NotEqual("")

}
