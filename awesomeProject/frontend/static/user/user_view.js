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
        </div>
    `;
    }

    function fetchPosts() {
        if (prefetchedItems.length > 0) {
            renderPosts(prefetchedItems);
        } else {
            fetchWithAuth(`${apiUrl}/all`)
                .then(renderPosts)
                .catch(error => console.error('Error fetching posts:', error));
        }
    }

    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(renderPost).join('') || '<p>No posts available</p>';
    }

    function prefetchItems() {
        fetchWithAuth(`${apiUrl}/prefetch`)
            .then(response => {
                if (response && response.prefetchedItems) {
                    prefetchedItems = response.prefetchedItems;
                    console.log('Items prefetched successfully');
                } else {
                    prefetchedItems = [];
                    console.log('No items to prefetch');
                }
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
                .then(() => {
                    currentEditId = null;
                    form.reset();
                    document.getElementById('submit-btn').textContent = 'Add Post';
                    document.getElementById('edit-btn').style.display = 'none';
                    prefetchItems();
                    fetchPosts();
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
    }

    const editBtn = document.getElementById('edit-btn');
    if (editBtn) {
        editBtn.addEventListener('click', function () {
            currentEditId = null;
            form.reset();
            document.getElementById('submit-btn').textContent = 'Add Post';
            this.style.display = 'none';
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
            'Authorization': `Bearer ${sessionToken}`
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
