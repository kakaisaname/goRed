package main

import (
	"github.com/tietang/go-eureka-client/eureka"
	"time"
)

//本机输入 http://127.0.0.1:8762/ 查看
func main1() {
	cfg := eureka.Config{
		DialTimeout: time.Second * 10, //超时参数
	}
	client := eureka.NewClientByConfig([]string{
		"http://127.0.0.1:8762/eureka",
	}, cfg)
	appName := "Go-Example" //应用名称
	instance := eureka.NewInstanceInfo(
		"test.com", appName, //test.com  注册名称
		"127.0.0.2",
		8080, 30, //30 心跳周期    false 是否禁用ssl
		false)
	client.RegisterInstance(appName, instance)
	client.Start()
	c := make(chan int, 1)
	<-c
}
