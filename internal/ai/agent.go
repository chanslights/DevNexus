package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Agent struct {
	ApiKey  string
	BaseURL string // 例如 https://api.deepseek.com/v1
	Model   string // 例如 deepseek-coder
}

// NewAgent 初始化
func NewAgent(apiKey string) *Agent {
	return &Agent{
		ApiKey:  apiKey,
		BaseURL: "https://api.deepseek.com", // 或者 https://api.openai.com/v1
		Model:   "deepseek-chat",            // 或者 gpt-3.5-turbo
	}
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response 结构体
type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func (a *Agent) AnalyzeLog(errorLog string) (string, error) {
	fmt.Println("🤖 [AI] 正在思考中... (分析错误原因)")
	prompt := fmt.Sprintf(`
你是一个 DevOps 专家。请分析下面的构建/部署错误日志，并给出简短的修复建议。
不要废话，直接说原因和解决办法。

错误日志：
%s
`, errorLog)
	reqBody := chatRequest{
		Model: a.Model,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	// 发送HTTP请求
	req, err := http.NewRequest("POST", a.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.ApiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(req.Body)
		return "", fmt.Errorf("AI API error: %s", string(body))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(req.Body).Decode(&chatResp); err != nil {
		return "", err
	}
	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}
	return "AI 没有任何建议", nil
}
