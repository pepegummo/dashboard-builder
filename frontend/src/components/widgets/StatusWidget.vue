<script setup lang="ts">
import type { WidgetElement } from '@/types'

const props = defineProps<{
  title: string
  value: string | number | null | undefined
  elements?: WidgetElement[]
}>()

/* Map known status values to inline bg/glow color pairs */
function statusColor(value: string | number | null | undefined): { bg: string; glow: string } {
  const key = String(value ?? '').toLowerCase()
  if (key === 'running') return { bg: '#2d8470', glow: 'rgba(45, 132, 112, 0.5)' }
  if (key === 'stopped' || key === 'error') return { bg: '#e74c3c', glow: 'rgba(231, 76, 60, 0.5)' }
  if (key === 'idle') return { bg: '#f1970f', glow: 'rgba(241, 151, 15, 0.5)' }
  return { bg: '#64748b', glow: 'rgba(100, 116, 139, 0.4)' }
}
</script>

<template>
  <div class="card widget-card db-widget-card shadow-sm">
    <div class="card-body" :class="elements?.length ? 'p-0 position-relative' : 'text-center justify-content-center'">
      <template v-if="elements?.length">
        <div class="kpi-elements-canvas">
          <div
            v-for="el in elements"
            :key="el.key"
            class="kpi-element-cell"
            :style="{ left: el.x + '%', top: el.y + '%', width: el.w + '%', height: el.h + '%' }"
          >
            <template v-if="el.key === 'title'">
              <h6 class="widget-title db-label mb-0">{{ title }}</h6>
            </template>
            <template v-else-if="el.key === 'badge'">
              <span
                class="status-badge db-status-pill text-uppercase fw-semibold"
                :style="{
                  backgroundColor: statusColor(value).bg,
                  boxShadow: `0 0 12px ${statusColor(value).glow}`,
                }"
              >{{ value ?? '—' }}</span>
            </template>
          </div>
        </div>
      </template>
      <template v-else>
        <h6 class="widget-title db-label">{{ title }}</h6>
        <span
          class="status-badge db-status-pill text-uppercase fw-semibold"
          :style="{
            backgroundColor: statusColor(value).bg,
            boxShadow: `0 0 12px ${statusColor(value).glow}`,
          }"
        >{{ value ?? '—' }}</span>
      </template>
    </div>
  </div>
</template>

<style scoped>
.db-widget-card {
  background-color: var(--db-card-bg) !important;
  border: 1px solid var(--db-border) !important;
  border-radius: 0.625rem !important;
  box-shadow: 0 0 0 1px var(--db-border-glow), 0 4px 24px rgba(0, 0, 0, 0.5) !important;
  color: var(--db-text);
  overflow: hidden;
}

.db-label {
  color: var(--db-text-muted) !important;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  font-size: inherit;
}

.db-status-pill {
  display: inline-block;
  color: #fff;
  border-radius: 999px;
  padding: 0.3em 0.9em;
  letter-spacing: 0.06em;
  transition: background-color 0.4s ease, box-shadow 0.4s ease;
}
</style>
