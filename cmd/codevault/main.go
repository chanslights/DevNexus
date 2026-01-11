package main

import (
	"log"
	"net/http"

	"github.com/chanslights/DevNexus/internal/codevault/git"
	"github.com/chanslights/DevNexus/pkg/utils"
)

func main() {
	log.Printf("DevNexus starting %s", utils.GetVersion())
	config := git.Config{
		RepoRoot: "./repos",
	}
	// 初始化Handler
	gitHandler := git.NewHandler(config)

	// 注册路由
	// 所有路由都交给gitHandler处理
	http.Handle("/", gitHandler)

	port := ":8080"
	log.Printf("CodeVault [Git Server] running on %s", port)
	log.Printf("Repo Storage: %s", config.RepoRoot)

	// 启动HTTP服务
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
