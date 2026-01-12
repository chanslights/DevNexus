package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chanslights/DevNexus/pkg/types"
	"github.com/chanslights/DevNexus/pkg/utils"
)

func main() {
	log.Printf("DevNexus starting %s", utils.GetVersion())
	log.Println("DevNexus OpsEngine [CI/CD Worker] is starting...")

	http.HandleFunc("/webhook", handleWebHook)

	port := ":8081"
	log.Printf("OpsEngine is listening on port %s for webhooks...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Failed to start OpsEngine: %v", err)
	}
}

func handleWebHook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// 1.解析JSON数据
	var payload types.WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}
	// 2.打印日志（假装开始构建）
	fmt.Println("------------------------------------------------")
	fmt.Printf("🔔 收到 Webhook 通知！\n")
	fmt.Printf("📦 仓库: %s\n", payload.RepoName)
	fmt.Printf("🌿 分支: %s\n", payload.Branch)
	fmt.Printf("🔑 Commit ID: %s\n", payload.CommitID)
	fmt.Println("🚀 正在触发流水线构建... (模拟中)")
	fmt.Println("------------------------------------------------")

	w.WriteHeader(200)
	w.Write([]byte("Webhook received successfully"))
}
