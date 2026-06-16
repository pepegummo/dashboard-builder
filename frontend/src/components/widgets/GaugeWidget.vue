<script setup lang="ts">
import { computed } from 'vue'
import type { WidgetElement } from '@/types'

const props = defineProps<{
  title: string
  value: number | null | undefined
  min: number
  max: number
  unit?: string
  elements?: WidgetElement[]
}>()

const percent = computed(() => {
  if (props.value === null || props.value === undefined) return 0
  const range = props.max - props.min || 1
  const pct = ((props.value - props.min) / range) * 100
  return Math.min(100, Math.max(0, pct))
})

const colorClass = computed(() => {
  if (percent.value >= 90) return 'bg-danger'
  if (percent.value >= 70) return 'bg-warning'
  return 'bg-success'
})

function formatNumber(value: number): string {
  return value.toLocaleString(undefined, { maximumFractionDigits: 2 })
}
</script>

<template>
  <div class="card widget-card shadow-sm">
    <div class="card-body" :class="elements?.length ? 'p-0 position-relative' : 'justify-content-center'">
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
            <template v-else-if="el.key === 'value'">
              <span class="gauge-value">
                {{ value === null || value === undefined ? '—' : formatNumber(value) }}
                <small v-if="unit" class="gauge-unit text-muted ms-1">{{ unit }}</small>
              </span>
            </template>
            <template v-else-if="el.key === 'bar'">
              <div class="progress gauge-track w-100">
                <div
                  class="progress-bar"
                  :class="colorClass"
                  role="progressbar"
                  :style="{ width: percent + '%' }"
                  :aria-valuenow="value ?? 0"
                  :aria-valuemin="min"
                  :aria-valuemax="max"
                ></div>
              </div>
            </template>
            <template v-else-if="el.key === 'minmax'">
              <div class="d-flex justify-content-between text-muted small w-100 gauge-minmax">
                <span>{{ formatNumber(min) }}</span>
                <span>{{ formatNumber(max) }}</span>
              </div>
            </template>
          </div>
        </div>
      </template>
      <template v-else>
        <h6 class="card-subtitle text-muted widget-title">{{ title }}</h6>
        <div class="d-flex justify-content-between align-items-baseline mb-2">
          <span class="gauge-value">
            {{ value === null || value === undefined ? '—' : formatNumber(value) }}
            <small v-if="unit" class="gauge-unit text-muted ms-1">{{ unit }}</small>
          </span>
        </div>
        <div class="progress gauge-track">
          <div
            class="progress-bar"
            :class="colorClass"
            role="progressbar"
            :style="{ width: percent + '%' }"
            :aria-valuenow="value ?? 0"
            :aria-valuemin="min"
            :aria-valuemax="max"
          ></div>
        </div>
        <div class="d-flex justify-content-between text-muted small mt-1 gauge-minmax">
          <span>{{ formatNumber(min) }}</span>
          <span>{{ formatNumber(max) }}</span>
        </div>
      </template>
    </div>
  </div>
</template>
