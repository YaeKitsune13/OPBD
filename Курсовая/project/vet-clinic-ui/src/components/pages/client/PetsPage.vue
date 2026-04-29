<script setup>
import { ref, reactive } from 'vue'
import BaseModal from '../../ui/BaseModal.vue'

const pets = ref([
  {
    id: 1,
    name: 'Барсик',
    type: 'Кот',
    breed: 'Британец',
    dob: '15.03.2021',
    weight: 4.2,
    avatar: '🐱',
  },
  {
    id: 2,
    name: 'Рыжик',
    type: 'Пёс',
    breed: 'Лабрадор',
    dob: '08.11.2023',
    weight: 28.5,
    avatar: '🐶',
  },
])

// Состояние модалок
const isAddModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const petToDelete = ref(null)

// Форма нового питомца
const newPet = reactive({
  name: '',
  type: 'Кошка',
  breed: '',
  dob: '',
  weight: '',
})

function savePet() {
  if (!newPet.name || !newPet.dob) return alert('Заполните обязательные поля')

  const id = Date.now()
  const avatar = newPet.type === 'Кошка' ? '🐱' : newPet.type === 'Собака' ? '🐶' : '🐇'

  pets.value.push({ id, ...newPet, avatar })
  isAddModalOpen.value = false

  // Сброс формы
  newPet.name = ''
  newPet.breed = ''
  newPet.dob = ''
  newPet.weight = ''
}

function confirmDelete(pet) {
  petToDelete.value = pet
  isDeleteModalOpen.value = true
}

function deletePet() {
  pets.value = pets.value.filter((p) => p.id !== petToDelete.value.id)
  isDeleteModalOpen.value = false
}
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
          <button class="btn btn-ghost btn-sm" style="flex: 1">✏ Изменить</button>
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
    <BaseModal
      :show="isDeleteModalOpen"
      title="Удалить питомца?"
      maxWidth="400px"
      @close="isDeleteModalOpen = false"
    >
      <p style="color: var(--text2); font-size: 14px">
        Вы уверены, что хотите удалить карточку питомца <b>{{ petToDelete?.name }}</b
        >? Это действие необратимо.
      </p>
      <template #footer>
        <button class="btn btn-ghost" @click="isDeleteModalOpen = false">Отмена</button>
        <button class="btn btn-danger" @click="deletePet">Удалить</button>
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
