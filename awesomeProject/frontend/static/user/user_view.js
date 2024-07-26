// This script manages the user view for the blog CRUD application
// It handles fetching, displaying, creating, and editing blog posts

document.addEventListener('DOMContentLoaded', function() {
    // Initialize variables and get DOM elements
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const apiUrl = 'http://localhost:8080/items';
    let currentEditId = null;

    // Set up event listeners for view buttons
    document.getElementById('view-by-id-btn').addEventListener('click', viewPostById);
    document.getElementById('view-all-btn').addEventListener('click', fetchPosts);

    // Function to view a single post by its ID
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

    // Function to render a single post HTML
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

    // Function to handle editing a post
    window.editPost = function (id) {
        fetch(`${apiUrl}/${id}`).then(res => res.json()).then(post => {
            // Populate form with post data
            form.blogname.value = post.blogname;
            form.author.value = post.author;
            form.content.value = post.content;
            currentEditId = id;
            // Update UI, change submit button text, and show edit button
            document.getElementById('submit-btn').textContent = 'Update Post';
            document.getElementById('edit-btn').style.display = 'inline-block';
        });
    };

    // Function to fetch all posts
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

    // Function to render all posts
    function renderPosts(posts) {
        postsContainer.innerHTML = posts.map(renderPost).join('') || '<p>No posts available</p>';
    }

    // Handle form submission for creating or updating a post
    form.onsubmit = function (e) {
        e.preventDefault();
        const post = {
            blogname: form.blogname.value,
            author: form.author.value,
            content: form.content.value
        };
        fetch(currentEditId ? `${apiUrl}/${currentEditId}` : apiUrl, {
            method: currentEditId ? 'PUT' : 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(post)
        })
            .then(() => {
                // Reset form and UI after successful submission
                currentEditId = null;
                form.reset();
                document.getElementById('submit-btn').textContent = 'Add Post';
                document.getElementById('edit-btn').style.display = 'none';
                fetchPosts();
            })
            .catch(error => console.error('Error saving post:', error));
    };

    // Handle canceling edit
    document.getElementById('edit-btn').addEventListener('click', function () {
        currentEditId = null;
        form.reset();
        document.getElementById('submit-btn').textContent = 'Add Post';
        this.style.display = 'none';
    });
})