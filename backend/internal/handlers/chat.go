package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ChatHandler forwards chat messages to the Groq API (OpenAI-compatible) using
// a Qwen model. It is stateless: the frontend sends the full conversation
// history plus a compact context string describing the dashboard on screen.
type ChatHandler struct {
	client *http.Client
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{client: &http.Client{Timeout: 30 * time.Second}}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Messages []chatMessage `json:"messages"`
	Context  string        `json:"context"`
}

const chatSystemPrompt = "You explain factory-dashboard widgets and live telemetry to operators. " +
	"Answer using the dashboard context provided below. Be concise and concrete."

// buildGroqBody assembles the OpenAI-compatible request payload. The system
// instruction + dashboard context are prepended as a system message. Split out
// so the assembly can be checked without a network call.
func buildGroqBody(model string, req chatRequest) map[string]any {
	system := chatSystemPrompt
	if strings.TrimSpace(req.Context) != "" {
		system += "\n\nDashboard context:\n" + req.Context
	}
	messages := make([]chatMessage, 0, len(req.Messages)+1)
	messages = append(messages, chatMessage{Role: "system", Content: system})
	messages = append(messages, req.Messages...)
	return map[string]any{
		"model":    model,
		"messages": messages,
		// Qwen3 is a reasoning model; keep its <think> trace out of the reply.
		"reasoning_format": "hidden",
	}
}

func (h *ChatHandler) Send(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		writeError(w, http.StatusServiceUnavailable, "chat is not configured (GROQ_API_KEY unset)")
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Messages) == 0 {
		writeError(w, http.StatusBadRequest, "messages required")
		return
	}

	model := os.Getenv("GROQ_MODEL")
	if model == "" {
		model = "qwen/qwen3-32b"
	}

	body, _ := json.Marshal(buildGroqBody(model, req))
	apiReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost,
		"https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to build request")
		return
	}
	apiReq.Header.Set("content-type", "application/json")
	apiReq.Header.Set("authorization", "Bearer "+apiKey)

	resp, err := h.client.Do(apiReq)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to reach chat service")
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		writeError(w, http.StatusBadGateway, "chat service error: "+string(raw))
		return
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil || len(parsed.Choices) == 0 {
		writeError(w, http.StatusBadGateway, "invalid chat service response")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"reply": parsed.Choices[0].Message.Content})
}
