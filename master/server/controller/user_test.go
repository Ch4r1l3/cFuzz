package controller

import (
	"net/http"
	"strconv"
	"testing"
)

func TestUserList(t *testing.T) {
	e := getAdminExpect(t)
	e.GET("/api/user").Expect().Status(http.StatusOK).JSON().Object().Value("data").Array().Length().Equal(1)
}

func TestUser1(t *testing.T) {
	e := getAdminExpect(t)
	postdata1 := map[string]interface{}{
		"username": "abc",
		"password": "123456",
	}
	postdata2 := map[string]interface{}{
		"username": "abcd",
		"password": "123456",
	}
	postdata3 := map[string]interface{}{
		"newPassword": "1234567",
	}
	postdata4 := map[string]interface{}{
		"username": "abcd",
		"password": "12345",
	}
	e.POST("/api/user").WithJSON(postdata1).Expect().Status(http.StatusBadRequest)
	e.POST("/api/user").WithJSON(postdata4).Expect().Status(http.StatusBadRequest)
	id := int(e.POST("/api/user").WithJSON(postdata2).Expect().Status(http.StatusCreated).JSON().Object().Value("id").Number().Raw())
	e.PUT("/api/user/" + strconv.Itoa(id)).WithJSON(postdata3).Expect().Status(http.StatusCreated)
	e.DELETE("/api/user/" + strconv.Itoa(id)).Expect().Status(http.StatusNoContent)
}

func TestUser2(t *testing.T) {
	e := getExpect(t)
	postdata1 := map[string]interface{}{
		"username": "abc",
		"password": "123456",
	}
	postdata2 := map[string]interface{}{
		"username": "abcd",
		"password": "123456",
	}
	postdata3 := map[string]interface{}{
		"oldPassword": "1234",
		"newPassword": "1234567",
	}
	postdata4 := map[string]interface{}{
		"oldPassword": "123456",
		"newPassword": "123456",
	}
	id := int(e.GET("/api/user/info").Expect().Status(http.StatusOK).JSON().Object().Value("id").Number().Raw())
	e.GET("/api/user").Expect().Status(http.StatusForbidden)
	e.POST("/api/user").WithJSON(postdata1).Expect().Status(http.StatusForbidden)
	e.POST("/api/user").WithJSON(postdata2).Expect().Status(http.StatusForbidden)
	e.PUT("/api/user/" + strconv.Itoa(id)).WithJSON(postdata3).Expect().Status(http.StatusBadRequest)
	e.PUT("/api/user/" + strconv.Itoa(id)).WithJSON(postdata4).Expect().Status(http.StatusCreated)
}
