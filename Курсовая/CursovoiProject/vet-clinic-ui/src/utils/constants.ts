export const PET_SPECIES = [
  { value: "dog", label: "Собака", emoji: "🐶" },
  { value: "cat", label: "Кошка", emoji: "🐱" },
  { value: "rabbit", label: "Кролик", emoji: "🐰" },
  { value: "bird", label: "Птица", emoji: "🐦" },
  { value: "other", label: "Другое", emoji: "🐾" },
];

export const SERVICES = [
  { id: 1, name: "Первичный осмотр", price: 800, duration: 30 },
  { id: 2, name: "Вакцинация", price: 1200, duration: 20 },
  { id: 3, name: "Хирургическая операция", price: 8000, duration: 120 },
  { id: 4, name: "УЗИ", price: 2000, duration: 40 },
  { id: 5, name: "Анализ крови", price: 1500, duration: 15 },
  { id: 6, name: "Стоматология", price: 3000, duration: 60 },
  { id: 7, name: "Дерматология", price: 1000, duration: 30 },
  { id: 8, name: "Кардиология", price: 2500, duration: 45 },
];

export const MEDICATIONS = [
  {
    id: 1,
    name: 'Антибиотик "Амоксициллин"',
    description: "Широкого спектра действия, 10 таблеток",
    price: 450,
    inStock: true,
    category: "Антибиотики",
  },
  {
    id: 2,
    name: "Витамины для собак",
    description: "Комплекс витаминов A,B,C,D, 30 таблеток",
    price: 680,
    inStock: true,
    category: "Витамины",
  },
  {
    id: 3,
    name: 'Капли от блох "Фронтлайн"',
    description: "Для собак весом до 10кг",
    price: 890,
    inStock: true,
    category: "Паразиты",
  },
  {
    id: 4,
    name: 'Обезболивающее "Мелоксикам"',
    description: "1мг/мл, раствор для инъекций 10мл",
    price: 320,
    inStock: false,
    category: "Обезболивающие",
  },
  {
    id: 5,
    name: "Витамины для кошек",
    description: "Омега-3 и таурин, 60 капсул",
    price: 590,
    inStock: true,
    category: "Витамины",
  },
  {
    id: 6,
    name: 'Антигельминтное "Дронтал"',
    description: "Для кошек, 8 таблеток",
    price: 740,
    inStock: true,
    category: "Паразиты",
  },
];

export const APPOINTMENT_STATUSES = {
  waiting: { label: "Ожидает", class: "badge-waiting" },
  confirmed: { label: "Принята", class: "badge-confirmed" },
  rejected: { label: "Отклонена", class: "badge-rejected" },
  done: { label: "Завершена", class: "badge-info" },
};
export const ROLES = {
  client: "client",
  doctor: "doctor",
};

export const API = {
  auth: {
    login: "/api/auth/login",
    register: "/api/auth/register",
  },
  pets: "/api/pets",
  appointments: "/api/appointments",
  visits: "/api/visits",
  medications: "/api/medications/orders",
};
