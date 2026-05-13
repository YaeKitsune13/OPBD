<script setup lang="ts">
import { ref } from 'vue'
import { useToast } from '../../../utils/useToast';
const { showToast } = useToast();

const appointments = ref([])

async function loadAppoinments() {
  const token = localStorage.getItem('token');
  const userRaw = localStorage.getItem('user');
  if (!userRaw) return;

  const userData = JSON.parse(userRaw);
  const userId = userData.user_id || userData.id;
  try {
    const appointmentsRes = await fetch(`/api/appointments/owner/${userId}`, {
      headers: { "Authorization": `Bearer ${token}` }
    });
    if (appointmentsRes.ok) appointments.value = await appointmentsRes.json();
    console.log(appointments.value);
  } catch (e) {
    showToast("Ошибка при загрузке данных", "error")
  }
}

loadAppoinments()

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
            <tr v-for="(app, index) in appointments" :key="app.id">
              <td class="mono text-muted"># {{ index+1 }}</td>
              <td class="td-main">{{ app.petLabel }}</td>
              <td>{{ app.doctorName }} · {{ app.specialty }}</td>
              <td class="mono">{{ app.scheduledDate }} · {{ app.scheduledTime }}</td>
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
