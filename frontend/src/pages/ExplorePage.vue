<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useDashboardStore } from '@/stores/dashboard.store'
import { useCatalogStore } from '@/stores/catalog.store'
import { useTelemetry } from '@/composables/useTelemetry'
import { api, apiErrorMessage } from '@/services/api'
import { DEFAULT_ELEMENTS } from '@/utils/widgetElements'
import WidgetRenderer from '@/components/widgets/WidgetRenderer.vue'
import type { ChatMessage, Dashboard, Widget } from '@/types'

const dashboardStore = useDashboardStore()
const catalog = useCatalogStore()
const { reading, history, start } = useTelemetry()

dashboardStore.fetchDashboards()
catalog.fetchMetrics()

const selectedId = ref('')
const dashboard = ref<Dashboard | null>(null)
const loading = ref(false)
const error = ref('')
const activePage = ref(0)

const pages = computed(() => dashboard.value?.pages ?? [])
const currentMachine = computed(() => pages.value[activePage.value]?.machine ?? null)
const widgets = computed(() => dashboard.value?.template?.widgets ?? [])
const template = computed(() => dashboard.value?.template ?? null)

const canvasStyle = computed(() => {
  if (!template.value) return undefined
  const ratio = template.value.width / template.value.height
  return {
    aspectRatio: `${template.value.width} / ${template.value.height}`,
    width: `min(100%, calc(58vh * ${ratio}))`,
    maxHeight: '58vh',
  }
})

function elementsFor(w: Widget) {
  return w.elements ?? DEFAULT_ELEMENTS[w.type] ?? []
}

async function loadDashboard() {
  if (!selectedId.value) {
    dashboard.value = null
    return
  }
  loading.value = true
  error.value = ''
  activePage.value = 0
  focused.value = null
  try {
    dashboard.value = await api.getDashboard(selectedId.value)
  } catch (err) {
    error.value = apiErrorMessage(err, 'Failed to load dashboard')
  } finally {
    loading.value = false
  }
}

watch(
  currentMachine,
  (machine) => {
    if (machine) start(machine.id)
  },
  { immediate: true },
)

// --- Element inspect overlay (read-only: hover highlight + click to ask) ---
const hovered = ref<{ widgetId: string; key: string } | null>(null)
const focused = ref<{ widgetId: string; key: string } | null>(null)

function onElementClick(w: Widget, key: string) {
  focused.value = { widgetId: w.id, key }
  send(`What is the "${key}" element on the "${w.title}" widget for?`)
}

// --- Chat ---
const messages = ref<ChatMessage[]>([])
const input = ref('')
const sending = ref(false)
const chatLog = ref<HTMLElement | null>(null)

function metricLine(metricKey: string): string {
  const m = catalog.metrics.find((x) => x.key === metricKey)
  if (!m) return metricKey || '(none)'
  const range = m.unit ? `${m.min}–${m.max} ${m.unit}` : `${m.min}–${m.max}`
  return `${m.label} (${m.key}), range ${range}`
}

function buildContext(): string {
  if (!dashboard.value) return ''
  const lines: string[] = []
  lines.push(`Dashboard: ${dashboard.value.name}`)
  if (dashboard.value.factory) lines.push(`Factory: ${dashboard.value.factory.name}`)
  if (currentMachine.value) lines.push(`Current machine: ${currentMachine.value.name}`)
  lines.push('Widgets:')
  for (const w of widgets.value) {
    const r = reading.value?.[w.metricKey]
    const live = r !== undefined ? `, live value: ${r}` : ''
    lines.push(`- ${w.type} "${w.title}" → ${metricLine(w.metricKey)}${live}`)
  }
  if (focused.value) {
    const w = widgets.value.find((x) => x.id === focused.value!.widgetId)
    if (w) lines.push(`Focused element: "${focused.value.key}" of widget "${w.title}" (${w.type}, ${w.metricKey})`)
  }
  return lines.join('\n')
}

async function send(text: string) {
  const content = text.trim()
  if (!content || sending.value || !dashboard.value) return
  messages.value.push({ role: 'user', content })
  sending.value = true
  await scrollChat()
  try {
    const { reply } = await api.chat(messages.value, buildContext())
    messages.value.push({ role: 'assistant', content: reply })
  } catch (err) {
    messages.value.push({ role: 'assistant', content: apiErrorMessage(err, 'Chat failed') })
  } finally {
    sending.value = false
    await scrollChat()
  }
}

function onSubmit() {
  const text = input.value
  input.value = ''
  send(text)
}

async function scrollChat() {
  await nextTick()
  if (chatLog.value) chatLog.value.scrollTop = chatLog.value.scrollHeight
}
</script>

<template>
  <div class="explore-page">
    <div class="d-flex align-items-center gap-2 mb-3">
      <label class="fw-semibold mb-0">Dashboard</label>
      <select
        v-model="selectedId"
        class="form-select w-auto"
        :disabled="dashboardStore.loading"
        @change="loadDashboard"
      >
        <option value="">Select a dashboard…</option>
        <option v-for="d in dashboardStore.dashboards" :key="d.id" :value="d.id">
          {{ d.name }}<template v-if="d.factory"> — {{ d.factory.name }}</template>
        </option>
      </select>
    </div>

    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status"></div>
    </div>
    <div v-else-if="error" class="alert alert-danger">{{ error }}</div>

    <div v-else-if="dashboard" class="explore-canvas-wrap">
      <div class="dashboard-canvas" :style="canvasStyle">
        <div
          v-for="w in widgets"
          :key="w.id"
          class="dashboard-canvas-item"
          :style="{ left: `${w.x}%`, top: `${w.y}%`, width: `${w.w}%`, height: `${w.h}%` }"
        >
          <WidgetRenderer :widget="w" :readings="reading" :history="history" :machine="currentMachine" />
          <div class="inspect-overlay">
            <div
              v-for="el in elementsFor(w)"
              :key="el.key"
              class="element-handle"
              :class="{
                'is-active': focused?.widgetId === w.id && focused?.key === el.key,
                'is-hovered': hovered?.widgetId === w.id && hovered?.key === el.key,
              }"
              :style="{ left: `${el.x}%`, top: `${el.y}%`, width: `${el.w}%`, height: `${el.h}%` }"
              @mouseenter="hovered = { widgetId: w.id, key: el.key }"
              @mouseleave="hovered = null"
              @click="onElementClick(w, el.key)"
            >
              <span class="element-label">{{ el.key }}</span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="pages.length > 1" class="d-flex justify-content-center gap-2 mt-2">
        <button
          v-for="(p, i) in pages"
          :key="p.id"
          type="button"
          class="page-dot"
          :class="{ active: activePage === i }"
          :aria-label="`Show ${p.machine?.name}`"
          @click="activePage = i"
        ></button>
      </div>
    </div>

    <p v-else class="text-secondary py-4">
      Pick a dashboard to view it and ask about its widgets.
    </p>

    <!-- Chat -->
    <div class="chat-panel mt-3">
      <div ref="chatLog" class="chat-log">
        <p v-if="!messages.length" class="text-secondary small m-0">
          Hover an element to see what it's called, click it to ask about it, or type a question below.
        </p>
        <div
          v-for="(m, i) in messages"
          :key="i"
          class="chat-msg"
          :class="m.role === 'user' ? 'chat-msg-user' : 'chat-msg-bot'"
        >
          {{ m.content }}
        </div>
        <div v-if="sending" class="chat-msg chat-msg-bot text-secondary">…</div>
      </div>
      <form class="d-flex gap-2 mt-2" @submit.prevent="onSubmit">
        <input
          v-model="input"
          class="form-control"
          placeholder="Ask about this dashboard…"
          :disabled="!dashboard || sending"
        />
        <button class="btn btn-primary" type="submit" :disabled="!dashboard || sending || !input.trim()">
          Send
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.inspect-overlay {
  position: absolute;
  inset: 1px;
  z-index: 20;
}
.chat-log {
  height: 30vh;
  overflow-y: auto;
  border: 1px solid var(--bs-border-color);
  border-radius: 0.5rem;
  padding: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.chat-msg {
  max-width: 80%;
  padding: 0.5rem 0.75rem;
  border-radius: 0.75rem;
  white-space: pre-wrap;
}
.chat-msg-user {
  align-self: flex-end;
  background-color: var(--bs-primary);
  color: #fff;
}
.chat-msg-bot {
  align-self: flex-start;
  background-color: var(--bs-secondary-bg);
}
</style>
