package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FetchAndParse 核心函数：拉取代码并解析配置
// repoURL: http://localhost:8080/demo.git
// commitID: 刚才 Webhook 传过来的 ID
func FethchAndParse(repoURL string, commitID string) (*PipelineConfig, error) {
	// 1.创建临时目录，用于存放代码
	// 类似于：/tmp/devnexus-build-123456
	workDir, err := os.MkdirTemp("", "devexus-build-*")
	if err != nil {
		return nil, fmt.Errorf("Failed to create temp dir: %v", err)
	}
	// ⚠️ 注意：实际生产中，任务结束后应该 defer os.RemoveAll(workDir) 清理目录
	// 这里为了方便你调试观察，我们先保留目录，不删除
	fmt.Printf("📂 工作空间已创建: %s\n", workDir)

	// 2.Clone代码
	// 这一步证明CodeVault在工作，OpsEngine像一个普通用户一样去拉取代码
	fmt.Printf("⬇️ 正在从 %s 拉取代码...\n", repoURL)
	cmd := exec.Command("git", "clone", repoURL, workDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git clone failed: %s, output: %s", err, out)
	}

	// 3.（可选）Checkout到指定的Commit ID,保证我们要构建的是用户刚刚Push的那个版本
	if commitID != "" {
		checkoutCmd := exec.Command("git", "checkout", commitID)
		checkoutCmd.Dir = workDir
		if err := checkoutCmd.Run(); err != nil {
			return nil, fmt.Errorf("git checkout failed: %v", err)
		}
	}

	// 4.读取.devnexus.yaml
	configPath := filepath.Join(workDir, ".devnexus.yaml")
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("repo missing .devnexus.yaml")
	} else if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	// 5.解析YAML
	var config PipelineConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Invalid yaml format: %v", err)
	}
	return &config, nil
}
