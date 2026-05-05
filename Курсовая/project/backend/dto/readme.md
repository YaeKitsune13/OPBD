## 2. Описание DTO (Data Transfer Objects) для вывода в UI

DTO используются для передачи агрегированных данных из базы в интерфейс пользователя.

### 2.1 DashboardDTO (Экран "Обзор" — `DashboardPage.vue`)
Используется для вывода сводной информации на главной странице клиента.
- `clientName` (String): Имя для приветствия ("Иван").
- `petsCount` (Integer): Общее число питомцев (например, 3).
- `pendingApps` (Integer): Количество заявок в статусе "Ожидание".
- `totalVisits` (Integer): Суммарное количество визитов за всё время.
- `nextAppointment` (Object): `{ date, time, petName, doctorName }` — данные ближайшего приёма.
- `recentAppointments` (Array): Список последних 3–5 записей для виджета.
- `petsShort` (Array): Краткий список питомцев `{ name, species, breed, weight }` для виджета.

### 2.2 PetCardDTO (Экран "Мои питомцы" — `PetsPage.vue`)
Данные одной карточки питомца.
- `petId` (Integer): Идентификатор.
- `name` (String): Кличка.
- `species` (String): Вид.
- `breed` (String): Порода.
- `dob` (String): Дата рождения.
- `weight` (Decimal): Текущий вес в кг.
- `avatar` (String): Эмодзи-аватар, выбирается по виду (🐱, 🐶, 🐇).

### 2.3 AppointmentRowDTO (Экран "Записи" — `AppointmentsPage.vue`)
Данные одной строки таблицы заявок.
- `appointmentId` (Integer): Номер #ID.
- `petLabel` (String): Эмодзи + кличка.
- `doctorName` (String): ФИО врача.
- `specialty` (String): Специализация.
- `scheduledDate` (String): Дата.
- `scheduledTime` (String): Время.
- `status` (Enum): 'waiting' | 'confirmed' | 'rejected'.

### 2.4 HealthJournalDTO (Экран "Журнал здоровья" — `HistoryPage.vue`)
Данные одной карточки визита.
- `visitId` (Integer): Идентификатор.
- `date` (String): Дата визита.
- `time` (String): Время.
- `doctor` (String): ФИО врача.
- `diagnosis` (String): Диагноз с иконкой.
- `details` (String): Список назначений.
- `analysis` (String | null): Результаты анализов.
- `recommendations` (String | null): Рекомендации.
- `price` (String): Итоговая стоимость.

### 2.5 WeightPointDTO (Экран "Динамика веса" — `WeightPage.vue`)
Одна точка на графике веса.
- `label` (String): Подпись по оси X (месяц, например "апр").
- `value` (Decimal): Значение веса в кг.
- `date` (String): Полная дата измерения.
- `doctorName` (String): Врач, взвесивший питомца.

### 2.6 TodayScheduleDTO (Экран "Расписание" — `TodayPage.vue`)
Одна запись в расписании врача на день.
- `appointmentId` (Integer): Идентификатор.
- `time` (String): Время приёма.
- `petLabel` (String): Эмодзи + кличка.
- `ownerName` (String): ФИО владельца.
- `breed` (String): Порода.
- `reason` (String): Причина визита.
- `status` (Enum): 'waiting' | 'confirmed'.

### 2.7 ConductVisitDTO (Экран "Ведение приёма" — `ConductPage.vue`)
Данные текущего приёма.
- `selectedPet` (Object): `{ id, name, avatar, breed, owner }` — выбранный питомец.
- `anamnesis` (String): Текст жалоб.
- `diagnosis` (String): Введённый диагноз.
- `assignments` (Array): `[{ id, name, type, price, qty }]` — список назначений.
- `totalCost` (Integer): Вычисляемая сумма.

### 2.8 AnalyticsSummaryDTO (Экран "Аналитика" — `AnalyticsPage.vue`)
Сводные данные по KPI клиники.
- `monthlyVisits` (Integer): Визитов за месяц.
- `totalRevenue` (String): Суммарная выручка.
- `revenueChange` (String): Изменение к прошлому месяцу ("+12%").
- `activeClients` (Integer): Активных клиентов.
- `avgCheck` (Integer): Средний чек в рублях.
- `popularServices` (Array): `[{ name, count, revenue }]`.
- `doctorLoad` (Array): `[{ name, visitCount, loadStatus }]`.

### 2.9 RevenueReportDTO (Экран "Выручка" — `RevenuePage.vue`)
Детализация выручки за период.
- `periodTotal` (Integer): Итого за период.
- `servicesTotal` (Integer): За услуги.
- `dailyRows` (Array): `[{ date, visits, services, meds, total }]` — строки таблицы.

### 2.10 ServiceItemDTO (Экран "Услуги" — `ServicesPage.vue`)
Позиция справочника услуг.
- `id` (String): Артикул (001 …).
- `name` (String): Наименование.
- `desc` (String): Описание.
- `price` (Integer): Стоимость в рублях.

### 2.11 MedicationItemDTO (Экран "Медикаменты" — `MedsPage.vue`)
Позиция аптечного склада.
- `id` (String): Артикул (M001 …).
- `name` (String): Название.
- `desc` (String): Описание / показания.
- `price` (Integer): Цена за единицу.
- `expiry` (String): Срок годности (YYYY-MM).
- `status` (Enum): 'ok' | 'expired'.
