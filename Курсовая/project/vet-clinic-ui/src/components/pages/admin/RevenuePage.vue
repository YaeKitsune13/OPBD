<script setup>
import { ref } from 'vue'

const dateFrom = ref('2026-04-01')
const dateTo = ref('2026-04-28')

const dailyRevenue = ref([
  { date: '28.04.2026', visits: 8, services: '9 400', meds: '3 200', total: '12 600' },
  { date: '27.04.2026', visits: 6, services: '7 100', meds: '2 850', total: '9 950' },
  { date: '26.04.2026', visits: 5, services: '5 800', meds: '1 900', total: '7 700' },
])
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Выручка</div>
        <div class="page-sub">Финансовые отчёты по периодам</div>
      </div>
      <div class="row">
        <input type="date" v-model="dateFrom" style="width: auto" />
        <input type="date" v-model="dateTo" style="width: auto" />
        <button class="btn btn-primary">Сформировать</button>
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">Итого за период</div>
        <div class="stat-value mono">248 460</div>
        <div class="stat-change">руб.</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Услуги</div>
        <div class="stat-value mono">164 200</div>
        <div class="stat-change up">66% от выручки</div>
      </div>
    </div>

    <div class="card">
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Дата</th>
              <th>Визитов</th>
              <th>Услуги</th>
              <th>Медикаменты</th>
              <th>Итого</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="day in dailyRevenue" :key="day.date">
              <td class="mono">{{ day.date }}</td>
              <td>{{ day.visits }}</td>
              <td>{{ day.services }} руб.</td>
              <td>{{ day.meds }} руб.</td>
              <td class="mono text-accent">{{ day.total }} руб.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
