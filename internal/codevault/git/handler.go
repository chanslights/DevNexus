package git

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config The configuration for the Git handler
type Config struct {
	RepoRoot string // The root directory of the Git repository
}

// Handler The handler for the Git protocol
type Handler struct {
	config Config
}

// NewHandler Create a new Git handler
func NewHandler(config Config) *Handler {
	// auto
	if err := os.MkdirAll(config.RepoRoot, 0755); err != nil {
		log.Printf("Warning:failed to create repository root: %v", err)
	}
	return &Handler{config: config}
}

// ServeHTTP 核心入口，让Handler实现http.Handler接口
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. 解析URL，例如： /my-project.git/info/refs
	// pathParts[1]是repo名字（my-project.git）
	// pathParts[2]是动作（info或者git-receive-pack)
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "Invalid path", 400)
		return
	}

	repoName := pathParts[0]
	actions := pathParts[1]

	// 2. 拼接仓库的物理路径
	repoPath := filepath.Join(h.config.RepoRoot, repoName)

	// 如果文件夹不存在，自动帮助用户初始化一个Git裸仓库
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Printf("Initializing new repo: %s", repoName)
		initCmd := exec.Command("git", "init", "--bare", repoPath)
		if err := initCmd.Run(); err != nil {
			http.Error(w, "Failed to init repo", 500)
			return
		}
	}
	// 3. 根据动作分发请求
	switch actions {
	case "info": // 处理 info/refs 握手
		if len(pathParts) < 3 || pathParts[2] != "refs" {
			http.Error(w, "Invalid info request", 400)
		}
		h.handleInfoRefs(w, r, repoPath)
	case "git-receive-pack": // 处理POST推送
		h.handleRPC(w, r, repoPath, "git-receive-pack")

	case "git-upload-pack": // 处理POST拉取(git clone)
		h.handleRPC(w, r, repoPath, "git-upload-pack")
	default:
		http.Error(w, "Method not allowed", 405)
	}

}

// handleInfoRefs 处理第一步：握手
func (h *Handler) handleInfoRefs(w http.ResponseWriter, r *http.Request, repoPath string) {
	service := r.URL.Query().Get("service") // 例如：git-receive-pack

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))

	// Git 协议规定的Smart HTTP响应头
	// 格式是：十六进制长度+字符串+换行
	// "# service=git-receive-pack\n" 的长度是 29 (1d)，加上前缀 "001d"
	packet := fmt.Sprintf("# service=%s\n", service)
	length := len(packet) + 4
	fmt.Fprintf(w, "%04x%s0000", length, packet)

	// 调用系统Git命令
	cmd := exec.Command("git", service[4:], "--stateless-rpc", "--advertise-refs", repoPath)
	cmd.Stdout = w // 把git的输出直接写给HTTP响应
	cmd.Run()
}

// handleRPC 处理第二步：数据传输
func (h *Handler) handleRPC(w http.ResponseWriter, r *http.Request, repoPath string, service string) {
	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", service))

	// 调用系统git命令处理数据流
	cmd := exec.Command("git", service[4:], "--stateless-rpc", repoPath)
	cmd.Stdin = r.Body // 核心，把客户端上传的数据直接塞给git命令
	cmd.Stdout = w     // 核心，把git命令的反馈直接塞回给客户端
	cmd.Run()
}
