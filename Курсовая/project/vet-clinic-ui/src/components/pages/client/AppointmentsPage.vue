<script setup>
import { ref } from 'vue'

const appointments = ref([
  {
    id: 1042,
    pet: '🐱 Барсик',
    doctor: 'Кузнецов А.В.',
    spec: 'Терапевт',
    date: '02.05.2026',
    time: '10:30',
    status: 'waiting',
  },
  {
    id: 1038,
    pet: '🐶 Рыжик',
    doctor: 'Попова М.С.',
    spec: 'Хирург',
    date: '10.05.2026',
    time: '14:00',
    status: 'confirmed',
  },
])

function cancel(id) {
  alert('Запись #' + id + ' будет отменена')
}
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Записи на приём</div>
        <div class="page-sub">Статусы и история заявок</div>
      </div>
      <button class="btn btn-primary" @click="$emit('navigate', 'book')">✚ Новая запись</button>
    </div>

    <div class="card">
      <div class="card-header">
        <span class="card-title">Активные заявки</span>
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>№</th>
              <th>Питомец</th>
              <th>Врач</th>
              <th>Дата и время</th>
              <th>Статус</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="app in appointments" :key="app.id">
              <td class="mono text-muted">#{{ app.id }}</td>
              <td class="td-main">{{ app.pet }}</td>
              <td>{{ app.doctor }} · {{ app.spec }}</td>
              <td class="mono">{{ app.date }} · {{ app.time }}</td>
              <td>
                <span v-if="app.status === 'waiting'" class="badge badge-waiting">Ожидание</span>
                <span v-else class="badge badge-confirmed">Подтверждено</span>
              </td>
              <td>
                <button class="btn btn-ghost btn-sm" @click="cancel(app.id)">Отменить</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
