import type { WidgetElement, WidgetType } from '@/types'

export const DEFAULT_ELEMENTS: Partial<Record<WidgetType, WidgetElement[]>> = {
  status: [
    { key: 'title', x: 0, y: 0,  w: 100, h: 30 },
    { key: 'badge', x: 0, y: 30, w: 100, h: 70 },
  ],
  kpi: [
    { key: 'title', x: 0, y: 0,  w: 100, h: 25 },
    { key: 'value', x: 0, y: 25, w: 100, h: 50 },
    { key: 'unit',  x: 0, y: 75, w: 100, h: 25 },
  ],
  gauge: [
    { key: 'title',  x: 0, y: 0,  w: 100, h: 20 },
    { key: 'value',  x: 0, y: 20, w: 100, h: 35 },
    { key: 'bar',    x: 0, y: 55, w: 100, h: 25 },
    { key: 'minmax', x: 0, y: 80, w: 100, h: 20 },
  ],
  line: [
    { key: 'title', x: 0, y: 0,  w: 100, h: 20 },
    { key: 'chart', x: 0, y: 20, w: 100, h: 80 },
  ],
  bar: [
    { key: 'title', x: 0, y: 0,  w: 100, h: 20 },
    { key: 'chart', x: 0, y: 20, w: 100, h: 80 },
  ],
  table: [
    { key: 'title', x: 0, y: 0,  w: 100, h: 15 },
    { key: 'table', x: 0, y: 15, w: 100, h: 85 },
  ],
}
