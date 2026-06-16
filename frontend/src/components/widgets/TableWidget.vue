<script setup lang="ts">
import type { WidgetElement } from '@/types'
import { useCatalogStore } from '@/stores/catalog.store'

defineProps<{
  title: string
  readings: Record<string, number | string> | null
  elements?: WidgetElement[]
}>()

const catalog = useCatalogStore()

function formatValue(key: string, value: number | string | undefined): string {
  if (value === undefined) return '—'
  if (typeof value === 'number') {
    const formatted = value.toLocaleString(undefined, { maximumFractionDigits: 2 })
    const unit = catalog.metricUnit(key)
    return unit ? `${formatted} ${unit}` : formatted
  }
  return String(value)
}
</script>

<template>
  <div class="card widget-card shadow-sm">
    <div class="card-body" :class="elements?.length ? 'p-0 position-relative' : ''">
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
            <template v-else-if="el.key === 'table'">
              <div class="position-absolute overflow-auto" style="inset: 0;">
                <div v-if="!readings" class="text-muted small p-2">Waiting for data…</div>
                <table v-else class="table table-sm mb-0 widget-table">
                  <tbody>
                    <tr v-for="metric in catalog.metrics" :key="metric.key">
                      <td class="text-muted">{{ metric.label }}</td>
                      <td class="text-end fw-semibold">{{ formatValue(metric.key, readings[metric.key]) }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </template>
          </div>
        </div>
      </template>
      <template v-else>
        <h6 class="card-subtitle text-muted widget-title">{{ title }}</h6>
        <div v-if="!readings" class="text-muted small">Waiting for data…</div>
        <table v-else class="table table-sm mb-0 widget-table">
          <tbody>
            <tr v-for="metric in catalog.metrics" :key="metric.key">
              <td class="text-muted">{{ metric.label }}</td>
              <td class="text-end fw-semibold">{{ formatValue(metric.key, readings[metric.key]) }}</td>
            </tr>
          </tbody>
        </table>
      </template>
    </div>
  </div>
</template>
