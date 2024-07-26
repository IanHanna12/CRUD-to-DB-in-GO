document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;

    document.getElementById('view-by-id-btn').addEventListener('click', viewPostById);
    document.getElementById('view-all-btn').addEventListener('click', fetchPosts);

    function viewPostById() {
        const id = document.getElementById('post-id-input').value;
        if (id) {
            fetch(`${apiUrl}/${id}`)
                .then(res => res.json())
                .then(post => {
                    postsContainer.innerHTML = renderPost(post);
                })
                .catch(error => {
                    console.error('Error fetching post:', error);
                    postsContainer.innerHTML = '<p>Error fetching post</p>';
                });
        }
    }

    function renderPost(post) {
        return `
            <div class="blog-post">
                <h3>${escapeHtml(post.blogname)}</h3>
                <p class="author">By: ${escapeHtml(post.author)}</p>
                <p class="content">${escapeHtml(post.content)}</p>
                <button onclick="editPost('${post.id}')" class="edit-btn">Edit</button>
            </div>
        `;
    }

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

    document.getElementById('edit-btn').addEventListener('click', function() {
        currentEditId = null;
        form.reset();
        document.getElementById('submit-btn').textContent = 'Add Post';
        this.style.display = 'none';
    });

    function escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    fetchPosts();
});
