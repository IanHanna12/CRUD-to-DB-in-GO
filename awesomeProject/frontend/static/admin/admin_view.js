document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;
    let prefetchedItems = [];

    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(post => `
            <div class="blog-post">
                <h3>${post.blogname}</h3>
                <p>By: ${post.author}</p>
                <p>${post.content}</p>
                <button onclick="editPost('${post.id}')">Edit</button>
                <button onclick="deletePost('${post.id}')">Delete</button>
            </div>
        `).join('') || '<p>No posts available</p>';
    }

    function prefetchItems() {
        fetch(`${apiUrl}/prefetch/all`, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('sessionID')}`
            },
            credentials: 'include'
        })
            .then(res => res.json())
            .then(data => {
                prefetchedItems = data.prefetchedItems;
                renderPosts(prefetchedItems);
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
            fetch(currentEditId ? `${apiUrl}/${currentEditId}` : apiUrl, {
                method: currentEditId ? 'PUT' : 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('sessionID')}`
                },
                body: JSON.stringify(post),
                credentials: 'include'
            }).then(() => {
                currentEditId = null;
                form.reset();
                prefetchItems();
            });
        };
    }

    window.editPost = function(id) {
        const post = prefetchedItems.find(item => item.id === id);
        if (post) {
            form.style.display = 'block';
            form.blogname.value = post.blogname;
            form.author.value = post.author;
            form.content.value = post.content;
            currentEditId = id;
        }
    };

    window.deletePost = function(id) {
        fetch(`${apiUrl}/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('sessionID')}`
            },
            credentials: 'include'
        }).then(() => {
            prefetchedItems = prefetchedItems.filter(item => item.id !== id);
            renderPosts(prefetchedItems);
        });
    };

    const viewAllBtn = document.getElementById('view-all-btn');
    if (viewAllBtn) {
        viewAllBtn.addEventListener('click', prefetchItems);
    }

    // Initial prefetch
    prefetchItems();
});