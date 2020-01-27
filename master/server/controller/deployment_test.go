package controller

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestDeploymentList(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/deployment").Expect().Status(http.StatusOK)
}

func TestDeployment1(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	postdata1 := map[string]interface{}{
		"content": "test1",
	}
	postdata2 := map[string]interface{}{
		"name":    "test1",
		"content": "111",
	}
	e.POST("/deployment").WithJSON(postdata1).Expect().Status(http.StatusBadRequest)
	e.POST("/deployment").WithJSON(postdata2).Expect().Status(http.StatusOK).JSON().Object().Value("id").NotEqual(0)
	obj := e.GET("/deployment").Expect().Status(http.StatusOK).JSON().Array().First().Object()
	id := int(obj.Value("id").Number().Raw())
	e.DELETE("/deployment/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}

func TestDeployment2(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	postdata1 := map[string]interface{}{
		"name":    "test1",
		"content": "111",
	}
	postdata2 := map[string]interface{}{
		"name":    "test2",
		"content": "222",
	}
	e.POST("/deployment").WithJSON(postdata1).Expect().Status(http.StatusOK)
	e.GET("/deployment").Expect().Status(http.StatusOK).JSON().Array().Length().Equal(1)
	obj := e.GET("/deployment").Expect().Status(http.StatusOK).JSON().Array().First().Object()
	obj.Value("name").Equal("test1")
	obj.Value("content").Equal("111")
	id := int(obj.Value("id").Number().Raw())
	e.PUT("/deployment/" + strconv.Itoa(id)).WithJSON(postdata2).Expect().Status(http.StatusNoContent)
	obj = e.GET("/deployment").Expect().Status(http.StatusOK).JSON().Array().First().Object()
	obj.Value("name").Equal("test2")
	obj.Value("content").Equal("222")
	e.DELETE("/deployment/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}
