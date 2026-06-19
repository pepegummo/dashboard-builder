package handlers

import (
	"strings"
	"testing"
)

func TestBuildGroqBody(t *testing.T) {
	req := chatRequest{
		Messages: []chatMessage{{Role: "user", Content: "hi"}},
		Context:  "Dashboard: Line A",
	}
	body := buildGroqBody("qwen/qwen3-32b", req)

	if body["model"] != "qwen/qwen3-32b" {
		t.Fatalf("model not passed through: %v", body["model"])
	}
	msgs, ok := body["messages"].([]chatMessage)
	if !ok || len(msgs) != 2 {
		t.Fatalf("expected system + user message, got: %v", body["messages"])
	}
	if msgs[0].Role != "system" || !strings.Contains(msgs[0].Content, chatSystemPrompt) || !strings.Contains(msgs[0].Content, "Line A") {
		t.Fatalf("system message missing base or context: %+v", msgs[0])
	}
	if msgs[1].Role != "user" || msgs[1].Content != "hi" {
		t.Fatalf("user message not passed through: %+v", msgs[1])
	}

	// Empty context must not append the "Dashboard context:" header.
	bare := buildGroqBody("m", chatRequest{Messages: req.Messages})
	bareMsgs := bare["messages"].([]chatMessage)
	if strings.Contains(bareMsgs[0].Content, "Dashboard context") {
		t.Fatal("empty context should not add context header")
	}
}
