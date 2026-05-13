package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	listenAddr    = ":8082"
	minimaxURL    = "https://api.minimax.chat/v1/chat/completions"
	model         = "MiniMax-M2.7"
	maxTokens     = 2000
	maxPerMinute  = 5
	windowSeconds = 60
)

type GenerateRequest struct {
	Topic string `json:"topic"`
}

type GenerateResponse struct {
	Words  []string `json:"words,omitempty"`
	Error  string   `json:"error,omitempty"`
}

type MinimaxMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MinimaxRequest struct {
	Model    string            `json:"model"`
	Messages []MinimaxMessage  `json:"messages"`
	MaxTokens int              `json:"max_tokens"`
}

type MinimaxChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type MinimaxResponse struct {
	Choices []MinimaxChoice `json:"choices"`
	BaseResp struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}

// Simple sliding-window rate limiter per IP
type rateLimiter struct {
	mu       sync.Mutex
	windows  map[string][]int64
}

func newRateLimiter() *rateLimiter {
	rl := &rateLimiter{windows: make(map[string][]int64)}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().Unix()
	cutoff := now - windowSeconds

	entries := rl.windows[ip]
	filtered := entries[:0]
	for _, t := range entries {
		if t > cutoff {
			filtered = append(filtered, t)
		}
	}
	rl.windows[ip] = filtered

	if len(filtered) >= maxPerMinute {
		return false
	}
	rl.windows[ip] = append(filtered, now)
	return true
}

func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		cutoff := time.Now().Unix() - windowSeconds
		for ip, entries := range rl.windows {
			filtered := entries[:0]
			for _, t := range entries {
				if t > cutoff {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				delete(rl.windows, ip)
			} else {
				rl.windows[ip] = filtered
			}
		}
		rl.mu.Unlock()
	}
}

var limiter = newRateLimiter()

func getClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	if fwd := r.Header.Get("X-Real-IP"); fwd != "" {
		return fwd
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func generateWordsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSON(w, http.StatusMethodNotAllowed, GenerateResponse{Error: "use POST"})
		return
	}

	ip := getClientIP(r)
	if !limiter.allow(ip) {
		writeJSON(w, http.StatusTooManyRequests, GenerateResponse{Error: "请求太频繁，请稍后再试"})
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, GenerateResponse{Error: "invalid body"})
		return
	}

	topic := strings.TrimSpace(req.Topic)
	if topic == "" {
		writeJSON(w, http.StatusBadRequest, GenerateResponse{Error: "请提供主题"})
		return
	}

	apiKey := os.Getenv("MINIMAX_API_KEY")
	if apiKey == "" {
		writeJSON(w, http.StatusInternalServerError, GenerateResponse{Error: "server not configured"})
		return
	}

	words, err := callMiniMax(apiKey, topic)
	if err != nil {
		log.Printf("MiniMax error: %v", err)
		writeJSON(w, http.StatusInternalServerError, GenerateResponse{Error: "生成失败，请重试"})
		return
	}

	writeJSON(w, http.StatusOK, GenerateResponse{Words: words})
}

func callMiniMax(apiKey, topic string) ([]string, error) {
	body := MinimaxRequest{
		Model: model,
		Messages: []MinimaxMessage{
			{Role: "user", Content: "Generate 25 Chinese words related to \"" + topic + "\" for a board game. One word per line, 1-4 characters each. Output only the words, no numbering, no explanation."},
		},
		MaxTokens: maxTokens,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", minimaxURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mr MinimaxResponse
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return nil, err
	}

	if mr.BaseResp.StatusCode != 0 {
		return nil, fmt.Errorf("MiniMax API error: %s", mr.BaseResp.StatusMsg)
	}

	if len(mr.Choices) == 0 || mr.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("empty response from MiniMax")
	}

	content := stripThinking(mr.Choices[0].Message.Content)
	lines := strings.Split(content, "\n")
	words := make([]string, 0, 30)
	for _, line := range lines {
		w := strings.TrimSpace(line)
		// Remove common numbering prefixes
		w = strings.TrimLeft(w, "0123456789.、) ）")
		w = strings.TrimSpace(w)
		if w != "" && len([]rune(w)) >= 1 && len([]rune(w)) <= 6 {
			words = append(words, w)
		}
		if len(words) >= 30 {
			break
		}
	}

	if len(words) < 25 {
		return nil, fmt.Errorf("not enough words generated: got %d", len(words))
	}

	return words[:25], nil
}

func stripThinking(s string) string {
	// Remove  ...  blocks from MiniMax thinking output
	for {
		start := strings.Index(s, "<think>")
		if start == -1 {
			break
		}
		end := strings.Index(s, "</think>")
		if end == -1 {
			break
		}
		s = s[:start] + s[end+len("</think>"):]
	}
	return strings.TrimSpace(s)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = listenAddr
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	http.HandleFunc("/api/generate-words", corsMiddleware(generateWordsHandler))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Printf("codenames backend starting on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
