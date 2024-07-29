document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.querySelector('form');

    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const username = document.querySelector('input[type="text"]').value;
            const password = document.querySelector('input[type="password"]').value;
            const isAdmin = document.getElementById('isAdmin').checked;

            fetch('http://localhost:8080/admin-login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password, isAdmin }),
            })
                .then(response => response.json())
                .then(data => {
                    data.sessionID = data.sessionID.toString();
                    if (data.success) {
                        localStorage.setItem('isAdmin', data.isAdmin);
                        if (data.isAdmin) {
                            localStorage.setItem('adminSessionID', data.sessionID);
                            window.location.href = 'http://localhost:8080/static/admin/admin_view.html';
                        } else {
                            window.location.href = 'http://localhost:8080/static/user/user_view.html';
                        }
                    } else {
                        alert('Login failed: ' + data.error);
                    }
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('An error occurred during login');
                });
        });
    }
});
