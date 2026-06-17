<script setup lang="ts">
import type { Widget, WidgetElement } from '@/types'
import WidgetElementEditor from './WidgetElementEditor.vue'

const props = defineProps<{
  widget: Widget
  isSelected: boolean
  isElementEditing: boolean
  activeElementKey?: string
  hoveredElementKey?: string
}>()

const emit = defineEmits<{
  select: []
  remove: []
  'update-elements': [WidgetElement[]]
  'select-element': [string]
}>()

function handleBodyClick() {
  emit('select')
}
</script>

<template>
  <div
    class="widget-editor-shell"
    :class="{ 'is-selected': isSelected }"
    @click="handleBodyClick"
  >
    <div class="widget-shell-toolbar">
      <span class="drag-handle"><i class="bi bi-grip-vertical"></i></span>
      <span class="widget-shell-title">{{ widget.title }}</span>
      <button
        type="button"
        class="btn-shell-config"
        title="Configure widget"
        @click.stop="$emit('select')"
      >
        <i class="bi bi-gear-fill"></i>
      </button>
      <button
        type="button"
        class="btn-shell-remove"
        title="Remove widget"
        @click.stop="$emit('remove')"
      >
        <i class="bi bi-x-lg"></i>
      </button>
    </div>
    <div class="widget-shell-resize-corner"></div>
    <slot />
    <WidgetElementEditor
      v-if="isElementEditing"
      :widget="widget"
      :activeKey="activeElementKey"
      :hoveredKey="hoveredElementKey"
      @update:elements="emit('update-elements', $event)"
      @select-element="emit('select-element', $event)"
    />
  </div>
</template>
