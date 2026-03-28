package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ModelConfig struct {
	APIKey             string            `json:"api_key"`
	APIBase            string            `json:"api_base"`
	Model              string            `json:"model"`
	Timeout            int               `json:"timeout"`
	TranslationOptions map[string]string `json:"translation_options"`
}

type Config struct {
	Port   string                 `json:"port"`
	Models map[string]ModelConfig `json:"models"`
}

var config Config
var defaultModel string

func loadConfig() {
	path := "config.json"
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		path = v
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	if len(config.Models) == 0 {
		log.Fatal("配置文件中没有模型")
	}
	if config.Port == "" {
		config.Port = "8787"
	}
	for k := range config.Models {
		defaultModel = k
		break
	}
}

// ========== 请求/响应结构 ==========

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type TranslatePayload struct {
	Model              string            `json:"model"`
	Messages           []Message         `json:"messages"`
	TranslationOptions map[string]string `json:"translation_options,omitempty"`
}

func makeID() string {
	return "chatcmpl-" + strings.ReplaceAll(uuid.New().String(), "-", "")[:24]
}

func buildResponse(content, model string) map[string]any {
	return map[string]any{
		"id":      makeID(),
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []map[string]any{
			{"index": 0, "message": map[string]string{"role": "assistant", "content": content}, "finish_reason": "stop"},
		},
		"usage": map[string]int{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
	}
}

func extractUserText(messages []Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			switch v := messages[i].Content.(type) {
			case string:
				return v
			case []any:
				var parts []string
				for _, item := range v {
					if m, ok := item.(map[string]any); ok && m["type"] == "text" {
						if t, ok := m["text"].(string); ok {
							parts = append(parts, t)
						}
					}
				}
				return strings.Join(parts, "\n")
			}
		}
	}
	return ""
}

func callTranslate(text, modelID string) (string, error) {
	cfg, ok := config.Models[modelID]
	if !ok {
		cfg = config.Models[defaultModel]
	}

	payload := TranslatePayload{
		Model:              cfg.Model,
		Messages:           []Message{{Role: "user", Content: text}},
		TranslationOptions: cfg.TranslationOptions,
	}
	body, _ := json.Marshal(payload)

	client := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	req, _ := http.NewRequest("POST", strings.TrimRight(cfg.APIBase, "/")+"/chat/completions", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("upstream %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty choices")
	}
	return result.Choices[0].Message.Content, nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func handleModels(w http.ResponseWriter, r *http.Request) {
	var data []map[string]any
	for id := range config.Models {
		data = append(data, map[string]any{
			"id": id, "object": "model", "created": 1700000000, "owned_by": "translate-proxy",
		})
	}
	writeJSON(w, 200, map[string]any{"object": "list", "data": data})
}

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid request"})
		return
	}

	modelID := req.Model
	if modelID == "" {
		modelID = defaultModel
	}

	userText := extractUserText(req.Messages)
	if strings.TrimSpace(userText) == "" {
		writeJSON(w, 200, buildResponse("请输入需要翻译的文本。", modelID))
		return
	}

	result, err := callTranslate(userText, modelID)
	if err != nil {
		if req.Stream {
			streamResult(w, modelID, fmt.Sprintf("翻译失败: %v", err))
			return
		}
		writeJSON(w, 502, map[string]string{"error": fmt.Sprintf("翻译失败: %v", err)})
		return
	}

	if req.Stream {
		streamResult(w, modelID, result)
		return
	}

	writeJSON(w, 200, buildResponse(result, modelID))
}

func streamResult(w http.ResponseWriter, model, content string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	chatID := makeID()
	created := time.Now().Unix()

	chunk, _ := json.Marshal(map[string]any{
		"id": chatID, "object": "chat.completion.chunk", "created": created, "model": model,
		"choices": []map[string]any{{"index": 0, "delta": map[string]string{"role": "assistant", "content": content}, "finish_reason": nil}},
	})
	fmt.Fprintf(w, "data: %s\n\n", chunk)

	done, _ := json.Marshal(map[string]any{
		"id": chatID, "object": "chat.completion.chunk", "created": created, "model": model,
		"choices": []map[string]any{{"index": 0, "delta": map[string]any{}, "finish_reason": "stop"}},
	})
	fmt.Fprintf(w, "data: %s\n\n", done)
	fmt.Fprintf(w, "data: [DONE]\n\n")

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func main() {
	loadConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/models", handleModels)
	mux.HandleFunc("/v1/chat/completions", handleChatCompletions)

	addr := ":" + config.Port
	log.Printf("translate-proxy listening on %s (%d models loaded)", addr, len(config.Models))
	log.Fatal(http.ListenAndServe(addr, mux))
}
