package controller

import (
	"net/http"
	"strconv"
	"testing"
)

func TestImageList(t *testing.T) {
	e := getExpect(t)
	e.GET("/api/image").Expect().Status(http.StatusOK)
}

func TestImage1(t *testing.T) {
	e := getExpect(t)
	postdata1 := map[string]interface{}{
		"content": "test1",
	}
	postdata2 := map[string]interface{}{
		"name":         "test1",
		"isDeployment": true,
		"content":      "111",
	}
	e.POST("/api/image").WithJSON(postdata1).Expect().Status(http.StatusBadRequest)
	e.POST("/api/image").WithJSON(postdata2).Expect().Status(http.StatusCreated).JSON().Object().Value("id").NotEqual(0)
	obj := e.GET("/api/image").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object()
	id := int(obj.Value("id").Number().Raw())
	ae := getAdminExpect(t)
	ae.GET("/api/image").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().NotEqual(0)
	e.DELETE("/api/image/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}

func TestImage2(t *testing.T) {
	e := getExpect(t)
	postdata1 := map[string]interface{}{
		"name":         "test1",
		"isDeployment": true,
		"content":      "111",
	}
	postdata2 := map[string]interface{}{
		"name":         "test2",
		"isDeployment": true,
		"content":      "222",
	}
	e.POST("/api/image").WithJSON(postdata1).Expect().Status(http.StatusCreated)
	e.GET("/api/image").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	e.GET("/api/image").WithQuery("offset", 0).WithQuery("limit", 0).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(0)
	e.GET("/api/image").WithQuery("offset", 0).WithQuery("limit", 1).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	e.GET("/api/image").WithQuery("name", "1").WithQuery("offset", 0).WithQuery("limit", 1).Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
	obj := e.GET("/api/image").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object()
	obj.Value("name").Equal("test1")
	obj.Value("content").Equal("111")
	id := int(obj.Value("id").Number().Raw())
	e.PUT("/api/image/" + strconv.Itoa(id)).WithJSON(postdata2).Expect().Status(http.StatusCreated)
	obj = e.GET("/api/image").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().First().Object()
	obj.Value("name").Equal("test2")
	obj.Value("content").Equal("222")
	e.DELETE("/api/image/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}
