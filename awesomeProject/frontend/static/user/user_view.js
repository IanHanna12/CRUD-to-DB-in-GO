document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;
    let prefetchedItems = [];

    const username = localStorage.getItem('username');
    const authorInput = document.getElementById('author');
    if (username && authorInput) {
        authorInput.value = username;
        authorInput.readOnly = true;
    }

    const viewByIdBtn = document.getElementById('view-by-id-btn');
    if (viewByIdBtn) {
        viewByIdBtn.addEventListener('click', viewPostById);
    }

    function viewPostById() {
        const id = document.getElementById('post-id-input').value;
        if (id) {
            const post = prefetchedItems.find(item => item.id === id);
            if (post) {
                postsContainer.innerHTML = renderPost(post);
            } else {
                fetchWithAuth(`${apiUrl}/single/${id}`)
                    .then(post => {
                        if (post) {
                            postsContainer.innerHTML = renderPost(post);
                        } else {
                            postsContainer.innerHTML = '<p>Post not found</p>';
                        }
                    })
                    .catch(error => {
                        console.error('Error fetching post:', error);
                        postsContainer.innerHTML = '<p>Error fetching post</p>';
                    });
            }
        } else {
            postsContainer.innerHTML = '<p>Please enter a valid post ID</p>';
        }
    }

    function renderPost(post) {
        return `
        <div class="blog-post">
            <h3>${post.blogname}</h3>
            <p class="author">By: ${post.author}</p>
            <p class="content">${post.content}</p>
            <button onclick="editPost('${post.id}')" class="edit-btn">Edit</button>
        </div>
    `;
    }

    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(renderPost).join('') || '<p>No posts available</p>';
    }

    function prefetchItems() {
        fetchWithAuth(`${apiUrl}/all`)
            .then(items => {
                prefetchedItems = items;
                console.log('Items prefetched successfully');
            })
            .catch(error => {
                console.error('Error prefetching items:', error);
                prefetchedItems = [];
            });
    }

    if (form) {
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
            })
                .then(updatedItem => {
                    if (currentEditId) {
                        const index = prefetchedItems.findIndex(item => item.id === currentEditId);
                        if (index !== -1) {
                            prefetchedItems[index] = updatedItem;
                        }
                        postsContainer.innerHTML = renderPost(updatedItem);
                    } else {
                        prefetchedItems.push(updatedItem);
                        postsContainer.innerHTML = renderPost(updatedItem);
                    }
                    currentEditId = null;
                    form.reset();
                    document.getElementById('submit-btn').textContent = 'Add Post';
                    document.getElementById('edit-btn').style.display = 'none';
                })
                .catch(error => console.error('Error saving post:', error));
        };
    }

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
            if (cookie.startsWith('AuthToken=')) {
                return cookie.split('=')[1];
            }
        }
        return null;
    }

    function fetchWithAuth(url, options = {}) {
        const AuthToken = getSessionTokenFromCookie();
        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${AuthToken}`
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

    prefetchItems();
});
