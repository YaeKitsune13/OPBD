<script setup>
// Принимаем данные для печати: инфо о питомце и массив визитов
defineProps(['pet', 'visits', 'type'])

const formatDate = () => new Date().toLocaleString('ru-RU')
</script>

<template>
  <div id="printRoot">
    <!-- Шапка клиники -->
    <div class="pr-header">
      <div class="pr-clinic-name">🐾 ВетКлиника</div>
      <div class="pr-clinic-sub">г. Москва, ул. Примерная, д. 1 · Тел: +7 (495) 000-00-00</div>
    </div>

    <div class="pr-doc-title">
      {{ type === 'weight' ? 'Отчет о динамике веса' : 'Журнал здоровья / Выписка' }}
    </div>
    <div class="pr-doc-num">Дата выдачи: {{ formatDate() }}</div>

    <!-- Инфо о питомце -->
    <div class="pr-section-title">Данные о животном</div>
    <div class="pr-pet-block" v-if="pet">
      <div class="pr-row">
        <span>Кличка:</span> <b>{{ pet.name }}</b>
      </div>
      <div class="pr-row"><span>Вид:</span> {{ pet.type }}</div>
      <div class="pr-row"><span>Порода:</span> {{ pet.breed }}</div>
      <div class="pr-row"><span>Владелец:</span> Иванов И.И.</div>
    </div>

    <!-- Таблица визитов -->
    <div class="pr-section-title">Записи</div>
    <div v-for="v in visits" :key="v.id" class="pr-visit">
      <div class="pr-visit-head">
        <b>{{ v.date }} · {{ v.time }}</b>
        <span>Врач: {{ v.doctor }}</span>
      </div>
      <div class="pr-visit-body">
        <div><b>Диагноз:</b> {{ v.diagnosis }}</div>
        <div style="font-size: 10pt; color: #555">{{ v.details }}</div>
        <div v-if="v.recommendations" style="margin-top: 5px">
          <b>Назначения:</b> {{ v.recommendations }}
        </div>
      </div>
    </div>

    <!-- Подписи -->
    <div class="pr-sign-block">
      <div>
        <div class="pr-sign-line"></div>
        <div class="pr-sign-label">Врач (подпись)</div>
      </div>
      <div>
        <div class="pr-sign-line"></div>
        <div class="pr-sign-label">Владелец (подпись)</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Эти стили активируются только при печати благодаря main.css @media print */
#printRoot {
  display: none;
  padding: 20px;
  background: white;
  color: black;
}

@media print {
  #printRoot {
    display: block !important;
  }
  .pr-header {
    text-align: center;
    border-bottom: 2px solid #000;
    margin-bottom: 20px;
    padding-bottom: 10px;
  }
  .pr-clinic-name {
    font-size: 20pt;
    font-weight: bold;
  }
  .pr-doc-title {
    text-align: center;
    text-transform: uppercase;
    margin: 10px 0;
    border: 1px solid #ccc;
    padding: 5px;
  }
  .pr-pet-block {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
    background: #f9f9f9;
    padding: 10px;
    border: 1px solid #ddd;
  }
  .pr-visit {
    border: 1px solid #eee;
    margin-top: 10px;
    page-break-inside: avoid;
  }
  .pr-visit-head {
    background: #f0f0f0;
    padding: 5px 10px;
    display: flex;
    justify-content: space-between;
  }
  .pr-visit-body {
    padding: 10px;
  }
  .pr-sign-block {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 50px;
    margin-top: 50px;
  }
  .pr-sign-line {
    border-bottom: 1px solid #000;
    height: 30px;
  }
  .pr-sign-label {
    font-size: 9pt;
    color: #666;
  }
}
</style>
