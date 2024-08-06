document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;
    let prefetchedItems = [];

    function renderPosts(posts) {
        let postsArray;
        if (Array.isArray(posts)) {
            postsArray = posts;
        } else if (posts && posts.prefetchedItems) {
            postsArray = posts.prefetchedItems;
        } else {
            postsArray = [];
        }

        if (postsArray.length > 0) {
            postsContainer.innerHTML = postsArray.map(post => `
                <div class="blog-post">
                    <h3>${post.blogname}</h3>
                    <p>By: ${post.author}</p>
                    <p>${post.content}</p>
                    <button onclick="editPost('${post.id}')">Edit</button>
                    <button onclick="deletePost('${post.id}')">Delete</button>
                </div>
            `).join('');
        } else {
            postsContainer.innerHTML = '<p>No posts available</p>';
        }
    }

    function prefetchItems() {
        fetchWithAuth(`${apiUrl}/prefetch/all`)
            .then(items => {
                if (Array.isArray(items)) {
                    prefetchedItems = items;
                } else if (items && items.prefetchedItems) {
                    prefetchedItems = items.prefetchedItems;
                } else {
                    prefetchedItems = [];
                }
                console.log('Items prefetched successfully');
            })
            .catch(error => {
                console.error('Error prefetching items:', error);
                prefetchedItems = [];
            });
    }

    function renderAllPosts() {
        renderPosts(prefetchedItems);
    }

    function viewPostById() {
        const id = document.getElementById('post-id-input').value;
        if (id) {
            fetchWithAuth(`${apiUrl}/single/${id}`)
                .then(post => {
                    if (post) {
                        postsContainer.innerHTML = `
                            <div class="blog-post">
                                <h3>${post.blogname}</h3>
                                <p>By: ${post.author}</p>
                                <p>${post.content}</p>
                                <button onclick="editPost('${post.id}')">Edit</button>
                                <button onclick="deletePost('${post.id}')">Delete</button>
                            </div>
                        `;
                    } else {
                        postsContainer.innerHTML = '<p>Post not found</p>';
                    }
                })
                .catch(error => {
                    console.error('Error fetching post:', error);
                    postsContainer.innerHTML = '<p>Error fetching post</p>';
                });
        } else {
            postsContainer.innerHTML = '<p>Please enter a valid post ID</p>';
        }
    }

    if (form) {
        form.onsubmit = function(e) {
            e.preventDefault();
            const post = {
                blogname: form.blogname.value,
                author: form.author.value,
                content: form.content.value
            };
            let url;
            let method;
            if (currentEditId) {
                url = `${apiUrl}/admin/update/${currentEditId}`;
                method = 'PUT';
            } else {
                url = `${apiUrl}/create`;
                method = 'POST';
            }
            fetchWithAuth(url, {
                method: method,
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(post)
            })
                .then(updatedItem => {
                    if (currentEditId) {
                        const index = prefetchedItems.findIndex(item => item.id === currentEditId);
                        if (index !== -1) {
                            prefetchedItems[index] = updatedItem;
                        }
                    } else {
                        prefetchedItems = [...prefetchedItems, updatedItem];
                    }
                    currentEditId = null;
                    form.reset();
                    document.getElementById('submit-btn').textContent = 'Add Post';
                    document.getElementById('edit-btn').style.display = 'none';
                    renderAllPosts();
                })
                .catch(error => console.error('Error saving post:', error));
        };
    }

    window.editPost = function(id) {
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
    }

    window.deletePost = function(id) {
        fetchWithAuth(`${apiUrl}/delete/${id}`, {
            method: 'DELETE'
        }).then(() => {
            prefetchedItems = prefetchedItems.filter(item => item.id !== id);
            renderAllPosts();
        }).catch(error => console.error('Error deleting post:', error));
    };

    const viewAllBtn = document.getElementById('view-all-btn');
    if (viewAllBtn) {
        viewAllBtn.addEventListener('click', renderAllPosts);
    }

    const viewByIdBtn = document.getElementById('view-by-id-btn');
    if (viewByIdBtn) {
        viewByIdBtn.addEventListener('click', viewPostById);
    }

    const deleteAllBtn = document.getElementById('delete-all-btn');
    if (deleteAllBtn) {
        deleteAllBtn.addEventListener('click', () => {
            fetchWithAuth(`${apiUrl}/all`, {
                method: 'DELETE'
            }).then(() => {
                prefetchedItems = [];
                renderAllPosts();
            }).catch(error => console.error('Error deleting all posts:', error));
        });
    }

    function fetchWithAuth(url, options = {}) {
        const sessionToken = getSessionTokenFromCookie();
        if (!options.headers) {
            options.headers = {};
        }
        options.headers['Authorization'] = `Bearer ${sessionToken}`;
        options.credentials = 'include';
        return fetch(url, options)
            .then(response => {
                if (!response.ok) {
                    throw new Error('HTTP error! status: ' + response.status);
                }
                return response.json();
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

    prefetchItems();
});
