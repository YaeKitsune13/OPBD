<script setup lang="ts">
import { onMounted, onUnmounted, ref, reactive, watch } from 'vue'
import { useToast } from '../../../utils/useToast';
const { showToast } = useToast();

const canvasRef = ref(null)
const pets = ref([])
const weightData = ref([])
const historyForm = reactive({ petId: '' })

let animId = null

async function getPetWeightHistory(petId) {
  const token = localStorage.getItem('token')
  try {
    const weightRes = await fetch(`/api/weight/pet/${petId}`, {
      headers: { "Authorization": `Bearer ${token}` }
    });
    if (weightRes.ok) weightData.value = await weightRes.json();
  } catch (e) {
    showToast("Ошибка при загрузке данных", "error")
  }
}

async function getUserPets() {
  const token = localStorage.getItem('token')
  const userRaw = localStorage.getItem('user')
  if (!userRaw) return;

  const userData = JSON.parse(userRaw)
  const userId = userData.user_id || userData.id
  try {
    const petsRes = await fetch(`/api/pets/owner/${userId}`, {
      headers: { "Authorization": `Bearer ${token}` }
    });
    if (petsRes.ok) pets.value = await petsRes.json();
    console.log(pets.value);
  } catch (e) {
    showToast("Ошибка при загрузке данных", "error")
  }
}

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

  ctx.clearRect(0, 0, W, H)

  const data = weightData.value
  if (!data || data.length === 0) return

  const values = data.map(d => d.value)
  const minV = Math.min(...values)
  const maxV = Math.max(...values)
  const range = maxV - minV || 1

  const pad = { left: 40, right: 20, top: 24, bottom: 30 }
  const chartW = W - pad.left - pad.right
  const chartH = H - pad.top - pad.bottom

  const xOf = (i) => pad.left + (i / (data.length - 1 || 1)) * chartW
  const yOf = (v) => pad.top + chartH - ((v - minV) / range) * chartH

  // Строим полный path один раз
  const fullPath = new Path2D()
  data.forEach((d, i) => {
    i === 0 ? fullPath.moveTo(xOf(i), yOf(d.value)) : fullPath.lineTo(xOf(i), yOf(d.value))
  })

  // Общая длина линии (приблизительно через сумму расстояний)
  let totalLength = 0
  for (let i = 1; i < data.length; i++) {
    const dx = xOf(i) - xOf(i - 1)
    const dy = yOf(data[i].value) - yOf(data[i - 1].value)
    totalLength += Math.sqrt(dx * dx + dy * dy)
  }

  const DURATION = 900  // мс
  let startTime = null
  animId = null

  function animate(ts) {
    if (!startTime) startTime = ts
    const elapsed = ts - startTime
    const progress = Math.min(elapsed / DURATION, 1)

    // Easing: ease-in-out
    const t = progress < 0.5
      ? 2 * progress * progress
      : 1 - Math.pow(-2 * progress + 2, 2) / 2

    ctx.clearRect(0, 0, W, H)

    // --- Фоновые горизонтальные линии ---
    ctx.strokeStyle = 'rgba(255,255,255,0.06)'
    ctx.lineWidth = 0.5
    for (let i = 0; i <= 4; i++) {
      const y = pad.top + (chartH / 4) * i
      ctx.beginPath()
      ctx.moveTo(pad.left, y)
      ctx.lineTo(W - pad.right, y)
      ctx.stroke()
    }

    // --- Ось Y (значения) ---
    ctx.fillStyle = 'rgba(255,255,255,0.35)'
    ctx.font = '10px sans-serif'
    ctx.textAlign = 'right'
    for (let i = 0; i <= 2; i++) {
      const v = minV + (range / 2) * i
      const y = yOf(v)
      ctx.fillText(v.toFixed(1), pad.left - 6, y + 3)
    }

    // --- Линия с clipping по прогрессу ---
    ctx.save()
    ctx.beginPath()
    ctx.rect(pad.left, 0, chartW * t, H)
    ctx.clip()

    // Тень под линией
    ctx.save()
    const grad = ctx.createLinearGradient(0, pad.top, 0, pad.top + chartH)
    grad.addColorStop(0, 'rgba(74, 222, 128, 0.18)')
    grad.addColorStop(1, 'rgba(74, 222, 128, 0)')
    const fillPath = new Path2D(fullPath)
    fillPath.lineTo(xOf(data.length - 1), pad.top + chartH)
    fillPath.lineTo(xOf(0), pad.top + chartH)
    fillPath.closePath()
    ctx.fillStyle = grad
    ctx.fill(fillPath)
    ctx.restore()

    // Сама линия
    ctx.strokeStyle = '#4ade80'
    ctx.lineWidth = 2
    ctx.lineJoin = 'round'
    ctx.lineCap = 'round'
    ctx.stroke(fullPath)

    ctx.restore()

    // --- Точки и подписи — появляются по мере прохождения линии ---
    data.forEach((d, i) => {
      const pointProgress = i / (data.length - 1 || 1)
      if (t < pointProgress) return  // ещё не дошли

      // Точка "вырастает"
      const gap = data.length > 1 ? 1 / (data.length - 1) : 1
      const localT = Math.min((t - pointProgress) / Math.min(gap * 0.8, 0.18), 1)
      const ease = 1 - Math.pow(1 - localT, 3)
      const r = 4 * ease

      const x = xOf(i)
      const y = yOf(d.value)

      // Внешнее кольцо
      ctx.fillStyle = 'rgba(74, 222, 128, 0.2)'
      ctx.beginPath()
      ctx.arc(x, y, r * 2, 0, Math.PI * 2)
      ctx.fill()

      // Точка
      ctx.fillStyle = '#4ade80'
      ctx.beginPath()
      ctx.arc(x, y, r, 0, Math.PI * 2)
      ctx.fill()

      // Центр белый
      ctx.fillStyle = '#111'
      ctx.beginPath()
      ctx.arc(x, y, r * 0.4, 0, Math.PI * 2)
      ctx.fill()

      // Значение сверху
      ctx.globalAlpha = ease
      ctx.fillStyle = '#ffffff'
      ctx.font = '500 11px sans-serif'
      ctx.textAlign = 'center'
      ctx.fillText(d.value + ' кг', x, y - 10)

      // Метка снизу
      ctx.fillStyle = 'rgba(255,255,255,0.45)'
      ctx.font = '10px sans-serif'
      ctx.fillText(d.label, x, H - 8)
      ctx.globalAlpha = 1
    })

    if (progress < 1) {
      animId = requestAnimationFrame(animate)
    }
  }

  if (animId) cancelAnimationFrame(animId)
  animId = requestAnimationFrame(animate)
}

function clearChart() {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  ctx.clearRect(0, 0, canvas.width, canvas.height)
}

onMounted(() => {
  drawChart()
  window.addEventListener('resize', drawChart)
})

onUnmounted(() => {
  window.removeEventListener('resize', drawChart)
})

watch(() => historyForm.petId, async (newId) => {
  if (animId) {
    cancelAnimationFrame(animId)
    animId = null
  }

  if (newId) {
    await getPetWeightHistory(newId)
    drawChart()
  } else {
    weightData.value = []
    clearChart()
  }
})

getUserPets()
</script>

<template>
  <div class="page-content">
    <div class="page-header">
      <div class="page-title">Динамика веса</div>
      <select class="pet-select" v-model="historyForm.petId">
        <option value="">— Выберите питомца —</option>
        <option v-for="pet in pets" :key="pet.petId" :value="pet.petId">
          {{ pet.avatar || '🐾' }} {{ pet.name }}
        </option>
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
            <tr v-for="row in weightData" :key="row.date">
              <td class="mono">{{ row.date }}</td>
              <td class="td-main">{{ row.value }}</td>
              <td>{{ row.doctorName }}</td>
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
