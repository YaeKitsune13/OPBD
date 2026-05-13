<script setup>
import { ref, onMounted } from 'vue'
import BaseModal from '../../ui/BaseModal.vue'
import { useToast } from '../../../utils/useToast'

const { showToast } = useToast()

const users = ref([])
const loading = ref(false)
const isDoctorModalOpen = ref(false)

// Данные для модалки назначения врача
const selectedUser = ref(null)
const speciality = ref('')

// 1. Загрузка всех пользователей
async function loadUsers() {
  const token = localStorage.getItem('token')
  try {
    const response = await fetch('/api/admin/users', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (response.ok) {
      users.value = await response.json()
    }
  } catch (e) {
    showToast("Ошибка при загрузке пользователей", "error")
  }
}

// 2. Открытие модалки
function openPromoteModal(user) {
  selectedUser.value = user
  speciality.value = ''
  isDoctorModalOpen.value = true
}

// 3. Процесс превращения в врача
async function promoteToDoctor() {
  if (!speciality.value) return
  
  loading.value = true
  const token = localStorage.getItem('token')
  const userId = selectedUser.value.user_id

  try {
    // Шаг А: Меняем роль пользователя на 'doctor'
    const roleRes = await fetch(`/api/admin/users/${userId}/role`, {
      method: 'PUT',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}` 
      },
      body: JSON.stringify({ role: 'doctor' })
    })

    // Шаг Б: Создаем запись в таблице врачей
    const docRes = await fetch('/api/admin/doctors', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}` 
      },
      body: JSON.stringify({
        user_id: userId,
        speciality: speciality.value
      })
    })

    if (roleRes.ok && docRes.ok) {
      showToast("Пользователь теперь врач!", "success")
      isDoctorModalOpen.value = false
      await loadUsers() // Обновляем список
    } else {
      showToast("Ошибка при сохранении данных", "error")
    }
  } catch (e) {
    showToast("Ошибка соединения с сервером", "error")
  } finally {
    loading.value = false
  }
}

onMounted(loadUsers)
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Управление пользователями</div>
        <div class="page-sub">Список всех аккаунтов и управление ролями</div>
      </div>
    </div>

    <div class="card">
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>ФИО</th>
              <th>Email / Телефон</th>
              <th>Роль</th>
              <th style="text-align: right">Действие</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.user_id">
              <td class="mono" style="color: var(--text3)">#{{ user.user_id }}</td>
              <td class="td-main">
                {{ user.last_name }} {{ user.first_name }} {{ user.middle_name }}
              </td>
              <td>
                <div style="font-size: 13px">{{ user.email }}</div>
                <div style="font-size: 11px; color: var(--text3)">{{ user.phone }}</div>
              </td>
              <td>
                <span class="badge" :class="{
                  'badge-confirmed': user.role === 'doctor',
                  'badge-waiting': user.role === 'client',
                  'badge-priority': user.role === 'admin'
                }">
                  {{ user.role }}
                </span>
              </td>
              <td style="text-align: right">
                <!-- Кнопка видна только для обычных клиентов -->
                <button 
                  v-if="user.role === 'client'"
                  class="btn btn-ghost btn-sm"
                  @click="openPromoteModal(user)"
                >
                  ✚ Сделать врачом
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Модалка назначения специализации -->
    <BaseModal :show="isDoctorModalOpen" title="Назначить врачом" @close="isDoctorModalOpen = false">
      <div v-if="selectedUser" style="margin-bottom: 20px">
        <p style="font-size: 14px; color: var(--text2)">
          Вы переводите пользователя <b>{{ selectedUser.last_name }} {{ selectedUser.first_name }}</b> в статус врача.
        </p>
      </div>

      <div class="form-group">
        <label>Выберите специализацию *</label>
        <select v-model="speciality">
          <option value="">-- Выберите из списка --</option>
          <option>Терапевт</option>
          <option>Хирург</option>
          <option>Офтальмолог</option>
          <option>Кардиолог</option>
          <option>Дерматолог</option>
          <option>Стоматолог</option>
        </select>
      </div>

      <template #footer>
        <button class="btn btn-ghost" @click="isDoctorModalOpen = false">Отмена</button>
        <button 
          class="btn btn-primary" 
          :disabled="!speciality || loading" 
          @click="promoteToDoctor"
        >
          {{ loading ? 'Сохранение...' : 'Подтвердить назначение' }}
        </button>
      </template>
    </BaseModal>
  </div>
</template>

<style scoped>
.badge-priority {
  background: #fef3c7;
  color: #92400e;
}
</style>