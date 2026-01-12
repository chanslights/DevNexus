package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Executor struct {
	cli *client.Client
}

// NewExecutor 初始化Docker客户端
func NewExecutor() (*Executor, error) {
	// FromEnv 会自动读取环境变量，连接本地的 Docker Daemon
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Executor{cli: cli}, nil
}

// RunStep 在容器内执行一个步骤
// ctx: 用于超时控制
// imageName: 镜像名 (如 "golang:1.21")
// commands: 要执行的 Shell 命令列表
// workDir: 宿主机上的代码目录 (会被挂载进容器)
func (e *Executor) RunStep(ctx context.Context, imageName string, commands []string, workDir string) (string, error) {
	fmt.Printf("🐳 [Docker] 准备在镜像 %s 中执行任务...\n", imageName)
	// 1. 拉取镜像 (必须先拉取，否则 Create 会报错)
	// 生产环境应该判断镜像是否存在，这里为了演示每次都 Pull
	reader, err := e.cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("pull image failed: %v", err)
	}
	// 把拉取进度扔掉(io.Discard)或者打印到控制台，防止刷屏
	io.Copy(io.Discard, reader)
	reader.Close()

	// 2. 拼接命令
	// 将 ["go version", "echo hello"] 变成 "/bin/sh -c 'go version && echo hello'"
	// 这样保证前一个命令失败，后面就不会执行
	shellCmd := strings.Join(commands, " && ")

	// 3. 创建容器 (Create)
	resp, err := e.cli.ContainerCreate(ctx,
		&container.Config{
			Image:      imageName,
			Cmd:        []string{"/bin/sh", "-c", shellCmd}, // 核心：执行用户的脚本
			WorkingDir: "/workspace",                        // 容器内的工作目录
			Tty:        false,
		},
		&container.HostConfig{
			// 核心技术：Bind Mount
			// 格式: 宿主机路径:容器内路径
			Binds: []string{workDir + ":/workspace"},
			// 自动删除：容器跑完就销毁，保持环境干净
			AutoRemove: false,
		},
		nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("create container failed: %v", err)
	}
	containerID := resp.ID
	fmt.Printf("🐳 [Docker] 容器已创建: %s\n", containerID[:12])
	defer func() {
		// 手动删除容器，清理资源
		e.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
	}()

	// 4. 启动容器 (Start)
	if err := e.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("start container failed: %v", err)
	}

	// 5. 获取日志流 (Logs)
	// 这一步非常关键，我们要实时看到容器里的输出
	out, err := e.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true, // 实时跟随
	})
	if err != nil {
		return "", err
	}

	// 创建一个Buffer来存日志
	var logBuf bytes.Buffer

	// 使用MultiWriter: 一份写到屏幕(os.Stdout)，一份写到 buffer(logBuf)
	multiWriter := io.MultiWriter(os.Stdout, &logBuf)

	// 拿到完整的日志字符串
	fullLogs := logBuf.String()

	// Docker 的日志流是多路复用的(Multiplexed)，不能直接 Print
	// 必须用 stdcopy 分离 Stdout 和 Stderr
	// 这里直接把容器的输出打印到 OpsEngine 的控制台
	stdcopy.StdCopy(os.Stdout, multiWriter, out)

	// 6. 等待容器结束 (Wait)
	// 这一步会阻塞，直到命令执行完毕
	statusCh, errCh := e.cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return "", fmt.Errorf("step failed with exit code: %d", status.StatusCode)
		}
	}

	fmt.Printf("✅ [Docker] 任务执行成功\n")
	return fullLogs, nil
}
