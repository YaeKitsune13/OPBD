<script setup>
import { ref, reactive, onMounted } from 'vue'
import BaseModal from '../../ui/BaseModal.vue'
import { useAuth } from '../../../utils/useAuth';
import { useToast } from '../../../utils/useToast';

const { logout } = useAuth();
const { showToast } = useToast();

const pets = ref([]) // Изначально пустой массив

// Состояние модалок
const isAddModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isEditModalOpen = ref(false)
const petToDelete = ref(null)

// Форма нового питомца
const newPet = reactive({
  name: '',
  type: 'Кошка',
  breed: '',
  dob: '',
  weight: '',
})

const editingPet = reactive({
  petId: 0,
  name: '',
  type: '',
  breed: '',
  dob: '',
  weight: '',
})

function openEditModal(pet) {
  editingPet.petId = pet.petId;
  editingPet.name = pet.name;
  editingPet.type = pet.species; // Обратите внимание: в DTO это species
  editingPet.breed = pet.breed;
  editingPet.dob = pet.dob;      // Дата уже в формате YYYY-MM-DD от сервера
  editingPet.weight = pet.weight;

  isEditModalOpen.value = true;
}

async function loadPets() {
  const token = localStorage.getItem('token');
  const userData = localStorage.getItem('user');
  if (!userData) return;

  const user = JSON.parse(userData);
  const ownerId = user.id;

  try {
    const response = await fetch(`/api/pets/owner/${ownerId}`, {
      headers: {
        "Accept": "application/json",
        "Authorization": `Bearer ${token}`
      }
    });

    if (response.ok) {
      const data = await response.json();
      // ВАЖНО: Присваиваем данные в ref
      pets.value = data;
    }
  } catch (e) {
    showToast("Ошибка при загрузке питомцев", "error");
  }
}
async function savePet() {
  // 1. Валидация перед отправкой
  if (!newPet.name || !newPet.dob || !newPet.type) {
    showToast("Заполните обязательные поля: Имя, Вид и Дата рождения", "error");
    return;
  }

  // 2. Достаем данные авторизации
  const token = localStorage.getItem('token');
  const userData = localStorage.getItem('user');
  if (!userData || !token) {
    showToast("Ошибка авторизации", "error");
    logout();
    return;
  }
  const ownerId = JSON.parse(userData).id;

  // 3. Подготовка данных
  // Вместо .toISOString() делаем чистую дату YYYY-MM-DD
  const formattedDate = newPet.dob.split('.').reverse().join('-');

  const payload = {
    name: newPet.name,
    breed: newPet.breed,
    dob: formattedDate, // Будет "2026-04-14"
    weight: Number(newPet.weight), // Гарантируем, что это число
    species: newPet.type,
    petId: 0,
    avatar: ""
  };

  try {
    const response = await fetch(`/api/pets?ownerId=${ownerId}`, {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    });

    if (response.status === 401) {
      showToast("Сессия истекла", "error");
      logout();
      return;
    }

    const result = await response.json();

    if (result.error) {
      showToast(result.error, "error");
      return;
    }

    // 4. Успех!
    showToast("Питомец успешно добавлен!", "success");

    // Обновляем список питомцев (вызываем твою функцию загрузки)
    await loadPets();

    // Закрываем модалку и чистим форму
    isAddModalOpen.value = false;
    resetForm();

  } catch (e) {
    showToast("Ошибка соединения с сервером", "error");
  }
}

async function updatePet() {
  if (!editingPet.name || !editingPet.dob) {
    showToast("Имя и дата рождения обязательны", "error");
    return;
  }

  const token = localStorage.getItem('token');

  const payload = {
    petId: editingPet.petId,
    name: editingPet.name,
    species: editingPet.type,
    breed: editingPet.breed,
    dob: editingPet.dob,
    weight: Number(editingPet.weight),
    avatar: ""
  };

  try {
    const response = await fetch(`/api/pets/${payload.petId}`, {
      method: 'PUT',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    });

    if (response.ok) {
      showToast("Данные обновлены", "success");
      isEditModalOpen.value = false;
      await loadPets(); // Обновляем список
    } else {
      const result = await response.json();
      showToast(result.error || "Ошибка при обновлении", "error");
    }
  } catch (e) {
    showToast("Ошибка соединения с сервером", "error");
  }
}

// Вспомогательная функция для очистки
function resetForm() {
  newPet.name = '';
  newPet.breed = '';
  newPet.dob = '';
  newPet.weight = '';
  newPet.type = '';
}

function confirmDelete(pet) {
  petToDelete.value = pet
  isDeleteModalOpen.value = true
}

async function deletePet() {
  if (!petToDelete.value) return;

  const token = localStorage.getItem('token');
  const petId = petToDelete.value.petId;

  try {
    const response = await fetch(`/api/pets/${petId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept': 'application/json'
      }
    });

    if (response.ok) {
      showToast("Питомец удален", "success");

      await loadPets();

    } else {
      const errorData = await response.json();
      showToast(errorData.error || "Ошибка при удалении", "error");
    }
  } catch (e) {
    showToast("Ошибка соединения с сервером", "error");
  } finally {
    isDeleteModalOpen.value = false;
    petToDelete.value = null;
  }
}

onMounted(() => {
  loadPets()
})
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Мои питомцы</div>
        <div class="page-sub">Карточки и медицинская информация</div>
      </div>
      <button class="btn btn-primary" @click="isAddModalOpen = true">✚ Добавить питомца</button>
    </div>

    <div class="pet-grid">
      <div v-for="pet in pets" :key="pet.id" class="pet-card">
        <div class="pet-avatar">{{ pet.avatar }}</div>
        <div class="pet-name">{{ pet.name }}</div>
        <div class="pet-meta">{{ pet.type }} · {{ pet.breed }}</div>
        <div class="pet-weight">⚖️ {{ pet.weight }} кг</div>

        <div style="display: flex; gap: 6px; margin-top: 10px">
          <button class="btn btn-ghost btn-sm" style="flex: 1" @click="openEditModal(pet)">
            ✏ Изменить
          </button>
          <button class="btn btn-danger btn-sm" style="flex: 1" @click="confirmDelete(pet)">
            🗑 Удалить
          </button>
        </div>
      </div>

      <div class="pet-card add-placeholder" @click="isAddModalOpen = true">
        <div class="add-icon">+</div>
        <div style="font-size: 12px">Добавить питомца</div>
      </div>
    </div>

    <!-- Модалка добавления -->
    <BaseModal :show="isAddModalOpen" title="Добавить питомца" @close="isAddModalOpen = false">
      <div class="form-grid">
        <div class="form-group">
          <label>Кличка *</label>
          <input v-model="newPet.name" type="text" placeholder="Только буквы" />
        </div>
        <div class="form-group">
          <label>Вид животного</label>
          <select v-model="newPet.type">
            <option>Кошка</option>
            <option>Собака</option>
            <option>Кролик</option>
            <option>Другое</option>
          </select>
        </div>
        <div class="form-group">
          <label>Порода</label>
          <input v-model="newPet.breed" type="text" placeholder="Напр. Британец" />
        </div>
        <div class="form-group">
          <label>Дата рождения *</label>
          <input v-model="newPet.dob" type="date" />
        </div>
        <div class="form-group">
          <label>Текущий вес (кг)</label>
          <input v-model="newPet.weight" type="number" step="0.1" placeholder="0.0" />
        </div>
      </div>
      <template #footer>
        <button class="btn btn-ghost" @click="isAddModalOpen = false">Отмена</button>
        <button class="btn btn-primary" @click="savePet">Сохранить</button>
      </template>
    </BaseModal>

    <!-- Модалка удаления -->
    <BaseModal :show="isDeleteModalOpen" title="Удалить питомца?" maxWidth="400px" @close="isDeleteModalOpen = false">
      <p style="color: var(--text2); font-size: 14px">
        Вы уверены, что хотите удалить карточку питомца <b>{{ petToDelete?.name }}</b>? Это действие необратимо.
      </p>
      <template #footer>
        <button class="btn btn-ghost" @click="isDeleteModalOpen = false">Отмена</button>
        <button class="btn btn-danger" @click="deletePet">Удалить</button>
      </template>
    </BaseModal>

    <!-- Модалка редактирования -->
    <BaseModal :show="isEditModalOpen" title="Редактировать питомца" @close="isEditModalOpen = false">
      <div class="form-grid">
        <div class="form-group">
          <label>Кличка *</label>
          <input v-model="editingPet.name" type="text" />
        </div>
        <div class="form-group">
          <label>Вид животного</label>
          <select v-model="editingPet.type">
            <option>Кошка</option>
            <option>Собака</option>
            <option>Кролик</option>
            <option>Другое</option>
          </select>
        </div>
        <div class="form-group">
          <label>Порода</label>
          <input v-model="editingPet.breed" type="text" />
        </div>
        <div class="form-group">
          <label>Дата рождения *</label>
          <input v-model="editingPet.dob" type="date" />
        </div>
        <div class="form-group">
          <label>Текущий вес (кг)</label>
          <input v-model="editingPet.weight" type="number" step="0.1" />
        </div>
      </div>
      <template #footer>
        <button class="btn btn-ghost" @click="isEditModalOpen = false">Отмена</button>
        <button class="btn btn-primary" @click="updatePet">Сохранить изменения</button>
      </template>
    </BaseModal>
  </div>
</template>

<style scoped>
.add-placeholder {
  border-style: dashed !important;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 180px;
  gap: 8px;
  color: var(--text3);
  cursor: pointer;
}

.add-placeholder:hover {
  color: var(--accent);
  border-color: var(--accent) !important;
}

.add-icon {
  font-size: 24px;
}
</style>
