document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const deleteAllBtn = document.getElementById('delete-all-btn');
    const viewByIdBtn = document.getElementById('view-by-id-btn');
    const viewAllBtn = document.getElementById('view-all-btn');
    const postIdInput = document.getElementById('post-id-input');
    const submitBtn = document.getElementById('submit-btn');
    const editBtn = document.getElementById('edit-btn');

    let currentEditId = null;

    function renderPosts(posts) {
        postsContainer.innerHTML = '';
        if (Array.isArray(posts) && posts.length > 0) {
            posts.forEach(post => {
                const postElement = document.createElement('div');
                postElement.classList.add('blog-post');
                postElement.innerHTML = `
                    <h3>${post.blogname}</h3>
                    <p class="author">By: ${post.author}</p>
                    <p class="content">${post.content}</p>
                    <div class="actions">
                        <button class="edit-btn" data-id="${post.id}">Edit</button>
                        <button class="delete-btn" data-id="${post.id}">Delete</button>
                    </div>
                `;
                postsContainer.appendChild(postElement);
            });
        } else {
            postsContainer.innerHTML = '<p>No posts to display</p>';
        }
    }

    function fetchPosts() {
        fetch('http://localhost:8080/items')
            .then(response => response.json())
            .then(posts => renderPosts(posts))
            .catch(error => {
                console.error('Error:', error);
                renderPosts([]);
            });
    }

    form.addEventListener('submit', function(e) {
        e.preventDefault();
        const post = {
            blogname: document.getElementById('blogname').value,
            author: document.getElementById('author').value,
            content: document.getElementById('content').value
        };

        const url = currentEditId ? `http://localhost:8080/items/${currentEditId}` : 'http://localhost:8080/items';
        const method = currentEditId ? 'PUT' : 'POST';

        fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(post)
        })
            .then(() => {
                currentEditId = null;
                submitBtn.textContent = 'Add Post';
                editBtn.style.display = 'none';
                fetchPosts();
                form.reset();
            })
            .catch(error => console.error('Error:', error));
    });

    postsContainer.addEventListener('click', function(e) {
        if (e.target.classList.contains('delete-btn')) {
            const id = e.target.getAttribute('data-id');
            fetch(`http://localhost:8080/items/${id}`, { method: 'DELETE' })
                .then(() => fetchPosts())
                .catch(error => console.error('Error:', error));
        } else if (e.target.classList.contains('edit-btn')) {
            const id = e.target.getAttribute('data-id');
            fetch(`http://localhost:8080/items/${id}`)
                .then(response => response.json())
                .then(post => {
                    document.getElementById('blogname').value = post.blogname;
                    document.getElementById('author').value = post.author;
                    document.getElementById('content').value = post.content;
                    currentEditId = post.id;
                    submitBtn.textContent = 'Update Post';
                    editBtn.style.display = 'inline-block';
                })
                .catch(error => console.error('Error:', error));
        }
    });

    deleteAllBtn.addEventListener('click', function() {
        fetch('http://localhost:8080/items', { method: 'DELETE' })
            .then(() => fetchPosts())
            .catch(error => console.error('Error:', error));
    });

    viewByIdBtn.addEventListener('click', function() {
        const id = postIdInput.value;
        fetch(`http://localhost:8080/items/${id}`)
            .then(response => response.json())
            .then(post => renderPosts([post]))
            .catch(() => {
                console.error('Post not found');
                renderPosts([]);
            });
    });

    viewAllBtn.addEventListener('click', fetchPosts);

    editBtn.addEventListener('click', function() {
        currentEditId = null;
        submitBtn.textContent = 'Add Post';
        editBtn.style.display = 'none';
        form.reset();
    });

    fetchPosts();
});