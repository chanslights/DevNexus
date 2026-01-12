package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/chanslights/DevNexus/internal/ai"
	"github.com/chanslights/DevNexus/internal/opsengine/docker"
	"github.com/chanslights/DevNexus/internal/opsengine/k8s"
	"github.com/chanslights/DevNexus/internal/opsengine/pipeline"
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

	// 初始化AI Agent
	aiAgent := ai.NewAgent("")

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
	fmt.Println("开始出发流水线构建...")

	// 2.1 构造clone地址。（目前都在本地构造，因此先拼一下地址）
	repoURL := fmt.Sprintf("http://localhost:8080/%s", payload.RepoName)

	// 2.2 调用Pipeline模块去拉取代码并解析
	// 这是一个耗时的操作，实际应该放入Go Channel队列里面异步执行。但当前为了演示，直接用go func跑
	go func() {
		config, workDir, err := pipeline.FetchAndParse(repoURL, payload.CommitID)
		if err != nil {
			log.Printf("❌ 流水线启动失败: %v", err)
			return
		}
		// ⚠️ 重要：任务结束后清理临时目录
		// defer os.RemoveAll(workDir)

		// 初始化Docker执行器
		executor, err := docker.NewExecutor()
		if err != nil {
			log.Printf("❌ Docker 客户端初始化失败: %v", err)
			return
		}

		k8sDeployer, err := k8s.NewDeployer()
		if err != nil {
			log.Printf("❌ K8s 客户端初始化失败: %v", err)
		}

		// 遍历执行每一个Stage
		ctx := context.Background()
		for _, stage := range config.Stages {
			// 遍历定义在循环外，用来接收日志
			var stepLogs string
			var stepErr error

			fmt.Printf("\n▶️  开始执行阶段: [%s]\n", stage.Name)

			if stage.Type == "kubernetes" {
				if k8sDeployer == nil {
					log.Printf("❌ K8s 未连接，无法部署")
					return
				}
				// 默认发布到 default 命名空间
				err := k8sDeployer.UpdateImage(ctx, "default", stage.Target, stage.NewImage)
				if err != nil {
					log.Printf("❌ 部署失败: %v", err)
					stepErr = err
					stepLogs = "Kubernetes Deployment Update Failed." // 简单占位
					return
				}
			} else {
				// 真正的执行
				_, err := executor.RunStep(ctx, stage.Image, stage.Script, workDir)
				if err != nil {
					log.Printf("❌ 阶段 [%s] 执行失败: %v\n", stage.Name, err)
					stepLogs, stepErr = executor.RunStep(ctx, stage.Image, stage.Script, workDir)
					return // 流水线中断
				}
			}

			// 错误处理与AI介入
			if stepErr != nil {
				log.Printf("❌ 阶段 [%s] 执行失败: %v", stage.Name, stepErr)
				// 呼叫 AI 进行分析
				fmt.Println("\n🚑 检测到构建失败，正在呼叫 AI 医生...")
				// 截取最后 2000 个字符的日志发给 AI (防止 Token 超出)
				logContext := stepLogs
				if len(logContext) > 2000 {
					logContext = logContext[len(logContext)-2000:]
				}
				suggestion, aiErr := aiAgent.AnalyzeLog(logContext)
				if aiErr != nil {
					fmt.Printf("⚠️ AI 分析失败: %v\n", aiErr)
				} else {
					fmt.Println("==================================================")
					fmt.Println("🤖 AI 诊断报告:")
					fmt.Println(suggestion)
					fmt.Println("==================================================")
				}
				return // 终止流水线
			}
		}
		fmt.Println("\n🎉🎉🎉 流水线全部执行成功！")
	}()

	w.WriteHeader(200)
	w.Write([]byte("Webhook received successfully"))
}
