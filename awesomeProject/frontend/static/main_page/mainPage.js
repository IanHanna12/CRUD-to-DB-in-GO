document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;

    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(post => `
            <div class="blog-post">
                <h3>${post.blogname}</h3>
                <p>By: ${post.author}</p>
                <p>${post.content}</p>
                <button onclick="editPost(${post.id})">Edit</button>
                <button onclick="deletePost(${post.id})">Delete</button>
            </div>
        `).join('') || '<p>No posts available</p>';
    }

    function fetchPosts() {
        fetch(apiUrl).then(res => res.json()).then(renderPosts);
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
        });
    };

    window.deletePost = function(id) {
        fetch(`${apiUrl}/${id}`, { method: 'DELETE' }).then(fetchPosts);
    };

    fetchPosts();
});