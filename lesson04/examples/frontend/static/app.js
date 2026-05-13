const API = '/api/tasks';

async function loadTasks() {
    const res = await fetch(API);
    const tasks = await res.json();
    const list = document.getElementById('task-list');
    list.innerHTML = '';
    tasks.forEach(task => {
        const li = document.createElement('li');
        li.className = 'task-item' + (task.done ? ' done' : '');
        li.innerHTML = `
            <input type="checkbox" ${task.done ? 'checked' : ''} data-id="${task.id}">
            <span class="task-title">${escapeHtml(task.title)}</span>
            <button class="delete-btn" data-id="${task.id}">Delete</button>
        `;
        list.appendChild(li);
    });
}

document.getElementById('task-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const input = document.getElementById('task-input');
    const title = input.value.trim();
    if (!title) return;

    await fetch(API, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title }),
    });
    input.value = '';
    loadTasks();
});

document.getElementById('task-list').addEventListener('click', async (e) => {
    if (e.target.classList.contains('delete-btn')) {
        const id = e.target.dataset.id;
        await fetch(`${API}/${id}`, { method: 'DELETE' });
        loadTasks();
    }
});

document.getElementById('task-list').addEventListener('change', async (e) => {
    if (e.target.type === 'checkbox') {
        const id = e.target.dataset.id;
        const done = e.target.checked;
        await fetch(`${API}/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ done }),
        });
        loadTasks();
    }
});

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

loadTasks();
