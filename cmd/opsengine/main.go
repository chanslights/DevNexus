package main

import (
	"log"
	"time"

	"github.com/chanslights/DevNexus/pkg/utils"
)

func main() {
	log.Printf("DevNexus starting %s", utils.GetVersion())
	log.Println("DevNexus OpsEngine [CI/CD Worker] is starting...")
	// 模拟一个死循环，作为后台守护进程（Daemon）
	// TODO: 未来这里会由RabbitMQ或者Webhook触发
	for {
		log.Println("OpsEngine: waiting for build tasks...")
		time.Sleep(5 * time.Second)
	}
}
