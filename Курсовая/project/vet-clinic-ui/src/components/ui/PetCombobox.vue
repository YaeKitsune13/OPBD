<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: Object,
  pets: Array, // [{ id, name, avatar, breed, owner }]
  placeholder: { type: String, default: 'Выберите питомца' },
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const containerRef = ref(null)

function selectPet(pet) {
  emit('update:modelValue', pet)
  isOpen.value = false
}

function toggleDropdown() {
  isOpen.value = !isOpen.value
}

const clickOutside = (e) => {
  if (containerRef.value && !containerRef.value.contains(e.target)) {
    isOpen.value = false
  }
}

onMounted(() => document.addEventListener('click', clickOutside))
onUnmounted(() => document.removeEventListener('click', clickOutside))
</script>

<template>
  <div class="pet-select-custom" ref="containerRef">
    <!-- Основное поле (выглядит как инпут, но это div) -->
    <div class="select-trigger" :class="{ 'is-open': isOpen }" @click="toggleDropdown">
      <div class="selected-content">
        <template v-if="modelValue">
          <span class="selected-avatar">{{ modelValue.avatar }}</span>
          <span class="selected-name">{{ modelValue.name }}</span>
        </template>
        <span v-else class="placeholder">{{ placeholder }}</span>
      </div>
      <div class="arrow" :class="{ rotated: isOpen }">▼</div>
    </div>

    <!-- Выпадающий список -->
    <Transition name="slide-fade">
      <div v-if="isOpen" class="dropdown-list">
        <div
          v-for="pet in pets"
          :key="pet.id"
          class="pet-option"
          :class="{ 'is-active': modelValue?.id === pet.id }"
          @click="selectPet(pet)"
        >
          <div class="option-avatar">{{ pet.avatar }}</div>
          <div class="option-info">
            <div class="option-name">{{ pet.name }}</div>
            <div class="option-sub">{{ pet.breed }} • {{ pet.owner }}</div>
          </div>
          <div v-if="modelValue?.id === pet.id" class="check-mark">✓</div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.pet-select-custom {
  position: relative;
  width: 100%;
  max-width: 320px; /* Сузил для компактности в заголовке */
  user-select: none; /* Запрещаем выделение текста */
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--surface2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 8px 12px;
  min-height: 40px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.select-trigger:hover {
  border-color: var(--border2);
  background: var(--surface3);
}

.select-trigger.is-open {
  border-color: var(--accent);
  background: var(--surface2);
}

.selected-content {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
}

.selected-avatar {
  font-size: 18px;
}

.selected-name {
  font-weight: 600;
  color: var(--text);
}

.placeholder {
  color: var(--text3);
}

.arrow {
  font-size: 9px;
  color: var(--text3);
  transition: transform 0.3s ease;
  margin-left: 10px;
}

.arrow.rotated {
  transform: rotate(180deg);
  color: var(--accent);
}

.dropdown-list {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--border2);
  border-radius: var(--radius);
  max-height: 300px;
  overflow-y: auto;
  z-index: 9999;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
}

.pet-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  cursor: pointer;
  transition: background 0.15s;
}

.pet-option:hover {
  background: var(--surface3);
}

.pet-option.is-active {
  background: var(--accent-dim);
}

.option-avatar {
  width: 32px;
  height: 32px;
  background: var(--surface2);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.option-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

.option-sub {
  font-size: 11px;
  color: var(--text3);
}

.check-mark {
  margin-left: auto;
  color: var(--accent);
  font-weight: bold;
}

/* Анимация */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.2s ease-out;
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Стилизация скроллбара */
.dropdown-list::-webkit-scrollbar {
  width: 4px;
}
.dropdown-list::-webkit-scrollbar-thumb {
  background: var(--border2);
  border-radius: 2px;
}
</style>
