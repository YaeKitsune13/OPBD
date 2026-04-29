<script setup>
import { onMounted, onUnmounted, ref } from 'vue'

const canvasRef = ref(null)
const weightData = [
  { label: 'ноя', value: 3.85 },
  { label: 'янв', value: 3.8 },
  { label: 'фев', value: 4.0 },
  { label: 'апр', value: 4.2 },
]

function drawChart() {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  const dpr = window.devicePixelRatio || 1
  const rect = canvas.parentElement.getBoundingClientRect()
  const W = rect.width
  const H = 180

  canvas.width = W * dpr
  canvas.height = H * dpr
  canvas.style.width = W + 'px'
  canvas.style.height = H + 'px'
  ctx.scale(dpr, dpr)

  // ... (весь твой код из drawWeightChart, но используем ctx и переменные W, H) ...
  // Используй массив weightData из этого скрипта

  // Короткий пример для проверки:
  ctx.strokeStyle = '#4ade80'
  ctx.lineWidth = 2
  ctx.beginPath()
  ctx.moveTo(0, H / 2)
  ctx.lineTo(W, H / 2)
  ctx.stroke()
}

onMounted(() => {
  drawChart()
  window.addEventListener('resize', drawChart)
})

onUnmounted(() => {
  window.removeEventListener('resize', drawChart)
})
</script>

<template>
  <div class="page-content">
    <div class="page-header">
      <div class="page-title">Динамика веса</div>
      <select class="pet-select">
        <option>🐱 Барсик</option>
      </select>
    </div>

    <div class="card">
      <div class="card-header">
        <span class="card-title">Барсик — История веса</span>
      </div>
      <div class="card-body">
        <div class="chart-area">
          <canvas ref="canvasRef" class="chart-canvas"></canvas>
        </div>
      </div>
    </div>

    <!-- Таблица -->
    <div class="card mt-12">
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Дата</th>
              <th>Вес (кг)</th>
              <th>Врач</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in 3" :key="row">
              <td class="mono">10.04.2026</td>
              <td class="td-main">4.20</td>
              <td>Кузнецов А.В.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chart-area {
  height: 180px;
  position: relative;
}
.chart-canvas {
  width: 100%;
  height: 100%;
  display: block;
}
.pet-select {
  width: auto;
  padding: 7px 12px;
  background: var(--surface2);
  color: white;
  border: 1px solid var(--border);
  border-radius: 4px;
}
</style>
