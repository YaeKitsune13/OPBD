<script setup>
import { ref, onMounted, reactive } from 'vue'
import BaseModal from '../../ui/BaseModal.vue'
import { useToast } from '../../../utils/useToast'

const { showToast } = useToast()

const users = ref([])
const loading = ref(false)

// Состояние для модалок
const isAddDoctorModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const userToDelete = ref(null)

// Форма нового врача
const doctorForm = reactive({
  first_name: '',
  last_name: '',
  middle_name: '',
  email: '',
  phone: '',
  password: '',
  speciality: ''
})

// 1. Загрузка пользователей
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

// 2. Добавление нового врача
async function createDoctor() {
  loading.value = true
  const token = localStorage.getItem('token')
  
  try {
    const response = await fetch('/api/admin/doctors/create-full', { // Предполагаемый эндпоинт для создания "под ключ"
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}` 
      },
      body: JSON.stringify(doctorForm)
    })

    if (response.ok) {
      showToast("Врач успешно добавлен", "success")
      isAddDoctorModalOpen.value = false
      // Сброс формы
      Object.assign(doctorForm, { first_name: '', last_name: '', middle_name: '', email: '', phone: '', password: '', speciality: '' })
      await loadUsers()
    } else {
      showToast("Ошибка при создании врача", "error")
    }
  } catch (e) {
    showToast("Ошибка соединения", "error")
  } finally {
    loading.value = false
  }
}

// 3. Удаление пользователя
function confirmDelete(user) {
  userToDelete.value = user
  isDeleteModalOpen.value = true
}

async function deleteUser() {
  if (!userToDelete.value) return
  
  loading.value = true
  const token = localStorage.getItem('token')
  try {
    const response = await fetch(`/api/admin/users/${userToDelete.value.user_id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      showToast("Пользователь удален", "success")
      isDeleteModalOpen.value = false
      await loadUsers()
    } else {
      showToast("Не удалось удалить пользователя", "error")
    }
  } catch (e) {
    showToast("Ошибка при удалении", "error")
  } finally {
    loading.value = false
    userToDelete.value = null
  }
}

onMounted(loadUsers)
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Управление пользователями</div>
        <div class="page-sub">Список всех аккаунтов и управление доступом</div>
      </div>
      <!-- КНОПКА ДОБАВИТЬ ДОКТОРА -->
      <button class="btn btn-primary" @click="isAddDoctorModalOpen = true">
        + Добавить врача
      </button>
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
                <!-- ЗАМЕНЕНА КНОПКА НА УДАЛЕНИЕ -->
                <button 
                  class="btn btn-ghost btn-sm btn-danger"
                  @click="confirmDelete(user)"
                >
                  🗑 Удалить
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Модалка добавления врача -->
    <BaseModal :show="isAddDoctorModalOpen" title="Новый врач" @close="isAddDoctorModalOpen = false">
      <div class="form-grid">
        <div class="form-group">
          <label>Фамилия</label>
          <input v-model="doctorForm.last_name" type="text" placeholder="Иванов">
        </div>
        <div class="form-group">
          <label>Имя</label>
          <input v-model="doctorForm.first_name" type="text" placeholder="Иван">
        </div>
        <div class="form-group">
          <label>Email</label>
          <input v-model="doctorForm.email" type="email" placeholder="example@mail.com">
        </div>
        <div class="form-group">
          <label>Пароль</label>
          <input v-model="doctorForm.password" type="password" placeholder="******">
        </div>
        <div class="form-group">
          <label>Специализация</label>
          <select v-model="doctorForm.speciality">
            <option value="">-- Выберите --</option>
            <option>Терапевт</option>
            <option>Хирург</option>
            <option>Кардиолог</option>
            <option>Стоматолог</option>
          </select>
        </div>
      </div>

      <template #footer>
        <button class="btn btn-ghost" @click="isAddDoctorModalOpen = false">Отмена</button>
        <button class="btn btn-primary" :disabled="loading" @click="createDoctor">
          {{ loading ? 'Создание...' : 'Создать врача' }}
        </button>
      </template>
    </BaseModal>

    <!-- Модалка подтверждения удаления -->
    <BaseModal :show="isDeleteModalOpen" title="Подтверждение удаления" @close="isDeleteModalOpen = false">
      <p v-if="userToDelete">
        Вы уверены, что хотите удалить пользователя <b>{{ userToDelete.last_name }} {{ userToDelete.first_name }}</b>?
        Это действие необратимо.
      </p>
      <template #footer>
        <button class="btn btn-ghost" @click="isDeleteModalOpen = false">Отмена</button>
        <button class="btn btn-danger" :disabled="loading" @click="deleteUser">
          {{ loading ? 'Удаление...' : 'Да, удалить' }}
        </button>
      </template>
    </BaseModal>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.badge-priority {
  background: #fef3c7;
  color: #92400e;
}
.btn-danger {
  color: #dc2626;
}
.btn-danger:hover {
  background: #fee2e2;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}
/* Чтобы некоторые поля были на всю ширину */
.form-group:nth-child(5) {
  grid-column: span 2;
}
</style>