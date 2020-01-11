package main

import (
	"fmt"
	"github.com/Ch4r1l3/cFuzz/master/server/config"
	"github.com/Ch4r1l3/cFuzz/master/server/service"
	"time"
)

func main() {
	config.Setup()
	service.Setup()
	service.TestProxy()

	deployment, _ := service.GenerateDeployment(1, "a", "nginx:1.7.9", 1)
	if err := service.CreateDeploy(deployment); err != nil {
		fmt.Println(err)
		return
	}
	err := service.CreateServiceByTaskID(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	<-time.After(time.Duration(10) * time.Second)

	service.DeleteServiceByTaskID(1)
	service.DeleteDeployByTaskID(1)

	for {
		time.Sleep(time.Second)
	}
}
