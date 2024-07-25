// mainPage.js

document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('post-form');
    const postsContainer = document.getElementById('posts-container');
    const deleteAllBtn = document.getElementById('delete-all-btn');
    const viewByIdBtn = document.getElementById('view-by-id-btn');
    const viewAllBtn = document.getElementById('view-all-btn');
    const postIdInput = document.getElementById('post-id-input');
    const submitBtn = document.getElementById('submit-btn');
    const editBtn = document.getElementById('edit-btn');

    let posts = JSON.parse(localStorage.getItem('posts')) || [];
    let currentEditId = null;

    function renderPosts(postsToRender = posts) {
        postsContainer.innerHTML = '';
        postsToRender.forEach(post => {
            const postElement = document.createElement('div');
            postElement.classList.add('blog-post');
            postElement.innerHTML = `
                <h3>${post.title}</h3>
                <p class="author">By: ${post.author}</p>
                <p class="content">${post.content}</p>
                <div class="actions">
                    <button class="edit-btn" data-id="${post.id}">Edit</button>
                    <button class="delete-btn" data-id="${post.id}">Delete</button>
                </div>
            `;
            postsContainer.appendChild(postElement);
        });
    }

    function savePosts() {
        localStorage.setItem('posts', JSON.stringify(posts));
    }

    form.addEventListener('submit', function(e) {
        e.preventDefault();
        const title = document.getElementById('blogname').value;
        const author = document.getElementById('author').value;
        const content = document.getElementById('content').value;

        if (currentEditId) {
            const index = posts.findIndex(post => post.id === currentEditId);
            posts[index] = { ...posts[index], title, author, content };
            currentEditId = null;
            submitBtn.textContent = 'Add Post';
            editBtn.style.display = 'none';
        } else {
            const newPost = {
                id: Date.now().toString(),
                title,
                author,
                content
            };
            posts.push(newPost);
        }

        savePosts();
        renderPosts();
        form.reset();
    });

    postsContainer.addEventListener('click', function(e) {
        if (e.target.classList.contains('delete-btn')) {
            const id = e.target.getAttribute('data-id');
            posts = posts.filter(post => post.id !== id);
            savePosts();
            renderPosts();
        } else if (e.target.classList.contains('edit-btn')) {
            const id = e.target.getAttribute('data-id');
            const postToEdit = posts.find(post => post.id === id);
            document.getElementById('blogname').value = postToEdit.title;
            document.getElementById('author').value = postToEdit.author;
            document.getElementById('content').value = postToEdit.content;
            currentEditId = id;
            submitBtn.textContent = 'Update Post';
            editBtn.style.display = 'inline-block';
        }
    });

    deleteAllBtn.addEventListener('click', function() {
        posts = [];
        savePosts();
        renderPosts();
    });

    viewByIdBtn.addEventListener('click', function() {
        const id = postIdInput.value;
        // find the post with the given id
        const post = posts.find(post => post.id === id);


        // render the post if found
        if (post) {
            renderPosts([post]);
        } else {
            alert('Post not found');
        }
    });

    viewAllBtn.addEventListener('click', function() {
        renderPosts();
    });

    editBtn.addEventListener('click', function() {
        currentEditId = null;
        submitBtn.textContent = 'Add Post';
        editBtn.style.display = 'none';
        form.reset();
    });

    renderPosts();
});
