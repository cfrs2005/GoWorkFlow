package executor

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// YouTubeASRExecutor YouTube ASR 下载执行器
type YouTubeASRExecutor struct {
	client *http.Client
}

// NewYouTubeASRExecutor 创建 YouTube ASR 执行器
func NewYouTubeASRExecutor() *YouTubeASRExecutor {
	return &YouTubeASRExecutor{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Name 返回执行器名称
func (e *YouTubeASRExecutor) Name() string {
	return "youtube_asr"
}

// Execute 执行任务
func (e *YouTubeASRExecutor) Execute(ctx context.Context, input map[string]interface{}, jobContext map[string]string) (map[string]interface{}, error) {
	// 获取 YouTube URL
	videoURL, ok := input["video_url"].(string)
	if !ok || videoURL == "" {
		return nil, fmt.Errorf("missing required parameter: video_url")
	}

	// 提取视频 ID
	videoID, err := e.extractVideoID(videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract video ID: %w", err)
	}

	// 获取字幕语言（默认英文，也支持中文）
	language := "en"
	if lang, ok := input["language"].(string); ok && lang != "" {
		language = lang
	}

	// 尝试多种方式获取字幕
	var transcript string

	// 方式 1: 使用 yt-dlp（如果安装了）
	transcript, err = e.getTranscriptWithYtDlp(ctx, videoID, language)
	if err == nil && transcript != "" {
		return map[string]interface{}{
			"video_id":   videoID,
			"transcript": transcript,
			"language":   language,
			"method":     "yt-dlp",
			"length":     len(transcript),
		}, nil
	}

	// 方式 2: 使用 YouTube Transcript API（Python 脚本）
	transcript, err = e.getTranscriptWithPython(ctx, videoID, language)
	if err == nil && transcript != "" {
		return map[string]interface{}{
			"video_id":   videoID,
			"transcript": transcript,
			"language":   language,
			"method":     "youtube-transcript-api",
			"length":     len(transcript),
		}, nil
	}

	// 方式 3: 模拟数据（用于演示）
	if transcript == "" {
		transcript = e.getMockTranscript()
		return map[string]interface{}{
			"video_id":   videoID,
			"transcript": transcript,
			"language":   language,
			"method":     "mock",
			"length":     len(transcript),
			"warning":    "Using mock data. Please install yt-dlp or youtube-transcript-api for real transcripts.",
		}, nil
	}

	return nil, fmt.Errorf("failed to get transcript from all methods")
}

// extractVideoID 从 URL 提取视频 ID
func (e *YouTubeASRExecutor) extractVideoID(videoURL string) (string, error) {
	// 支持多种 YouTube URL 格式
	// https://www.youtube.com/watch?v=VIDEO_ID
	// https://youtu.be/VIDEO_ID
	// https://www.youtube.com/embed/VIDEO_ID

	patterns := []string{
		`(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([a-zA-Z0-9_-]{11})`,
		`^([a-zA-Z0-9_-]{11})$`, // 直接是 video ID
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(videoURL)
		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("invalid YouTube URL format")
}

// getTranscriptWithYtDlp 使用 yt-dlp 获取字幕
func (e *YouTubeASRExecutor) getTranscriptWithYtDlp(ctx context.Context, videoID, language string) (string, error) {
	// 检查 yt-dlp 是否安装
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return "", fmt.Errorf("yt-dlp not found")
	}

	// 构建命令
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"--skip-download",
		"--write-auto-sub",
		"--sub-lang", language,
		"--sub-format", "json3",
		"--print", "%(subtitles)s",
		fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("yt-dlp failed: %w, output: %s", err, string(output))
	}

	// 解析输出
	transcript := e.parseYtDlpOutput(string(output))
	return transcript, nil
}

// getTranscriptWithPython 使用 Python youtube-transcript-api
func (e *YouTubeASRExecutor) getTranscriptWithPython(ctx context.Context, videoID, language string) (string, error) {
	// Python 脚本
	pythonScript := fmt.Sprintf(`
from youtube_transcript_api import YouTubeTranscriptApi
import json
try:
    transcript = YouTubeTranscriptApi.get_transcript('%s', languages=['%s', 'en'])
    text = ' '.join([entry['text'] for entry in transcript])
    print(text)
except Exception as e:
    print(f"Error: {e}", file=sys.stderr)
    exit(1)
`, videoID, language)

	// 执行 Python 脚本
	cmd := exec.CommandContext(ctx, "python3", "-c", pythonScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("python script failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// parseYtDlpOutput 解析 yt-dlp 输出
func (e *YouTubeASRExecutor) parseYtDlpOutput(output string) string {
	// 简单解析，实际可能需要更复杂的处理
	lines := strings.Split(output, "\n")
	var transcript strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "[") {
			transcript.WriteString(line)
			transcript.WriteString(" ")
		}
	}
	return strings.TrimSpace(transcript.String())
}

// getMockTranscript 返回模拟字幕（用于演示）
func (e *YouTubeASRExecutor) getMockTranscript() string {
	return `Welcome to this video about artificial intelligence and machine learning.
In this tutorial, we'll explore the fundamentals of neural networks and deep learning.
First, let's discuss what makes neural networks so powerful.
Neural networks are inspired by the human brain and consist of interconnected layers of nodes.
Each node processes information and passes it to the next layer.
Deep learning has revolutionized many fields including computer vision, natural language processing, and robotics.
Companies like Google, Facebook, and OpenAI are using these technologies to build amazing products.
In the next section, we'll dive into the mathematics behind backpropagation and gradient descent.
These algorithms allow neural networks to learn from data by adjusting their weights.
The key insight is that we can calculate how much each weight contributes to the error.
Then we update the weights to reduce this error, gradually improving the model's performance.
This process is repeated many times until the model converges to a good solution.
Thank you for watching, and I hope you found this introduction helpful.
Don't forget to subscribe and hit the notification bell for more AI content.`
}

// TranscriptEntry YouTube 字幕条目
type TranscriptEntry struct {
	Text     string  `json:"text"`
	Start    float64 `json:"start"`
	Duration float64 `json:"duration"`
}
