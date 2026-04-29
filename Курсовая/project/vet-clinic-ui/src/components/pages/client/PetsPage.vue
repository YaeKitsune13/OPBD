<script setup>
import { ref } from 'vue'

const emit = defineEmits(['open-modal'])

// В будущем эти данные придут с бэкенда
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
  {
    id: 3,
    name: 'Снежок',
    type: 'Кролик',
    breed: 'Вислоухий',
    dob: '20.06.2024',
    weight: 2.1,
    avatar: '🐇',
  },
])
</script>

<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title">Мои питомцы</div>
        <div class="page-sub">Карточки и медицинская информация</div>
      </div>
      <button class="btn btn-primary" @click="emit('open-modal', 'add-pet')">
        ✚ Добавить питомца
      </button>
    </div>

    <div class="pet-grid">
      <!-- Карточка существующего питомца -->
      <div v-for="pet in pets" :key="pet.id" class="pet-card">
        <div class="pet-avatar">{{ pet.avatar }}</div>
        <div class="pet-name">{{ pet.name }}</div>
        <div class="pet-meta">{{ pet.type }} · {{ pet.breed }}</div>
        <div class="pet-meta mt-4">Дата рождения: {{ pet.dob }}</div>
        <div class="pet-weight">⚖️ {{ pet.weight }} кг</div>

        <div style="display: flex; gap: 6px; margin-top: 10px">
          <button
            class="btn btn-ghost btn-sm"
            style="flex: 1"
            @click="emit('open-modal', 'edit-pet', pet)"
          >
            ✏ Изменить
          </button>
          <button
            class="btn btn-danger btn-sm"
            style="flex: 1"
            @click="emit('open-modal', 'delete-pet', pet)"
          >
            🗑 Удалить
          </button>
        </div>
      </div>

      <!-- Кнопка-заглушка "Добавить" -->
      <div class="pet-card add-placeholder" @click="emit('open-modal', 'add-pet')">
        <div class="add-icon">+</div>
        <div style="font-size: 12px">Добавить питомца</div>
      </div>
    </div>
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
  transition: all 0.2s;
}
.add-placeholder:hover {
  color: var(--accent);
  border-color: var(--accent) !important;
}
.add-icon {
  font-size: 24px;
}
</style>
