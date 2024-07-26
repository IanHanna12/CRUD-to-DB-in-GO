document.addEventListener('DOMContentLoaded', function() {
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
            fetch(`${apiUrl}/${id}`)
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
        fetch(apiUrl)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                if (Array.isArray(data)) {
                    renderPosts(data);
                } else {
                    console.error('Invalid response data:', data);
                    renderPosts([]);
                }
            })
            .catch(error => {
                console.error('Error fetching posts:', error);
                renderPosts([]);
            });
    }


    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(renderPost).join('') || '<p>No posts available</p>';
    }

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

    window.editPost = function(id) {
        fetch(`${apiUrl}/${id}`).then(res => res.json()).then(post => {
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
            fetch(`${apiUrl}/${id}`, { method: 'DELETE' }).then(fetchPosts);
        }
    };

    function deleteAllPosts() {
        if (confirm('Are you sure you want to delete all posts?')) {
            fetch(apiUrl, { method: 'DELETE' }).then(fetchPosts);
        }
    }

    fetchPosts();
});
