<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  LinearScale,
  LineElement,
  PointElement,
  Tooltip,
  type ChartData,
  type ChartOptions,
} from 'chart.js'
import type { HistoryPoint } from '@/composables/useTelemetry'
import type { WidgetElement } from '@/types'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Filler)

const props = defineProps<{
  title: string
  history: HistoryPoint[]
  unit?: string
  elements?: WidgetElement[]
}>()

const chartData = computed<ChartData<'line'>>(() => ({
  labels: props.history.map((p) => new Date(p.t).toLocaleTimeString()),
  datasets: [
    {
      label: props.title,
      data: props.history.map((p) => p.v),
      borderColor: '#00d2ff',
      backgroundColor: 'rgba(0, 210, 255, 0.08)',
      tension: 0.35,
      fill: true,
      pointRadius: 0,
      borderWidth: 2,
    },
  ],
}))

const chartOptions = computed<ChartOptions<'line'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  animation: false,
  scales: {
    x: { display: false },
    y: {
      grid: { color: 'rgba(255,255,255,0.06)' },
      ticks: {
        color: '#64748b',
        callback: (value) => (props.unit ? `${value} ${props.unit}` : `${value}`),
      },
      border: { color: 'rgba(255,255,255,0.1)' },
    },
  },
  plugins: {
    legend: { display: false },
  },
}))
</script>

<template>
  <div class="card widget-card db-widget-card shadow-sm">
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
              <h6 class="widget-title db-label mb-0">{{ title }}</h6>
            </template>
            <template v-else-if="el.key === 'chart'">
              <div class="position-absolute" style="inset: 0;">
                <Line v-if="history.length" :data="chartData" :options="chartOptions" />
                <div v-else class="d-flex align-items-center justify-content-center h-100 db-waiting small">
                  Waiting for data…
                </div>
              </div>
            </template>
          </div>
        </div>
      </template>
      <template v-else>
        <h6 class="widget-title db-label">{{ title }}</h6>
        <div class="chart-container">
          <Line v-if="history.length" :data="chartData" :options="chartOptions" />
          <div v-else class="d-flex align-items-center justify-content-center h-100 db-waiting small">
            Waiting for data…
          </div>
        </div>
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

.db-waiting {
  color: var(--db-text-muted);
}
</style>
