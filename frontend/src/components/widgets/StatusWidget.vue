<script setup lang="ts">
import type { WidgetElement } from '@/types'

const props = defineProps<{
  title: string
  value: string | number | null | undefined
  elements?: WidgetElement[]
}>()

const statusStyles: Record<string, string> = {
  running: 'bg-success',
  idle: 'bg-warning text-dark',
  stopped: 'bg-secondary',
  error: 'bg-danger',
}

function badgeClass(value: string | number | null | undefined): string {
  const key = String(value ?? '').toLowerCase()
  return statusStyles[key] ?? 'bg-secondary'
}
</script>

<template>
  <div class="card widget-card shadow-sm">
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
              <h6 class="card-subtitle text-muted widget-title mb-0">{{ title }}</h6>
            </template>
            <template v-else-if="el.key === 'badge'">
              <span class="badge status-badge text-uppercase" :class="badgeClass(value)">{{ value ?? '—' }}</span>
            </template>
          </div>
        </div>
      </template>
      <template v-else>
        <h6 class="card-subtitle text-muted widget-title">{{ title }}</h6>
        <span class="badge status-badge text-uppercase" :class="badgeClass(value)">
          {{ value ?? '—' }}
        </span>
      </template>
    </div>
  </div>
</template>
