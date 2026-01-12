package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type Executor struct{}

// NewExecutor 初始化Executor（本地执行模式）
func NewExecutor() (*Executor, error) {
	return &Executor{}, nil
}

// RunStep 在本地执行命令（临时解决方案）
func (e *Executor) RunStep(ctx context.Context, imageName string, commands []string, workDir string) error {
	fmt.Printf("🔧 [Local] 准备在工作目录 %s 中执行任务...\n", workDir)

	// 切换到工作目录
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir) // 执行完后恢复原目录

	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}

	// 依次执行每个命令
	for i, cmd := range commands {
		fmt.Printf("🔧 [Local] 执行步骤 %d: %s\n", i+1, cmd)

		// 创建命令
		command := exec.Command("sh", "-c", cmd)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		// 执行命令
		if err := command.Run(); err != nil {
			return fmt.Errorf("command failed: %s, error: %v", cmd, err)
		}
	}

	fmt.Printf("✅ [Local] 任务执行成功\n")
	return nil
}
