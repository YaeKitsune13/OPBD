<script setup>
import { ref } from 'vue'

const stats = ref([
  { label: 'Визитов за месяц', value: '147', change: '+18 vs март', trend: 'up' },
  { label: 'Выручка', value: '248k', change: '+12%', trend: 'up', mono: true },
  { label: 'Активных клиентов', value: '89', change: '+5', trend: 'up' },
  { label: 'Ср. чек', value: '1 687', change: 'руб.', trend: '', mono: true },
])

const popularServices = ref([
  { name: 'Первичный осмотр', count: 68, revenue: '57 800' },
  { name: 'Вакцинация', count: 41, revenue: '34 850' },
  { name: 'УЗИ брюшной полости', count: 22, revenue: '33 000' },
])
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Аналитическая сводка</div>
        <div class="page-sub">Апрель 2026</div>
      </div>
      <button class="btn btn-ghost" @click="alert('Экспорт запущен')">⬇ Экспорт</button>
    </div>

    <div class="stats-grid">
      <div v-for="s in stats" :key="s.label" class="stat-card">
        <div class="stat-label">{{ s.label }}</div>
        <div class="stat-value" :class="{ 'text-accent': s.trend === 'up', mono: s.mono }">
          {{ s.value }}
        </div>
        <div class="stat-change" :class="s.trend">{{ s.change }}</div>
      </div>
    </div>

    <div class="grid-2">
      <div class="card">
        <div class="card-header"><span class="card-title">Популярные услуги</span></div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Услуга</th>
                <th>Раз</th>
                <th>Выручка</th>
              </tr>
            </thead>
            <tbody>
                <!-- AnalyticsPage.vue -->
                <tbody>
                  <!-- Было ser, стал ser в ключе тоже -->
                  <tr v-for="ser in popularServices" :key="ser.name">
                    <td class="td-main">{{ ser.name }}</td>
                    <td class="mono">{{ ser.count }}</td>
                    <td class="mono text-accent">{{ ser.revenue }} руб.</td>
                  </tr>
                </tbody>
            </tbody>
          </table>
        </div>
      </div>

      <div class="card">
        <div class="card-header"><span class="card-title">Нагрузка врачей</span></div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Врач</th>
                <th>Приёмов</th>
                <th>Статус</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td class="td-main">Кузнецов А.В.</td>
                <td class="mono">52</td>
                <td><span class="badge badge-confirmed">Высокая</span></td>
              </tr>
              <tr>
                <td class="td-main">Попова М.С.</td>
                <td class="mono">38</td>
                <td><span class="badge badge-waiting">Средняя</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
