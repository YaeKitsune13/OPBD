async function fetchPets() {
    const token = localStorage.getItem('token'); // Достаем наш JWT

    const response = await fetch('/api/pets/owner/1', {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            // Вот эта строчка сообщает бэкенду, КТО делает запрос
            'Authorization': `Bearer ${token}` 
        }
    });

    return await response.json();
}