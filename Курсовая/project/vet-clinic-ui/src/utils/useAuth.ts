import { ref } from 'vue'

// Создаем глобальные переменные вне функции, чтобы они были общими для всего приложения
const isAuthenticated = ref(!!localStorage.getItem('token'))

export function useAuth() {
    
    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        isAuthenticated.value = false;
        // Можно даже сделать релоад страницы, чтобы сбросить всё состояние
        window.location.reload(); 
    }

    return {
        isAuthenticated,
        logout
    }
}