package executor

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cfrs2005/GoWorkFlow/pkg/logger"
)

// HTMLReportExecutor HTML 报告生成执行器
type HTMLReportExecutor struct {
	outputDir string
}

// NewHTMLReportExecutor 创建 HTML 报告执行器
func NewHTMLReportExecutor(outputDir string) *HTMLReportExecutor {
	// 确保输出目录存在
	if outputDir == "" {
		outputDir = "./reports"
	}
	os.MkdirAll(outputDir, 0755)

	return &HTMLReportExecutor{
		outputDir: outputDir,
	}
}

// Name 返回执行器名称
func (e *HTMLReportExecutor) Name() string {
	return "html_report"
}

// Execute 执行任务
func (e *HTMLReportExecutor) Execute(ctx context.Context, input map[string]interface{}, jobContext map[string]string) (map[string]interface{}, error) {
	// 从 job context 获取分析结果
	videoID := jobContext["video_id"]
	if videoID == "" {
		videoID = "unknown"
	}

	// 获取各部分内容
	summary := e.getStringValue(input, "summary", "摘要生成中...")
	mindmap := e.getStringValue(input, "mindmap", "思维导图生成中...")
	keyPoints := e.getStringValue(input, "key_points", "重点分析生成中...")
	insights := e.getStringValue(input, "insights", "个人认知生成中...")

	// 准备报告数据
	reportData := ReportData{
		Title:       fmt.Sprintf("YouTube 视频分析报告 - %s", videoID),
		VideoID:     videoID,
		VideoURL:    fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		Summary:     e.markdownToHTML(summary),
		Mindmap:     e.markdownToHTML(mindmap),
		KeyPoints:   e.markdownToHTML(keyPoints),
		Insights:    e.markdownToHTML(insights),
	}

	// 生成 HTML 报告
	htmlContent, err := e.generateHTML(reportData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}

	// 保存到文件
	filename := fmt.Sprintf("youtube_analysis_%s_%d.html", videoID, time.Now().Unix())
	filepath := filepath.Join(e.outputDir, filename)

	if err := os.WriteFile(filepath, []byte(htmlContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	logger.Infof("HTML report generated: %s", filepath)

	return map[string]interface{}{
		"report_path": filepath,
		"report_url":  fmt.Sprintf("/reports/%s", filename),
		"filename":    filename,
		"size":        len(htmlContent),
	}, nil
}

// getStringValue 安全获取字符串值
func (e *HTMLReportExecutor) getStringValue(data map[string]interface{}, key, defaultValue string) string {
	if val, ok := data[key].(string); ok && val != "" {
		return val
	}
	return defaultValue
}

// markdownToHTML 简单的 Markdown 转 HTML
func (e *HTMLReportExecutor) markdownToHTML(markdown string) template.HTML {
	// 简单的 Markdown 转换
	html := markdown

	// 标题
	html = strings.ReplaceAll(html, "\n### ", "\n<h3>")
	html = strings.ReplaceAll(html, "\n## ", "\n<h2>")
	html = strings.ReplaceAll(html, "\n# ", "\n<h1>")
	html = strings.ReplaceAll(html, "\n</h3>", "</h3>\n")
	html = strings.ReplaceAll(html, "\n</h2>", "</h2>\n")
	html = strings.ReplaceAll(html, "\n</h1>", "</h1>\n")

	// 给标题加上结束标签（简单处理）
	lines := strings.Split(html, "\n")
	var result []string
	for _, line := range lines {
		if strings.HasPrefix(line, "<h") && !strings.Contains(line, "</h") {
			line = line + strings.Replace(line[:4], "<", "</", 1)
		}
		result = append(result, line)
	}
	html = strings.Join(result, "\n")

	// 粗体
	html = strings.ReplaceAll(html, "**", "<strong>")

	// 列表
	html = strings.ReplaceAll(html, "\n- ", "\n<li>")
	html = strings.ReplaceAll(html, "\n* ", "\n<li>")

	// 段落
	html = strings.ReplaceAll(html, "\n\n", "</p><p>")
	html = "<p>" + html + "</p>"

	return template.HTML(html)
}

// generateHTML 生成 HTML 报告
func (e *HTMLReportExecutor) generateHTML(data ReportData) (string, error) {
	tmpl := template.Must(template.New("report").Parse(reportTemplate))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ReportData 报告数据
type ReportData struct {
	Title       string
	VideoID     string
	VideoURL    string
	GeneratedAt string
	Summary     template.HTML
	Mindmap     template.HTML
	KeyPoints   template.HTML
	Insights    template.HTML
}

// HTML 报告模板
const reportTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #7C3AED 0%, #A78BFA 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }

        .header h1 {
            font-size: 2rem;
            margin-bottom: 10px;
        }

        .header .meta {
            font-size: 0.9rem;
            opacity: 0.9;
        }

        .header .meta a {
            color: white;
            text-decoration: underline;
        }

        .content {
            padding: 40px;
        }

        .section {
            margin-bottom: 40px;
            padding: 30px;
            background: #f9fafb;
            border-radius: 12px;
            border-left: 4px solid #7C3AED;
        }

        .section-title {
            font-size: 1.5rem;
            color: #7C3AED;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
        }

        .section-title::before {
            content: '';
            width: 8px;
            height: 8px;
            background: #7C3AED;
            border-radius: 50%;
            margin-right: 12px;
        }

        .section-content {
            color: #4b5563;
            font-size: 1rem;
        }

        .section-content h1 {
            font-size: 1.8rem;
            color: #1f2937;
            margin: 20px 0 10px;
        }

        .section-content h2 {
            font-size: 1.5rem;
            color: #374151;
            margin: 18px 0 8px;
        }

        .section-content h3 {
            font-size: 1.2rem;
            color: #4b5563;
            margin: 15px 0 8px;
        }

        .section-content p {
            margin: 10px 0;
            line-height: 1.8;
        }

        .section-content li {
            margin-left: 20px;
            margin-bottom: 8px;
            list-style-type: disc;
        }

        .section-content strong {
            color: #7C3AED;
            font-weight: 600;
        }

        .footer {
            background: #f3f4f6;
            padding: 20px;
            text-align: center;
            color: #6b7280;
            font-size: 0.9rem;
        }

        @media print {
            body {
                background: white;
                padding: 0;
            }
            .container {
                box-shadow: none;
            }
        }

        @media (max-width: 768px) {
            .container {
                border-radius: 0;
            }
            .header, .content {
                padding: 20px;
            }
            .section {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📹 YouTube 视频分析报告</h1>
            <div class="meta">
                <p>视频 ID: <a href="{{.VideoURL}}" target="_blank">{{.VideoID}}</a></p>
                <p>生成时间: {{.GeneratedAt}}</p>
                <p>由 GoWorkFlow 自动生成</p>
            </div>
        </div>

        <div class="content">
            <!-- 阅读摘要 -->
            <div class="section">
                <h2 class="section-title">📝 阅读摘要</h2>
                <div class="section-content">
                    {{.Summary}}
                </div>
            </div>

            <!-- 思维导图 -->
            <div class="section">
                <h2 class="section-title">🗺️ 思维导图</h2>
                <div class="section-content">
                    {{.Mindmap}}
                </div>
            </div>

            <!-- 重点分析 -->
            <div class="section">
                <h2 class="section-title">⭐ 重点分析</h2>
                <div class="section-content">
                    {{.KeyPoints}}
                </div>
            </div>

            <!-- 个人认知 -->
            <div class="section">
                <h2 class="section-title">💡 个人认知</h2>
                <div class="section-content">
                    {{.Insights}}
                </div>
            </div>
        </div>

        <div class="footer">
            <p>Powered by GoWorkFlow & BigModel GLM-4-Air</p>
            <p>© 2024 All Rights Reserved</p>
        </div>
    </div>
</body>
</html>`
