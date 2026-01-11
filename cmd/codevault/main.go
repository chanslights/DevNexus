package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chanslights/DevNexus/pkg/utils"
)

func main() {
	log.Printf("DevNexus starting %s", utils.GetVersion())
	// 这里未来会放入Git协议处理逻辑
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to DevNexus CodeVault(MyGit Server)")
	})
	port := ":8080"
	log.Printf("CodeVault [Git Server] running on %s", port)
	// 启动HTTP服务
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
