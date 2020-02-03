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
	e.GET("/api/deployment").Expect().Status(http.StatusOK)
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
	e.POST("/api/deployment").WithJSON(postdata1).Expect().Status(http.StatusBadRequest)
	e.POST("/api/deployment").WithJSON(postdata2).Expect().Status(http.StatusCreated).JSON().Object().Value("id").NotEqual(0)
	obj := e.GET("/api/deployment").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object()
	id := int(obj.Value("id").Number().Raw())
	e.DELETE("/api/deployment/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
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
	e.POST("/api/deployment").WithJSON(postdata1).Expect().Status(http.StatusCreated)
	e.GET("/api/deployment").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	e.GET("/api/deployment/simplist").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	e.GET("/api/deployment/count").Expect().Status(http.StatusOK).JSON().Object().Value("count").Equal(1)
	e.GET("/api/deployment").WithQuery("offset", 0).WithQuery("limit", 0).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(0)
	e.GET("/api/deployment").WithQuery("offset", 0).WithQuery("limit", 1).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	e.GET("/api/deployment").WithQuery("name", "1").WithQuery("offset", 0).WithQuery("limit", 1).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	obj := e.GET("/api/deployment").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object()
	obj.Value("name").Equal("test1")
	obj.Value("content").Equal("111")
	id := int(obj.Value("id").Number().Raw())
	e.PUT("/api/deployment/" + strconv.Itoa(id)).WithJSON(postdata2).Expect().Status(http.StatusCreated)
	obj = e.GET("/api/deployment").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object()
	obj.Value("name").Equal("test2")
	obj.Value("content").Equal("222")
	e.DELETE("/api/deployment/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}
