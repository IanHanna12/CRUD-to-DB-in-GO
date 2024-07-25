document.addEventListener('DOMContentLoaded', function() {
    const isAdmin = localStorage.getItem('isAdmin') === 'true';
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;

    const createBtn = document.getElementById('create-btn');
    const readBtn = document.getElementById('read-btn');
    const updateBtn = document.getElementById('update-btn');
    const deleteBtn = document.getElementById('delete-btn');

    if (isAdmin) {
        enableAdminFunctionality();
    } else {
        disableNonAdminFunctionality();
    }

    function disableNonAdminFunctionality() {
        [createBtn, updateBtn, deleteBtn].forEach(btn => {
            if (btn) btn.style.display = 'none';
        });
        if (form) form.style.display = 'none';
        if (readBtn) readBtn.addEventListener('click', fetchPosts);
    }

    function enableAdminFunctionality() {
        [createBtn, readBtn, updateBtn, deleteBtn].forEach(btn => {
            if (btn) btn.addEventListener('click', getButtonFunction(btn.id));
        });
        if (form) form.style.display = 'block';
    }

    function getButtonFunction(buttonId) {
        const functions = {
            'create-btn': createItem,
            'read-btn': fetchPosts,
            'update-btn': updateItem,
            'delete-btn': deleteItem
        };
        return functions[buttonId];
    }

    function createItem() {
        currentEditId = null;
        form.reset();
        form.style.display = 'block';
    }

    function updateItem() {
        const id = prompt('Enter the ID of the post to update:');
        if (id) {
            fetch(`${apiUrl}/${id}`)
                .then(res => res.json())
                .then(post => {
                    form.style.display = 'block';
                    form.blogname.value = post.blogname;
                    form.author.value = post.author;
                    form.content.value = post.content;
                    currentEditId = id;
                });
        }
    }

    function deleteItem() {
        const id = prompt('Enter the ID of the post to delete:');
        if (id) {
            fetch(`${apiUrl}/${id}`, { method: 'DELETE' }).then(fetchPosts);
        }
    }

    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(post => `
            <div class="blog-post">
                <h3>${post.blogname}</h3>
                <p>By: ${post.author}</p>
                <p>${post.content}</p>
                ${isAdmin ? `
                    <button onclick="editPost(${post.id})">Edit</button>
                    <button onclick="deletePost(${post.id})">Delete</button>
                ` : ''}
            </div>
        `).join('') || '<p>No posts available</p>';
    }

    function fetchPosts() {
        fetch(apiUrl).then(res => res.json()).then(renderPosts);
    }

    if (form) {
        form.onsubmit = function(e) {
            e.preventDefault();
            const post = {
                blogname: form.blogname.value,
                author: form.author.value,
                content: form.content.value
            };
            fetch(currentEditId ? `${apiUrl}/${currentEditId}` : apiUrl, {
                method: currentEditId ? 'PUT' : 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(post)
            }).then(() => {
                currentEditId = null;
                form.reset();
                fetchPosts();
            });
        };
    }

    window.editPost = function(id) {
        if (isAdmin) {
            fetch(`${apiUrl}/${id}`).then(res => res.json()).then(post => {
                form.style.display = 'block';
                form.blogname.value = post.blogname;
                form.author.value = post.author;
                form.content.value = post.content;
                currentEditId = id;
            });
        } else {
            alert('You do not have permission to edit posts.');
        }
    };

    window.deletePost = function(id) {
        if (isAdmin) {
            fetch(`${apiUrl}/${id}`, { method: 'DELETE' }).then(fetchPosts);
        } else {
            alert('You do not have permission to delete posts.');
        }
    };

    fetchPosts();
});
