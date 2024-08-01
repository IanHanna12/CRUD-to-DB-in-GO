document.addEventListener('DOMContentLoaded', function() {
    const isAdmin = localStorage.getItem('isAdmin') === 'true';
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;

    document.getElementById('view-by-id-btn').addEventListener('click', viewPostById);
    document.getElementById('view-all-btn').addEventListener('click', fetchPosts);
    document.getElementById('delete-all-btn').addEventListener('click', deleteAllPosts);

    function viewPostById() {
        const id = document.getElementById('post-id-input').value;
        if (id) {
            fetchWithAuth(`${apiUrl}/single/${id}`)
                .then(res => res.json())
                .then(post => {
                    postsContainer.innerHTML = renderPost(post);
                });
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
        fetchWithAuth(`${apiUrl}/prefetch`)
            .then(res => res.json())
            .then(items => {
                console.log('Prefetched items:', items);
                renderPosts(items);
            })
            .catch(error => console.error('Error prefetching items:', error));
    }

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
        }).then(() => {
            currentEditId = null;
            form.reset();
            fetchPosts();
        });
    };

    window.editPost = function(id) {
        fetchWithAuth(`${apiUrl}/single/${id}`)
            .then(res => res.json())
            .then(post => {
                form.blogname.value = post.blogname;
                form.author.value = post.author;
                form.content.value = post.content;
                currentEditId = id;
                document.getElementById('submit-btn').textContent = 'Update Post';
                document.getElementById('edit-btn').style.display = 'inline-block';
            });
    };

    window.deletePost = function(id) {
        if (confirm('Are you sure you want to delete this post?')) {
            fetchWithAuth(`${apiUrl}/delete/${id}`, { method: 'DELETE' }).then(fetchPosts);
        }
    };

    function deleteAllPosts() {
        if (confirm('Are you sure you want to delete all posts?')) {
            fetchWithAuth(`${apiUrl}/all`, { method: 'DELETE' }).then(fetchPosts);
        }
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

    // Call prefetchItems when the page loads
    prefetchItems();
});
