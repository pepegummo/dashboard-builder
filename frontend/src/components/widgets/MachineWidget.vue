<script setup lang="ts">
import type { Machine, WidgetElement } from '@/types'

defineProps<{
  title: string
  machine: Machine | null
  elements?: WidgetElement[]
}>()

const statusBadge: Record<string, string> = {
  running: 'bg-success',
  idle: 'bg-warning text-dark',
  stopped: 'bg-secondary',
  error: 'bg-danger',
}

function badgeClass(status: string): string {
  return statusBadge[status.toLowerCase()] ?? 'bg-secondary'
}
</script>

<template>
  <div class="card widget-card db-widget-card shadow-sm h-100">
    <div class="card-body d-flex flex-column justify-content-center gap-2">
      <h6 class="widget-title db-label mb-0">{{ title }}</h6>
      <div v-if="!machine" class="db-waiting small">Waiting for machine…</div>
      <template v-else>
        <div class="machine-name fw-bold">{{ machine.name }}</div>
        <div class="machine-type small db-text-muted">{{ machine.type }}</div>
        <div>
          <span class="badge" :class="badgeClass(machine.status)">{{ machine.status }}</span>
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
.machine-name {
  font-size: 1.1rem;
}
.machine-type {
  color: var(--db-text-muted);
}
</style>
