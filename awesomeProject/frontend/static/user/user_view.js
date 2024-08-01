document.addEventListener('DOMContentLoaded', function() {
    const isAdmin = localStorage.getItem('isAdmin') === 'true';
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;

    const viewByIdBtn = document.getElementById('view-by-id-btn');

    viewByIdBtn.addEventListener('click', viewPostById);

    function viewPostById() {
        const id = document.getElementById('post-id-input').value;
        if (id) {
            fetchWithAuth(`${apiUrl}/single/${id}`)
                .then(post => {
                    postsContainer.innerHTML = renderPost(post);
                })
                .catch(error => console.error('Error fetching post:', error));
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
        fetchWithAuth(`${apiUrl}/all`)
            .then(renderPosts)
            .catch(error => console.error('Error fetching posts:', error));
    }

    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(renderPost).join('') || '<p>No posts available</p>';
    }

    function prefetchItems() {
        fetchWithAuth(`${apiUrl}/prefetch`)
            .then(items => {
                console.log('Prefetched items:', items);
                renderPosts(items);
            })
            .catch(error => console.error('Error prefetching items:', error));
    }

    if (form) {
        form.onsubmit = function(e) {
            e.preventDefault();
            const post = {
                blogname: form.blogname.value,
                author: form.author.value,
                content: form.content.value
            };
            fetchWithAuth(currentEditId ? `${apiUrl}/update/${currentEditId}` : `${apiUrl}/create`, {
                method: currentEditId ? 'PUT' : 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(post)
            })
                .then(() => {
                    currentEditId = null;
                    form.reset();
                    document.getElementById('submit-btn').textContent = 'Add Post';
                    document.getElementById('edit-btn').style.display = 'none';
                    fetchPosts();
                })
                .catch(error => console.error('Error saving post:', error));
        };
    }

    window.editPost = function(id) {
        currentEditId = id;
        document.getElementById('submit-btn').textContent = 'Update Post';
        document.getElementById('edit-btn').style.display = 'inline-block';
        fetchWithAuth(`${apiUrl}/single/${id}`)
            .then(post => {
                form.blogname.value = post.blogname;
                form.author.value = post.author;
                form.content.value = post.content;
                currentEditId = id;
                document.getElementById('submit-btn').textContent = 'Update Post';
                document.getElementById('edit-btn').style.display = 'inline-block';
            })
            .catch(error => console.error('Error editing post:', error));
    };

    const editBtn = document.getElementById('edit-btn');
    if (editBtn) {
        editBtn.addEventListener('click', function() {
            currentEditId = null;
            form.reset();
            document.getElementById('submit-btn').textContent = 'Add Post';
            this.style.display = 'none';
        });
    }

    function fetchWithAuth(url, options = {}) {
        options.headers = {
            ...options.headers,
            'isAdmin': isAdmin.toString()
        };
        return fetch(url, options)
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => {
                        throw new Error(text || `HTTP error! status: ${response.status}`);
                    });
                }
                return response.json();
            });
    }

    prefetchItems();
});