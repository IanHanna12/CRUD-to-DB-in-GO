document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.querySelector('form');

    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const username = document.querySelector('input[type="text"]').value;
            const password = document.querySelector('input[type="password"]').value;
            const isAdmin = document.getElementById('isAdmin').checked;

            localStorage.setItem('isAdmin', isAdmin);

            // Redirect to appropriate view
            if (isAdmin) {
                window.location.href = 'http://localhost:8080/static/admin/admin_view.html';
            } else {
                window.location.href = 'http://localhost:8080/static/user/user_view.html';
            }
        });
    }
});
