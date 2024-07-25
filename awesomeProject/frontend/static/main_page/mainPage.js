document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('login-form');
    const postsContainer = document.getElementById('posts-container');
    const postForm = document.getElementById('post-form');
    const deleteAllBtn = document.getElementById('delete-all-btn');

    const isAdmin = localStorage.getItem('isAdmin') === 'true';

    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;

            fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            })
                .then(response => response.json())
                .then(response => {
                    if (response.success) {
                        localStorage.setItem('isAdmin', response.isAdmin);
                        window.location.href = '/static/main_page/mainpage.html';
                    } else {
                        alert('Login failed. Please try again.');
                    }
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('An error occurred. Please try again.');
                });
        });
    }

    if (postsContainer) {
        loadItems();
    }
});

function loadItems() {
    fetch('http://localhost:8080/items')
        .then(response => response.json())
        .then(items => {
            const itemList = document.getElementById('posts-container');
            if (itemList) {
                const isAdmin = localStorage.getItem('isAdmin') === 'true';
                itemList.innerHTML = items.map(item => `
                <div class="blog-post">
                    <h3>${item.blogname}</h3>
                    <p><strong>Author:</strong> ${item.author}</p>
                    <p>${item.content}</p>
                    ${isAdmin ? `<button onclick="deleteItem('${item.id}')">Delete</button>` : ''}
                </div>
            `).join('');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('An error occurred. Please try again.');
        });
}