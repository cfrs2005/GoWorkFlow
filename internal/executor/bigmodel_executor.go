package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BigModelExecutor 智谱 AI BigModel 执行器
type BigModelExecutor struct {
	client *http.Client
	apiKey string
	apiURL string
}

// NewBigModelExecutor 创建 BigModel 执行器
func NewBigModelExecutor(apiKey string) *BigModelExecutor {
	return &BigModelExecutor{
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		apiKey: apiKey,
		apiURL: "https://open.bigmodel.cn/api/paas/v4/chat/completions",
	}
}

// Name 返回执行器名称
func (e *BigModelExecutor) Name() string {
	return "bigmodel_analysis"
}

// Execute 执行任务
func (e *BigModelExecutor) Execute(ctx context.Context, input map[string]interface{}, jobContext map[string]string) (map[string]interface{}, error) {
	// 从 job context 获取转录文本
	transcript := jobContext["transcript"]
	if transcript == "" {
		// 如果 context 中没有，尝试从 input 获取
		if t, ok := input["transcript"].(string); ok {
			transcript = t
		} else {
			return nil, fmt.Errorf("missing transcript in job context")
		}
	}

	// 检查文本长度
	if len(transcript) > 10000 {
		transcript = transcript[:10000] + "..." // 截断过长的文本
	}

	// 生成多个分析结果
	results := make(map[string]interface{})

	// 1. 生成阅读摘要
	summary, err := e.generateContent(ctx, transcript, "summary")
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}
	results["summary"] = summary

	// 2. 生成思维导图
	mindmap, err := e.generateContent(ctx, transcript, "mindmap")
	if err != nil {
		return nil, fmt.Errorf("failed to generate mindmap: %w", err)
	}
	results["mindmap"] = mindmap

	// 3. 重点分析
	keyPoints, err := e.generateContent(ctx, transcript, "key_points")
	if err != nil {
		return nil, fmt.Errorf("failed to generate key points: %w", err)
	}
	results["key_points"] = keyPoints

	// 4. 个人认知
	insights, err := e.generateContent(ctx, transcript, "insights")
	if err != nil {
		return nil, fmt.Errorf("failed to generate insights: %w", err)
	}
	results["insights"] = insights

	return results, nil
}

// generateContent 生成特定类型的内容
func (e *BigModelExecutor) generateContent(ctx context.Context, transcript, contentType string) (string, error) {
	// 根据内容类型构建不同的提示词
	prompt := e.buildPrompt(transcript, contentType)

	// 如果没有配置 API key，返回模拟数据
	if e.apiKey == "" || e.apiKey == "your_api_key_here" {
		return e.getMockContent(contentType, transcript), nil
	}

	// 调用 BigModel API
	request := BigModelRequest{
		Model: "glm-4-air",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		TopP:        0.9,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.apiURL, bytes.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))

	resp, err := e.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned error: %s, body: %s", resp.Status, string(body))
	}

	var response BigModelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}

// buildPrompt 构建提示词
func (e *BigModelExecutor) buildPrompt(transcript, contentType string) string {
	prompts := map[string]string{
		"summary": fmt.Sprintf(`请基于以下视频转录内容，生成一篇简洁的阅读摘要（300-500字）。
要求：
1. 概括视频的主要主题和核心观点
2. 使用清晰的段落结构
3. 突出关键信息和重要结论

视频转录内容：
%s

请生成摘要：`, transcript),

		"mindmap": fmt.Sprintf(`请基于以下视频转录内容，生成一个结构化的思维导图（使用 Markdown 格式）。
要求：
1. 使用分层的列表结构
2. 第一层是主题
3. 第二层是关键要点
4. 第三层是具体细节
5. 使用简洁的短语

视频转录内容：
%s

请生成思维导图（Markdown 格式）：`, transcript),

		"key_points": fmt.Sprintf(`请基于以下视频转录内容，提取并分析重点内容。
要求：
1. 列出 5-8 个关键要点
2. 每个要点包含简短的解释
3. 标注重要性等级（⭐⭐⭐ 高 / ⭐⭐ 中 / ⭐ 低）
4. 使用 Markdown 格式

视频转录内容：
%s

请生成重点分析：`, transcript),

		"insights": fmt.Sprintf(`请基于以下视频转录内容，提供个人认知和深度思考。
要求：
1. 从不同角度分析视频内容的价值
2. 提出可能的延伸思考或应用场景
3. 指出内容的局限性或可改进之处
4. 总结个人收获和启发
5. 使用 Markdown 格式

视频转录内容：
%s

请生成个人认知：`, transcript),
	}

	prompt, ok := prompts[contentType]
	if !ok {
		return transcript
	}

	return prompt
}

// getMockContent 返回模拟内容（用于演示）
func (e *BigModelExecutor) getMockContent(contentType, transcript string) string {
	// 分析转录文本的关键词
	keywords := e.extractKeywords(transcript)

	mockContent := map[string]string{
		"summary": fmt.Sprintf(`# 视频摘要

本视频主要探讨了关于 %s 的核心概念和实际应用。

## 主要内容

视频开头介绍了基础概念，随后深入讲解了相关技术细节。讲者通过实例演示，帮助观众更好地理解这些抽象的概念。

## 核心观点

1. **理论基础**：详细阐述了 %s 的理论框架
2. **实践应用**：展示了多个真实场景的应用案例
3. **未来趋势**：探讨了该领域的发展方向和潜在机会

## 总结

这是一个内容丰富、结构清晰的教学视频，适合初学者了解相关知识，也为有一定基础的学习者提供了新的视角。`, keywords, keywords),

		"mindmap": fmt.Sprintf(`# 视频思维导图

## 🎯 核心主题：%s

### 📚 基础概念
- 定义和背景
  - 历史发展
  - 核心理论
- 相关技术
  - 技术栈
  - 工具链

### 💡 关键要点
- 主要特性
  - 优势分析
  - 应用场景
- 实现方法
  - 技术细节
  - 最佳实践

### 🚀 实践应用
- 案例分析
  - 成功案例
  - 经验总结
- 实施步骤
  - 准备工作
  - 执行计划

### 🔮 未来展望
- 发展趋势
- 挑战与机遇`, keywords),

		"key_points": fmt.Sprintf(`# 重点分析

## ⭐⭐⭐ 高优先级要点

### 1. 核心概念理解
%s 是本视频的核心主题，理解其本质对后续学习至关重要。

### 2. 实际应用场景
视频展示了多个真实案例，这些案例可以直接应用到实际工作中。

## ⭐⭐ 中优先级要点

### 3. 技术实现细节
讲解了具体的技术实现方法，包括工具选择和配置方式。

### 4. 常见问题解答
总结了学习过程中容易遇到的问题及其解决方案。

## ⭐ 补充要点

### 5. 扩展阅读资源
提供了额外的学习资源和参考材料。

### 6. 社区和生态
介绍了相关的开源项目和社区资源。`, keywords),

		"insights": fmt.Sprintf(`# 个人认知与思考

## 💭 内容价值分析

这个视频从多个维度展现了 %s 的全貌，不仅包括理论知识，还结合了实践经验。

### 理论层面
视频构建了完整的知识框架，帮助观众建立系统性的理解。

### 实践层面
通过具体案例，展示了理论如何转化为实际应用。

## 🎯 应用场景延伸

1. **教育培训**：可以作为教学材料使用
2. **项目开发**：提供了可借鉴的实施方案
3. **技术研究**：为进一步研究提供了方向

## 🔍 批判性思考

### 优势
- 内容结构清晰
- 讲解深入浅出
- 案例丰富实用

### 可改进之处
- 某些技术细节可以进一步展开
- 可以增加更多对比分析
- 建议补充最新的发展动态

## 💡 个人收获

通过这个视频，我对 %s 有了更深入的理解，特别是在实际应用方面获得了很多启发。建议结合相关文档和实践项目，进一步巩固所学知识。`, keywords, keywords),
	}

	content, ok := mockContent[contentType]
	if !ok {
		return "内容生成中..."
	}

	return content
}

// extractKeywords 从文本中提取关键词（简单实现）
func (e *BigModelExecutor) extractKeywords(text string) string {
	keywords := []string{
		"artificial intelligence", "machine learning", "deep learning",
		"neural networks", "data science", "technology",
	}

	text = text[:min(500, len(text))]
	for _, keyword := range keywords {
		if contains(text, keyword) {
			return keyword
		}
	}

	return "相关主题"
}

func contains(text, substr string) bool {
	return len(text) >= len(substr) && (text[:len(substr)] == substr || contains(text[1:], substr))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BigModel API 请求结构
type BigModelRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// BigModel API 响应结构
type BigModelResponse struct {
	ID      string   `json:"id"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
