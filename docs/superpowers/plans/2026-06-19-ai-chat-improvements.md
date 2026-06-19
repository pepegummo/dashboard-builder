# AI Chat Improvements Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the in-app AI chat warmer/adaptive in tone and give it memory of the current conversation session.

**Architecture:** Two independent changes — rewrite the hardcoded system prompt in the Go backend, and accumulate conversation turns in a Vue ref so the full history is sent on every call. No new files, no new dependencies.

**Tech Stack:** Go 1.26 (backend), Vue 3 + TypeScript (frontend), Vite (`npm run build` = type-check)

## Global Constraints

- Backend test command: `go test ./...` from `backend/`
- Frontend type-check: `npm run build` from `frontend/`
- No new dependencies
- No new files

---

### Task 1: Rewrite system prompt in `chat.go`

**Files:**
- Modify: `backend/internal/handlers/chat.go` — replace `chatSystemPrompt` constant
- Test: `backend/internal/handlers/chat_test.go` — existing test already covers it (it checks `strings.Contains(system, chatSystemPrompt)` using the variable, so it auto-validates the new value)

**Interfaces:**
- Produces: `chatSystemPrompt` string constant — same name, new value. All other code unchanged.

- [ ] **Step 1: Open `backend/internal/handlers/chat.go` and replace the `chatSystemPrompt` constant**

  Find (lines 34–36):
  ```go
  const chatSystemPrompt = "You explain factory-dashboard widgets and live telemetry to operators. " +
  	"Answer using the dashboard context provided below. Be concise and concrete. " +
  	"When your answer refers to a specific widget, cite it as [N] (e.g. [0], [2]) using the index numbers from the widget list in the dashboard context."
  ```

  Replace with:
  ```go
  const chatSystemPrompt = "You are a helpful assistant for factory operators monitoring a live dashboard. " +
  	"Answer conversationally and warmly — plain language by default, but match the technical depth of the question. " +
  	"For simple questions be brief (1–2 sentences). When asked to explain or compare, give a clear structured answer. " +
  	"Always cite specific widgets as [N] (e.g. [0], [2]) when your answer refers to them."
  ```

- [ ] **Step 2: Run the existing test to verify it still passes**

  ```bash
  cd backend && go test ./internal/handlers/ -v -run TestBuildGroqBody
  ```

  Expected output:
  ```
  --- PASS: TestBuildGroqBody (0.00s)
  PASS
  ```

  The test checks `strings.Contains(msgs[0].Content, chatSystemPrompt)` using the variable name, so it validates whatever the constant is set to.

- [ ] **Step 3: Commit**

  ```bash
  git add backend/internal/handlers/chat.go
  git commit -m "feat: rewrite AI chat system prompt for warmer adaptive tone"
  ```

---

### Task 2: Add invisible conversation history to `ExplorePanel.vue`

**Files:**
- Modify: `frontend/src/components/ExplorePanel.vue`

**Interfaces:**
- Consumes: `ChatMessage` from `@/types` — `{ role: 'user' | 'assistant', content: string }`
- Consumes: `api.chat(messages: ChatMessage[], context: string)` from `@/services/api`
- Produces: no interface changes — props and emits are unchanged

- [ ] **Step 1: Add `ChatMessage` import to the script block**

  In `frontend/src/components/ExplorePanel.vue`, find the existing import line:
  ```ts
  import { ref, watch } from 'vue'
  ```

  Replace with:
  ```ts
  import { ref, watch } from 'vue'
  import type { ChatMessage } from '@/types'
  ```

- [ ] **Step 2: Add the `messages` ref and update `submit()`**

  Find the existing refs and `submit` function (lines 8–34):
  ```ts
  const question = ref('')

  watch(
    () => props.prefill,
    (v) => { if (v) question.value = v },
  )
  const answer = ref('')
  const loading = ref(false)
  const errorMsg = ref('')

  async function submit() {
    const q = question.value.trim()
    if (!q || loading.value) return
    loading.value = true
    errorMsg.value = ''
    question.value = ''
    try {
      const { reply } = await api.chat([{ role: 'user', content: q }], props.context)
      answer.value = reply
      const indices = [...reply.matchAll(/\[(\d+)\]/g)].map((m) => parseInt(m[1]))
      emit('highlight', [...new Set(indices)])
    } catch (err) {
      errorMsg.value = apiErrorMessage(err, 'Failed to get answer')
    } finally {
      loading.value = false
    }
  }
  ```

  Replace with:
  ```ts
  const question = ref('')
  const messages = ref<ChatMessage[]>([])

  watch(
    () => props.prefill,
    (v) => { if (v) question.value = v },
  )
  const answer = ref('')
  const loading = ref(false)
  const errorMsg = ref('')

  watch(
    () => props.context,
    () => { messages.value = []; answer.value = '' },
  )

  async function submit() {
    const q = question.value.trim()
    if (!q || loading.value) return
    loading.value = true
    errorMsg.value = ''
    question.value = ''
    messages.value.push({ role: 'user', content: q })
    try {
      const { reply } = await api.chat(messages.value, props.context)
      messages.value.push({ role: 'assistant', content: reply })
      answer.value = reply
      const indices = [...reply.matchAll(/\[(\d+)\]/g)].map((m) => parseInt(m[1]))
      emit('highlight', [...new Set(indices)])
    } catch (err) {
      messages.value.pop() // remove the user message that failed
      errorMsg.value = apiErrorMessage(err, 'Failed to get answer')
    } finally {
      loading.value = false
    }
  }
  ```

  Key changes:
  - `messages` ref accumulates the full conversation turn by turn
  - `watch(props.context)` resets both `messages` and `answer` when the dashboard or machine page changes
  - On error, pop the user message so a retry doesn't double-push it
  - `api.chat` now receives the full `messages.value` array instead of a single-element array

- [ ] **Step 3: Type-check the frontend**

  ```bash
  cd frontend && npm run build
  ```

  Expected: build completes with no TypeScript errors.

- [ ] **Step 4: Manual smoke test**

  1. Start the backend: `cd backend && go run ./cmd/server`
  2. Start the frontend: `cd frontend && npm run dev`
  3. Open `http://localhost:5174`, go to Explore, pick a dashboard
  4. Ask: "what is the temperature widget?"
  5. Ask a follow-up: "is that normal?" — AI should reference the temperature context from the previous turn without you repeating it
  6. Switch to a different dashboard — the panel should clear and the AI should have no memory of the previous conversation

- [ ] **Step 5: Commit**

  ```bash
  git add frontend/src/components/ExplorePanel.vue
  git commit -m "feat: add conversation history to AI chat panel"
  ```
