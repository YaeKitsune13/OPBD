<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from '../../../utils/useToast';
import { useAuth } from '../../../utils/useAuth';

const emit = defineEmits(['navigate'])
const { showToast } = useToast();
const { logout } = useAuth();

// --- СОСТОЯНИЕ ---
const pets = ref([])
const doctors = ref([]) // Теперь это массив объектов с сервера
const loading = ref(false)

const appointmentForm = reactive({
  petId: '',
  doctorId: '',
  specialization: '',
  date: '',
  comment: ''
})

// --- ЛОГИКА ЗАГРУЗКИ ДАННЫХ ---

async function fetchData() {
  const token = localStorage.getItem('token');
  const userData = JSON.parse(localStorage.getItem('user'));
  if (!userData) return;

  try {
    // 1. Загружаем питомцев
    const petsRes = await fetch(`/api/pets/owner/${userData.id}`, {
      headers: { "Authorization": `Bearer ${token}` }
    });
    if (petsRes.ok) pets.value = await petsRes.json();

    // 2. Загружаем врачей
    // Предполагаем, что есть эндпоинт /api/doctors
    const docsRes = await fetch(`/api/doctors`, {
      headers: { "Authorization": `Bearer ${token}` }
    });
    if (docsRes.ok) doctors.value = await docsRes.json();

  } catch (e) {
    showToast("Ошибка при загрузке данных", "error");
  }
}

// Фильтруем врачей по выбранной специализации
const filteredDoctors = computed(() => {
  if (!appointmentForm.specialization) return [];
  return doctors.value.filter(d => d.specialization === appointmentForm.specialization);
});

// Список уникальных специализаций для селекта
const specializations = computed(() => {
  const specs = doctors.value.map(d => d.specialization);
  return [...new Set(specs)]; // Удаляем дубликаты
});

// --- ОТПРАВКА ЗАЯВКИ ---

async function sendRequest() {
  // Валидация
  if (!appointmentForm.petId || !appointmentForm.doctorId || !appointmentForm.date) {
    showToast("Пожалуйста, заполните все обязательные поля", "warning");
    return;
  }

  const token = localStorage.getItem('token');
  loading.value = true;

  const payload = {
    petId: Number(appointmentForm.petId),
    doctorId: Number(appointmentForm.doctorId),
    appointmentDate: appointmentForm.date, // "2024-05-20"
    description: appointmentForm.comment,
    status: "scheduled"
  };

  try {
    const response = await fetch('/api/appointments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    });

    if (response.ok) {
      showToast("Запись успешно создана!", "success");
      emit('navigate', 'appointments'); // Переход к списку записей
    } else {
      const err = await response.json();
      showToast(err.error || "Ошибка при записи", "error");
    }
  } catch (e) {
    showToast("Ошибка соединения с сервером", "error");
  } finally {
    loading.value = false;
  }
}

onMounted(fetchData);
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Записаться на приём</div>
        <div class="page-sub">Выберите питомца, врача и удобное время</div>
      </div>
    </div>

    <div class="card">
      <div class="card-body">
        <div class="form-grid">
          <!-- ПИТОМЕЦ -->
          <div class="form-group">
            <label>Питомец *</label>
            <select v-model="appointmentForm.petId">
              <option value="">— Выберите питомца —</option>
              <option v-for="pet in pets" :key="pet.petId" :value="pet.petId">
                {{ pet.avatar }} {{ pet.name }}
              </option>
            </select>
          </div>

          <!-- СПЕЦИАЛИЗАЦИЯ -->
          <div class="form-group">
            <label>Специализация *</label>
            <select v-model="appointmentForm.specialization" @change="appointmentForm.doctorId = ''">
              <option value="">— Выберите специализацию —</option>
              <option v-for="spec in specializations" :key="spec" :value="spec">
                {{ spec }}
              </option>
            </select>
          </div>

          <!-- ВРАЧ -->
          <div class="form-group">
            <label>Врач *</label>
            <select v-model="appointmentForm.doctorId" :disabled="!appointmentForm.specialization">
              <option value="">— Выберите врача —</option>
              <option v-for="doc in filteredDoctors" :key="doc.id" :value="doc.id">
                {{ doc.name }}
              </option>
            </select>
          </div>

          <!-- ДАТА -->
          <div class="form-group">
            <label>Дата приёма *</label>
            <input v-model="appointmentForm.date" type="date" :min="new Date().toISOString().split('T')[0]" />
          </div>

          <!-- КОММЕНТАРИЙ -->
          <div class="form-group full">
            <label>Комментарий / Жалобы</label>
            <textarea v-model="appointmentForm.comment" placeholder="Опишите проблему..."></textarea>
          </div>
        </div>

        <div class="mt-12">
          <button class="btn btn-primary" :disabled="loading" @click="sendRequest">
            {{ loading ? 'Отправка...' : 'Подтвердить запись' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>