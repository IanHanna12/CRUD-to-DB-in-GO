document.addEventListener('DOMContentLoaded', function () {
    const isAdmin = localStorage.getItem('isAdmin') === 'true';
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;
    let prefetchedItems = [];

    document.getElementById('view-by-id-btn').addEventListener('click', viewPostById);
    document.getElementById('view-all-btn').addEventListener('click', fetchPosts);

    fetchWithAuth(`${apiUrl}/all-prefetch/delete`, {
        method: 'DELETE'
    })
        .then(response => {
            if (response.status === 200) {
                return response.json();
            } else {
                throw new Error('Error deleting all posts');
            }
        })
        .then(data => {
            prefetchedItems = [];
            postsContainer.innerHTML = '<p>No posts available</p>';
        })
        .catch(error => console.error('Error deleting all posts:', error));

    document.getElementById('delete-all-btn').addEventListener('click', deleteAllPosts);


    function viewPostById() {
        const id = document.getElementById('post-id-input').value;
        if (id) {
            const post = prefetchedItems.find(item => item.id === id);
            if (post) {
                postsContainer.innerHTML = renderPost(post);
            } else {
                fetchWithAuth(`${apiUrl}/single/${id}`)
                    .then(post => {
                        postsContainer.innerHTML = renderPost(post);
                    })
                    .catch(error => console.error('Error fetching post:', error));
            }
        }
    }

    function renderPost(post) {
        return `
            <div class="blog-post">
                <h3>${post.blogname}</h3>
                <p class="author">By: ${post.author}</p>
                <p class="content">${post.content}</p>
                <button onclick="editPost('${post.id}')" class="edit-btn">Edit</button>
                <button onclick="deletePost('${post.id}')" class="delete-btn">Delete</button>
            </div>
        `;
    }

    function fetchPosts() {
        fetchWithAuth(`${apiUrl}/all`)
            .then(renderPosts)
            .catch(error => console.error('Error fetching posts:', error));
    }

    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(renderPost).join('') || '<p>No posts available</p>';
    }

    function prefetchItems() {
        fetchWithAuth(`${apiUrl}/all-users`)
            .then(items => {
                prefetchedItems = items;
                console.log('Prefetched items:', items);
                renderPosts(items);
            })
            .catch(error => console.error('Error prefetching items:', error));
    }

    form.onsubmit = function (e) {
        e.preventDefault();
        const post = {
            blogname: form.blogname.value,
            author: form.author.value,
            content: form.content.value
        };
        fetchWithAuth(currentEditId ? `${apiUrl}/update/${currentEditId}` : `${apiUrl}/create`, {
            method: currentEditId ? 'PUT' : 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(post)
        }).then(updatedItem => {
            if (currentEditId) {
                const index = prefetchedItems.findIndex(item => item.id === currentEditId);
                if (index !== -1) {
                    prefetchedItems[index] = updatedItem;
                }
            } else {
                prefetchedItems.push(updatedItem);
            }
            currentEditId = null;
            form.reset();
            document.getElementById('submit-btn').textContent = 'Add Post';
            document.getElementById('edit-btn').style.display = 'none';
            renderPosts(prefetchedItems);
        }).catch(error => console.error('Error saving post:', error));
    };

    window.editPost = function (id) {
        currentEditId = id;
        document.getElementById('submit-btn').textContent = 'Update Post';
        document.getElementById('edit-btn').style.display = 'inline-block';
        const post = prefetchedItems.find(item => item.id === id);
        if (post) {
            populateForm(post);
        } else {
            fetchWithAuth(`${apiUrl}/single/${id}`)
                .then(populateForm)
                .catch(error => console.error('Error editing post:', error));
        }
    };

    function populateForm(post) {
        form.blogname.value = post.blogname;
        form.author.value = post.author;
        form.content.value = post.content;
        currentEditId = post.id;
        document.getElementById('submit-btn').textContent = 'Update Post';
        document.getElementById('edit-btn').style.display = 'inline-block';
        form.author.readOnly = true;
    }

    const editBtn = document.getElementById('edit-btn');
    if (editBtn) {
        editBtn.addEventListener('click', function () {
            currentEditId = null;
            form.reset();
            document.getElementById('submit-btn').textContent = 'Add Post';
            this.style.display = 'none';
            form.author.readOnly = false;
        });
    }

    function getSessionTokenFromCookie() {
        const cookies = document.cookie.split('; ');
        for (let cookie of cookies) {
            if (cookie.startsWith('session_token=')) {
                return cookie.split('=')[1];
            }
        }
        return null;
    }

    function fetchWithAuth(url, options = {}) {
        const sessionToken = getSessionTokenFromCookie();
        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${sessionToken}`,
            'isAdmin': isAdmin.toString()
        };
        options.credentials = 'include';
        return fetch(url, options)
            .then(response => {
                if (!response.ok) {
                    throw new Error('HTTP error! status: ' + response.status);
                }
                return response.json();
            });
    }

    // Call prefetchItems when the page loads
    prefetchItems();
});